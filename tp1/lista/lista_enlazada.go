package lista

type nodoLista[T any] struct {
	valor     T
	siguiente *nodoLista[T]
}

type listaEnlazada[T any] struct {
	primero *nodoLista[T]
	ultimo  *nodoLista[T]
	largo   int
}

func crearNodoLista[T any](valor T) *nodoLista[T] {
	nuevoNodoLista := new(nodoLista[T])
	nuevoNodoLista.valor = valor

	return nuevoNodoLista
}

func CrearListaEnlazada[T any]() Lista[T] {
	return new(listaEnlazada[T])
}

func (lista *listaEnlazada[T]) EstaVacia() bool {
	return lista.primero == nil || lista.ultimo == nil

}

func (lista *listaEnlazada[T]) Largo() int {
	return lista.largo
}

func (lista *listaEnlazada[T]) panicEstaVacia() { // aunque no se cambie nada se usa un puntero para no copiar la struct

	if lista.EstaVacia() { // lista no puede ser nil
		panic("La lista esta vacia")
	}
}

// VerPrimero obtiene el valor del primero de la lista. Si está vacía, entra en pánico con un mensaje
// "La lista esta vacia".
func (lista *listaEnlazada[T]) VerPrimero() T {
	lista.panicEstaVacia()

	return lista.primero.valor

}

func (lista *listaEnlazada[T]) VerUltimo() T {
	lista.panicEstaVacia()

	return lista.ultimo.valor

}

// Agrega un nuevo elemento a la lista, al final de la misma.
func (lista *listaEnlazada[T]) InsertarUltimo(valor T) {
	nuevoUltimo := crearNodoLista(valor)

	if lista.EstaVacia() {
		lista.primero = nuevoUltimo
	} else {
		lista.ultimo.siguiente = nuevoUltimo
	}
	lista.ultimo = nuevoUltimo

	lista.largo++
}

// Inserta como primero en la lista

func (lista *listaEnlazada[T]) InsertarPrimero(valor T) {

	// Se podria setear directamente a lista.primero.... pero para poder usar lista.EstaVacia sin confusion alguna decidimos que no
	nuevoPrimero := crearNodoLista(valor)
	nuevoPrimero.siguiente = lista.primero

	if lista.EstaVacia() {
		lista.ultimo = nuevoPrimero
	}

	lista.primero = nuevoPrimero

	lista.largo++

}

// Saca el primer elemento de la lista. Si la lista tiene elementos, se quita el primero de la misma,
// y se devuelve ese valor. Si está vacía, entra en pánico con un mensaje "La lista esta vacia".
func (lista *listaEnlazada[T]) BorrarPrimero() T {
	lista.panicEstaVacia()

	res := lista.primero.valor

	lista.primero = lista.primero.siguiente

	if lista.primero == nil {
		lista.ultimo = nil
	}

	lista.largo--

	return res

}

func (lista *listaEnlazada[T]) Iterador() IteradorLista[T] {
	return crearIterador[T](lista)
}

func (lista *listaEnlazada[T]) Iterar(haceAlgo func(T) bool) {
	nodoActual := lista.primero

	for nodoActual != nil && haceAlgo(nodoActual.valor) {
		nodoActual = nodoActual.siguiente
	}
}

type iteradorListaEnlazada[T any] struct {
	listaReferencia *listaEnlazada[T]
	actual          *nodoLista[T]
	anterior        *nodoLista[T]
}

func crearIterador[T any](listaReferencia *listaEnlazada[T]) IteradorLista[T] {
	iterador := new(iteradorListaEnlazada[T])
	iterador.listaReferencia = listaReferencia
	iterador.actual = listaReferencia.primero
	return iterador
}

func (iterador *iteradorListaEnlazada[T]) HaySiguiente() bool {
	return iterador.actual != nil
}

func (iterador *iteradorListaEnlazada[T]) esUltimo(ultimo *nodoLista[T]) bool {
	return ultimo == iterador.listaReferencia.ultimo
}
func (iterador *iteradorListaEnlazada[T]) panicIteradorRecorrido() {
	if !iterador.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
}

func (iterador *iteradorListaEnlazada[T]) VerActual() T {
	iterador.panicIteradorRecorrido()
	return iterador.actual.valor
}

func (iterador *iteradorListaEnlazada[T]) Siguiente() T {
	iterador.panicIteradorRecorrido()
	iterador.anterior = iterador.actual
	iterador.actual = iterador.actual.siguiente
	return iterador.anterior.valor
}

func (iterador *iteradorListaEnlazada[T]) siEsPrimero(nuevoValor *nodoLista[T]) {
	if iterador.actual == iterador.listaReferencia.primero {
		iterador.listaReferencia.primero = nuevoValor
	} else {
		iterador.anterior.siguiente = nuevoValor
	}
}

// Inserta en la posicion actual, haciendo el actual sea el siguiente
func (iterador *iteradorListaEnlazada[T]) Insertar(valor T) {
	nuevoActual := crearNodoLista(valor)

	iterador.siEsPrimero(nuevoActual)

	nuevoActual.siguiente = iterador.actual
	iterador.actual = nuevoActual

	if iterador.esUltimo(iterador.anterior) {
		iterador.listaReferencia.ultimo = iterador.actual
	}
	iterador.listaReferencia.largo++
}

// Borra el actual y devuelve su valor
// Si se usa en un iterador que itero todos loselementos, entra en panico con
// el mensaje "El iterador termino de iterar"
func (iterador *iteradorListaEnlazada[T]) Borrar() T {
	iterador.panicIteradorRecorrido()
	valor := iterador.actual.valor

	iterador.siEsPrimero(iterador.actual.siguiente)
	if iterador.esUltimo(iterador.actual) {
		iterador.listaReferencia.ultimo = iterador.anterior
	}
	iterador.actual = iterador.actual.siguiente
	iterador.listaReferencia.largo--
	return valor
}
