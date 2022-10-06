package cola

type nodo[T any] struct {
	valor     T
	siguiente *nodo[T]
}

type colaEnlazada[T any] struct {
	primero *nodo[T]
	ultimo  *nodo[T]
}

func crearNodo[T any](valor T) *nodo[T] {
	nuevo_nodo := new(nodo[T])
	nuevo_nodo.valor = valor

	return nuevo_nodo
}

func CrearColaEnlazada[T any]() Cola[T] {
	return new(colaEnlazada[T])
}

func (cola *colaEnlazada[T]) EstaVacia() bool {
	return cola.primero == nil
}

func (cola *colaEnlazada[T]) panicEstaVacia() { // aunque no se cambie nada se usa un puntero para no copiar la struct

	if cola.EstaVacia() { // cola no puede ser nil
		panic("La cola esta vacia")
	}
}

// VerPrimero obtiene el valor del primero de la cola. Si está vacía, entra en pánico con un mensaje
// "La cola esta vacia".
func (cola *colaEnlazada[T]) VerPrimero() T {
	cola.panicEstaVacia()

	return cola.primero.valor

}

// Encolar agrega un nuevo elemento a la cola, al final de la misma.
func (cola *colaEnlazada[T]) Encolar(valor T) {
	nuevo_ultimo := crearNodo(valor)

	if cola.EstaVacia() {
		cola.primero = nuevo_ultimo
	} else {
		cola.ultimo.siguiente = nuevo_ultimo
	}

	cola.ultimo = nuevo_ultimo

}

// Desencolar saca el primer elemento de la cola. Si la cola tiene elementos, se quita el primero de la misma,
// y se devuelve ese valor. Si está vacía, entra en pánico con un mensaje "La cola esta vacia".
func (cola *colaEnlazada[T]) Desencolar() T {
	cola.panicEstaVacia()

	res := cola.primero.valor

	cola.primero = cola.primero.siguiente

	if cola.EstaVacia() {
		cola.ultimo = nil
	}

	return res

}
