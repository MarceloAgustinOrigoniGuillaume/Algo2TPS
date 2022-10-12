package hashCuckoo

import "fmt"

const _CAPACIDAD_INICIAL = 128
const _MAXIMA_CARGA = 9 // esta constante tendria unidad de 10%, osea 9 = 90%
const ERROR_FUNCION_HASH = "Error: mala funcion de hash, redimension requeridad demasiadas veces seguidas"
const ERROR_NO_ESTABA = "La clave no pertenece al diccionario"
const ERROR_ITERADOR_TERMINO = "El iterador termino de iterar"

func toBytes[K comparable](objeto K) []byte{
	return []byte(fmt.Sprintf("%v",objeto))
}

func funcionHashingGenerica1(bytes []byte) int{
	if(len(bytes) == 0){
		return 0
	}

	return int(bytes[0] << 1)+ int(bytes[len(bytes)-1] >> 1) + len(bytes)
}

func funcionHashingGenerica2(bytes []byte) int{
	res := 256

	if(len(bytes) != 0){
		res += int(bytes[0])+ int(bytes[len(bytes)-1] << 2) + len(bytes)*2
	}
	return res
}

func funcionHashingGenerica3(bytes []byte) int{
	res := 100
	for i,dato:= range bytes{
		res += int(dato << 2) - i
	}

	if(res <0){
		res = -res
	}
	return res

}


type elementoCuckoo[K comparable, V any] struct{
	clave K
	valor V
	indiceFuncion int
}

func crearElementoCuckoo[K comparable, V any](clave K, valor V) *elementoCuckoo[K,V]{
	elemento:= new(elementoCuckoo[K,V])

	elemento.clave = clave
	elemento.valor = valor
	return elemento
}

func (elemento *elementoCuckoo[K,V]) damePosicionHash(funcionesHash []func([]byte) int) int{
	return funcionesHash[elemento.indiceFuncion](toBytes(elemento.clave))
}



type hashCuckoo[K comparable, V any] struct{
	elementos []*elementoCuckoo[K,V]
	cantidad int
	funcionesHash []func([]byte) int
}




func CrearHash[K comparable, V any]() Diccionario[K,V]{
	hash := new(hashCuckoo[K,V])

	hash.elementos = make([]*elementoCuckoo[K,V],_CAPACIDAD_INICIAL)
	hash.funcionesHash = []func([]byte) int {funcionHashingGenerica1,funcionHashingGenerica2,funcionHashingGenerica3}

	return hash
} 

// las unidades serian de 10%, aca lo que se haria es comparar el 100% de la cantidad con el (_MAXIMA_CARGA*10)% de la longitud(90%)
func (hash *hashCuckoo[K,V]) superoCargaPermitida() bool{
	return 10*hash.cantidad >= len(hash.elementos)*_MAXIMA_CARGA 
}

func reemplazoCuckoo[K comparable, V any](elementos []*elementoCuckoo[K,V],nuevoElemento *elementoCuckoo[K,V],funcionesHash []func([]byte) int) *elementoCuckoo[K,V]{
	pos:= funcionesHash[nuevoElemento.indiceFuncion](toBytes(nuevoElemento.clave)) % len(elementos)
	aGuardar:= elementos[pos]
	elementos[pos] = nuevoElemento

	if(aGuardar != nil){
		aGuardar.indiceFuncion++
		if aGuardar.indiceFuncion == len(funcionesHash){
			aGuardar.indiceFuncion = 0
		}
	}

	return aGuardar	
}

func insertarCuckoo[K comparable, V any](elementos []*elementoCuckoo[K,V],nuevoElemento *elementoCuckoo[K,V], funcionesHash []func([]byte) int) bool{
	posicionando := reemplazoCuckoo(elementos,nuevoElemento,funcionesHash)

	for (posicionando != nil && (posicionando != nuevoElemento || nuevoElemento.indiceFuncion != 0)) {
		posicionando = reemplazoCuckoo(elementos,posicionando,funcionesHash)
	}

	return posicionando == nil
}


func (hash *hashCuckoo[K,V]) buscarPosicionCuckoo(clave K) int{
	for _,funcionHash := range hash.funcionesHash{
		indice := funcionHash(toBytes(clave)) % len(hash.elementos)
		if(hash.elementos[indice] != nil && hash.elementos[indice].clave == clave){
			return indice
		}
	}

	return -1
}




func (hash *hashCuckoo[K,V]) reintentarInsertadoCuckoo(nuevo *elementoCuckoo[K,V]){
		// no deberia pasar mas de una vez
		intentos := 2
		for !insertarCuckoo(hash.elementos,nuevo,hash.funcionesHash){
			hash.redimensionar(2*len(hash.elementos))
			intentos--

			if(intentos <0){
				panic(ERROR_FUNCION_HASH) // no va a ser infinito...
			}

		}
}


func (hash *hashCuckoo[K,V]) redimensionar(nuevoLargo int){
	otra_vez := true
	multiplicador := 1
	var elementosNew []*elementoCuckoo[K,V]
	for otra_vez && multiplicador < 8{ // deberia ocurrir solo una vez, dios te salve si las funciones de hash son malas
		elementosNew = make([]*elementoCuckoo[K,V], multiplicador*nuevoLargo)
		otra_vez = false
		hash.Iterar(func (clave K, valor V) bool {
			if(!insertarCuckoo(elementosNew,crearElementoCuckoo(clave,valor),hash.funcionesHash)){ // no deberia pasar
				otra_vez = true
			}
			return !otra_vez
		})

		multiplicador *= 2
	}
	
	if(multiplicador == 8){
		panic(ERROR_FUNCION_HASH) // no va a ser infinito...
	}

	hash.elementos = elementosNew
}



func (hash *hashCuckoo[K,V]) Guardar(clave K, valor V){
	indice := hash.buscarPosicionCuckoo(clave)

	if(indice == -1){
		nuevo:= crearElementoCuckoo(clave,valor)
		hash.cantidad++

		if(hash.superoCargaPermitida() || !insertarCuckoo(hash.elementos,nuevo,hash.funcionesHash)){
			hash.reintentarInsertadoCuckoo(nuevo)
		}

	} else{
		hash.elementos[indice].valor = valor
	}

}


func (hash *hashCuckoo[K,V]) Pertenece(clave K) bool{
	return hash.buscarPosicionCuckoo(clave) != -1
}

func (hash *hashCuckoo[K,V]) Obtener(clave K) V{
	i:= hash.buscarPosicionCuckoo(clave)
	if(i == -1){
		panic(ERROR_NO_ESTABA)
	}
	return hash.elementos[i].valor
}

func (hash *hashCuckoo[K,V]) Borrar(clave K) V{
	i:= hash.buscarPosicionCuckoo(clave)
	
	if(i == -1){
		panic(ERROR_NO_ESTABA)
	}

	elem := hash.elementos[i].valor
	hash.elementos[i] = nil
	hash.cantidad--
	return elem
}

func (hash *hashCuckoo[K,V]) Cantidad() int{
	return hash.cantidad
}
	
func (hash *hashCuckoo[K,V]) Iterar(visitar func(clave K, dato V) bool){
	i:= 0
	for (i<len(hash.elementos) && (hash.elementos[i] == nil || visitar(hash.elementos[i].clave, hash.elementos[i].valor))){
		i++
	}	
}






// Iterador externo

func (hash *hashCuckoo[K,V]) Iterador() IterDiccionario[K,V]{
	return crearIteradorCuckoo(hash)
}

type iteradorCuckoo[K comparable, V any] struct{
	referencia *hashCuckoo[K,V]
	posActual int
}

func crearIteradorCuckoo[K comparable, V any](referencia *hashCuckoo[K,V]) IterDiccionario[K,V]{
	iterador := new(iteradorCuckoo[K,V])

	iterador.referencia = referencia
	iterador.posActual = -1
	iterador.iterarSiguiente()
	return iterador
}

func (iterador *iteradorCuckoo[K,V]) iterarSiguiente(){
	iterador.posActual++
	for iterador.posActual < len(iterador.referencia.elementos) && iterador.referencia.elementos[iterador.posActual] == nil{
		iterador.posActual++
	}
}
func (iterador *iteradorCuckoo[K,V]) panicTermino(){
	if(!iterador.HaySiguiente()){
		panic(ERROR_ITERADOR_TERMINO)
	}
}

func (iterador *iteradorCuckoo[K,V]) HaySiguiente() bool {
	return iterador.posActual < len(iterador.referencia.elementos)
}

func (iterador *iteradorCuckoo[K,V]) VerActual() (K, V) {
	iterador.panicTermino()
	elemento:= iterador.referencia.elementos[iterador.posActual]
	return elemento.clave,elemento.valor
}

func (iterador *iteradorCuckoo[K,V]) Siguiente() K {
	iterador.panicTermino()
	claveActual := iterador.referencia.elementos[iterador.posActual].clave
	iterador.iterarSiguiente()
	return claveActual



	
}