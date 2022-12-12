package diccionario

//const (
//	ERROR_NO_ESTABA               = "La clave no pertenece al diccionario"
//	ERROR_ITERADOR_TERMINO        = "El iterador termino de iterar"
//)

type hashMap[K comparable, V any] struct {
	hashInterno map[K]V
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	hash := new(hashMap[K, V])
	hash.hashInterno = make(map[K]V)
	return hash
}


func (hash *hashMap[K, V]) Guardar(clave K, valor V) bool {
	_,estaba := hash.hashInterno[clave]

	hash.hashInterno[clave] = valor

	return estaba
}


func (hash *hashMap[K, V]) Pertenece(clave K) bool {
	_,estaba := hash.hashInterno[clave]
	return estaba
}

func (hash *hashMap[K, V]) dameElemento(clave K) V{
	elem,estaba := hash.hashInterno[clave]

	if !estaba {
		panic(ERROR_NO_ESTABA)
	}

	return elem	

}

func (hash *hashMap[K, V]) Obtener(clave K) V {
	return hash.dameElemento(clave)
}

func (hash *hashMap[K, V]) Borrar(clave K) V {
	elem := hash.dameElemento(clave)


	delete(hash.hashInterno,clave)

	return elem
}

func (hash *hashMap[K, V]) Cantidad() int {
	return len(hash.hashInterno)
}

func (hash *hashMap[K, V]) Iterar(visitar func(clave K, dato V) bool) {

	for clave,valor := range hash.hashInterno{
		if !visitar(clave,valor){
			return
		}
	}
}

// Iterador externo

func (hash *hashMap[K, V]) Iterador() IterDiccionario[K, V] {
	return creariteradorMap(hash)
}

type iteradorMap[K comparable, V any] struct {
	referencia *hashMap[K, V]
	claves []K
	posActual  int
}

func creariteradorMap[K comparable, V any](referencia *hashMap[K, V]) IterDiccionario[K, V] {
	iterador := new(iteradorMap[K, V])

	iterador.referencia = referencia
	iterador.claves = make([]K, referencia.Cantidad())
	i:= 0
	referencia.Iterar(func(clave K, _ V) bool{
		iterador.claves[i] = clave
		i++
		return true
	})
	return iterador
}

func (iterador *iteradorMap[K, V]) panicTermino() {
	if !iterador.HaySiguiente() {
		panic(ERROR_ITERADOR_TERMINO)
	}
}

func (iterador *iteradorMap[K, V]) HaySiguiente() bool {
	return iterador.posActual < len(iterador.claves)
}

func (iterador *iteradorMap[K, V]) VerActual() (K, V) {
	iterador.panicTermino()
	clave:= iterador.claves[iterador.posActual]
	return clave, iterador.referencia.Obtener(clave)
}

func (iterador *iteradorMap[K, V]) Siguiente() K {
	iterador.panicTermino()
	iterador.posActual++

	return iterador.claves[iterador.posActual-1]

}
