package diccionario

import "fmt"

const _CAPACIDAD_INICIAL = 127
const _MAXIMA_CARGA = 6 // esta constante tendria unidad de 10%, osea 6 = 60%
const ERROR_NO_ESTABA = "La clave no pertenece al diccionario"
const ERROR_ITERADOR_TERMINO = "El iterador termino de iterar"

func toBytes[K comparable](objeto K) []byte{
	return []byte(fmt.Sprintf("%v",objeto))
}

func _JenkinsHashFunction(bytes []byte) int{
	res := 0
	for i:= 0; i< len(bytes) ; i++{
		res += int(bytes[i])
		res += res << 10;
		res ^= res >> 6;
	}

	return res
}


type elementoCerrado[K comparable, V any] struct{
	clave K
	valor V
	estado int
}

func crearElementoCerrado[K comparable, V any](clave K, valor V) elementoCerrado[K,V]{
	return elementoCerrado[K,V]{clave,valor,1}
}

func crearElementoCerradoVacio[K comparable, V any]() elementoCerrado[K,V]{
	return elementoCerrado[K,V]{}
}


func crearTabla[K comparable, V any](largo int) []elementoCerrado[K,V]{
	tabla:= make([]elementoCerrado[K,V],largo)
	for i:= range tabla{
		tabla[i] = crearElementoCerradoVacio[K,V]()
	}
	return tabla
}

type hashCerrado[K comparable, V any] struct{
	elementos []elementoCerrado[K,V]
	cantidad int
}




func CrearHash[K comparable, V any]() Diccionario[K,V]{
	hash := new(hashCerrado[K,V])
	hash.elementos = crearTabla[K,V](_CAPACIDAD_INICIAL)
	return hash
} 

// las unidades serian de 10%, aca lo que se haria es comparar el 100% de la cantidad con el (_MAXIMA_CARGA*10)% de la longitud(60%)
func (hash *hashCerrado[K,V]) superoCargaPermitida() bool{
	return 10*hash.cantidad >= len(hash.elementos)*_MAXIMA_CARGA 
}


func iterarPosicionCerrado(posInicial int,maximo int,visitar func(int) bool){

	seguir := visitar(posInicial)

	i:= posInicial+1


	for seguir && i<maximo{
		seguir = visitar(i)
		i++
	}

	i = 0
	for seguir && i<posInicial{
		seguir = visitar(i)
		i++
	}
}

func insertarCerrado[K comparable, V any](elementos []elementoCerrado[K,V],nuevoElemento elementoCerrado[K,V]) bool{

	agregoNuevo := true

	iterarPosicionCerrado(_JenkinsHashFunction(toBytes(nuevoElemento.clave)) % len(elementos),len(elementos),
	 func(indice int) bool{
	 	if(elementos[indice].estado == 0 || elementos[indice].clave == nuevoElemento.clave){
			agregoNuevo = elementos[indice].estado == 0
			nuevoElemento.estado = 1
			elementos[indice] = nuevoElemento
			return false
		}
		return true
	 })

	 return agregoNuevo
}


func (hash *hashCerrado[K,V]) buscarPosicion(clave K) int{
	res := -1
	iterarPosicionCerrado(_JenkinsHashFunction(toBytes(clave)) % len(hash.elementos),len(hash.elementos),
	 func(indice int) bool{
		if(hash.elementos[indice].estado == 0){
			return false
		}

		if(hash.elementos[indice].clave == clave){
			res = indice
			return false
		}

		return true
	})

	return res
}


func (hash *hashCerrado[K,V]) redimensionar(){
	nuevos :=  crearTabla[K,V](2 * len(hash.elementos))
	
	hash.Iterar(func (clave K, valor V) bool {
		insertarCerrado(nuevos,crearElementoCerrado(clave,valor))
		return true
	})

	hash.elementos = nuevos

}



func (hash *hashCerrado[K,V]) Guardar(clave K, valor V){
	if(hash.superoCargaPermitida()){
		hash.redimensionar()
	}

	if(insertarCerrado(hash.elementos,crearElementoCerrado(clave,valor))){
		hash.cantidad ++ 
	}
}


func (hash *hashCerrado[K,V]) Pertenece(clave K) bool{
	return hash.buscarPosicion(clave) != -1
}

func (hash *hashCerrado[K,V]) Obtener(clave K) V{
	i:= hash.buscarPosicion(clave)
	if(i == -1){
		panic(ERROR_NO_ESTABA)
	}
	return hash.elementos[i].valor
}

func (hash *hashCerrado[K,V]) Borrar(clave K) V{
	i:= hash.buscarPosicion(clave)
	
	if(i == -1){
		panic(ERROR_NO_ESTABA)
	}

	elem := hash.elementos[i].valor
	hash.elementos[i] = crearElementoCerradoVacio[K,V]()
	hash.elementos[i].estado = -1
	hash.cantidad--
	return elem
}

func (hash *hashCerrado[K,V]) Cantidad() int{
	return hash.cantidad
}
	
func (hash *hashCerrado[K,V]) Iterar(visitar func(clave K, dato V) bool){
	i:= 0
	for (i<len(hash.elementos) && (hash.elementos[i].estado != 1 || visitar(hash.elementos[i].clave, hash.elementos[i].valor))){
		i++
	}	
}






// Iterador externo

func (hash *hashCerrado[K,V]) Iterador() IterDiccionario[K,V]{
	return creariteradorCerrado(hash)
}

type iteradorCerrado[K comparable, V any] struct{
	referencia *hashCerrado[K,V]
	posActual int
}

func creariteradorCerrado[K comparable, V any](referencia *hashCerrado[K,V]) IterDiccionario[K,V]{
	iterador := new(iteradorCerrado[K,V])

	iterador.referencia = referencia
	iterador.posActual = -1
	iterador.iterarSiguiente()
	return iterador
}

func (iterador *iteradorCerrado[K,V]) iterarSiguiente(){
	iterador.posActual++
	for iterador.posActual < len(iterador.referencia.elementos) && iterador.referencia.elementos[iterador.posActual].estado != 1{
		iterador.posActual++
	}
}
func (iterador *iteradorCerrado[K,V]) panicTermino(){
	if(!iterador.HaySiguiente()){
		panic(ERROR_ITERADOR_TERMINO)
	}
}

func (iterador *iteradorCerrado[K,V]) HaySiguiente() bool {
	return iterador.posActual < len(iterador.referencia.elementos)
}

func (iterador *iteradorCerrado[K,V]) VerActual() (K, V) {
	iterador.panicTermino()
	elemento:= iterador.referencia.elementos[iterador.posActual]
	return elemento.clave,elemento.valor
}

func (iterador *iteradorCerrado[K,V]) Siguiente() K {
	iterador.panicTermino()
	claveActual := iterador.referencia.elementos[iterador.posActual].clave
	iterador.iterarSiguiente()
	return claveActual



	
}