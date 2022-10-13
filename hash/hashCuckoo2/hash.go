package diccionario

import "fmt"
import "reflect"
import "hash/xxh3"

const _CAPACIDAD_INICIAL = 127
const _MAXIMA_CARGA = 7 // esta constante tendria unidad de 10%, osea 9 = 90%
const ERROR_FUNCION_HASH = "Error: mala funcion de hash, redimension requeridad demasiadas veces seguidas"
const ERROR_NO_ESTABA = "La clave no pertenece al diccionario"
const ERROR_ITERADOR_TERMINO = "El iterador termino de iterar"

func toBytes(objeto interface{}) []byte{
	switch objeto.(type){
		case string: // se chequea el tipo para saber cuando se puede usar una forma mas rapida
			return []byte(reflect.ValueOf(objeto).String())
		default:
			return []byte(fmt.Sprintf("%v",objeto))//*((*[]byte) unsafe.Pointer(reflect.ValueOf(objeto).Pointer()) )
	}
	//return []byte(fmt.Sprintf("%v",objeto))
}




func creatividad2(bytes []byte) uint64{
	if(len(bytes) == 0){
		return 0
	}
	i:= 0
	i2 := len(bytes)-1
	var res uint64= (uint64(bytes[0] << 1 | bytes[i2] >> 1)+ uint64(bytes[i2] >> 1))
	res^= res << 6
	res^= res >> 3
	i++
	i2--
	for i<3 && i < i2{
		res+= (uint64(bytes[i] << 1 | bytes[i2] >> 1)+ uint64(bytes[i2] >> 1))
		res= (res ^(res << 6)) ^ (res >> 3)
		i++
		i2--
	}

	return res
}

func _JenkinsHashFunction(bytes []byte) uint64{
	var res uint64 = 0
	for i:= 0; i< len(bytes) ; i++{
		res += uint64(bytes[i])
		res += res << 10;
		res ^= res >> 6;
	}

	return res
}

func puraCreatividad(bytes []byte) uint64{

	return xxh3.Hash(bytes) 
	/*
	res := 127 + len(bytes)

	for _,dato:= range bytes{
		res += int(dato <<2)
		res ^= (res<<3 ^ res >> 3)<<1
	}

	if(res <0){
		res = -res
	}

	return res
	*/
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

func hasheame(bytes []byte,funcion func([]byte) uint64,max int) int {
	return int(funcion(bytes) % uint64(max))
}

func (elemento *elementoCuckoo[K,V]) damePosicionHash(funcionesHash []func([]byte) uint64,max int) int{
	return hasheame(toBytes(elemento.clave),funcionesHash[elemento.indiceFuncion],max)
}



type hashCuckoo[K comparable, V any] struct{
	elementos []*elementoCuckoo[K,V]
	cantidad int
	funcionesHash []func([]byte) uint64
}

func crearTabla[K comparable, V any](largo int) []*elementoCuckoo[K,V]{
	return make([]*elementoCuckoo[K,V],largo)
}



func CrearHash[K comparable, V any]() Diccionario[K,V]{
	hash := new(hashCuckoo[K,V])

	hash.elementos = crearTabla[K,V](_CAPACIDAD_INICIAL)
	hash.funcionesHash = []func([]byte) uint64 {puraCreatividad,_JenkinsHashFunction}

	return hash
} 

// las unidades serian de 10%, aca lo que se haria es comparar el 100% de la cantidad con el (_MAXIMA_CARGA*10)% de la longitud(90%)
func (hash *hashCuckoo[K,V]) superoCargaPermitida() bool{
	return 10*hash.cantidad >= len(hash.elementos)*_MAXIMA_CARGA 
}

func reemplazoCuckoo[K comparable, V any](elementos []*elementoCuckoo[K,V],nuevoElemento *elementoCuckoo[K,V],funcionesHash []func([]byte) uint64) *elementoCuckoo[K,V]{
	pos:= nuevoElemento.damePosicionHash(funcionesHash,len(elementos))
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

func insertarCuckoo[K comparable, V any](elementos []*elementoCuckoo[K,V],nuevoElemento *elementoCuckoo[K,V], funcionesHash []func([]byte) uint64) bool{
	posicionando := reemplazoCuckoo(elementos,nuevoElemento,funcionesHash)

	for (posicionando != nil && (posicionando != nuevoElemento || nuevoElemento.indiceFuncion != 0)) {
		posicionando = reemplazoCuckoo(elementos,posicionando,funcionesHash)
	}

	return posicionando == nil
}


func (hash *hashCuckoo[K,V]) buscarPosicionCuckoo(clave K) int{
	bytes := toBytes(clave)
	for _,funcionHash := range hash.funcionesHash{
		indice := hasheame(bytes,funcionHash,len(hash.elementos))
		if(hash.elementos[indice] != nil && hash.elementos[indice].clave == clave){
			return indice
		}
	}

	return -1
}




func (hash *hashCuckoo[K,V]) reintentarInsertadoCuckoo(nuevo *elementoCuckoo[K,V]){
		// no deberia pasar mas de una vez
		hash.redimensionar(2*len(hash.elementos))
		if !insertarCuckoo(hash.elementos,nuevo,hash.funcionesHash){
			panic(ERROR_FUNCION_HASH) // no va a ser infinito...
		}
}


func (hash *hashCuckoo[K,V]) redimensionar(nuevoLargo int){

	//fmt.Printf("\nREDIMENSIONANDO CUANDO %d porciento\n",((100.0)*hash.cantidad)/len(hash.elementos))

	elementosNew := crearTabla[K,V](2*nuevoLargo)
	hash.Iterar(func (clave K, valor V) bool {
		if(!insertarCuckoo(elementosNew,crearElementoCuckoo(clave,valor),hash.funcionesHash)){ // no deberia pasar
			panic(ERROR_FUNCION_HASH) // no va a ser infinito...
		}
		return true
	})
	
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