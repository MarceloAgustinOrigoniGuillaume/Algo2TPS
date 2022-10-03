package lista

type IteradorLista[T any] interface {

	//HarSiguiente devuelve verdadero si hay un elemento actual al que ver
	HaySiguiente() bool

	// Devuelve el elemento Actual
	//Si se usa en un iterador que itero todos loselementos, entra en panico con
	//el mensaje "El iterador termino de iterar"
	VerActual() T

	// Devuelve el valor actual y se mueve al siguiente
	//Si se usa en un iterador que itero todos loselementos, entra en panico con
	//el mensaje "El iterador termino de iterar"
	Siguiente() T

	// Inserta en la posicion actual, haciendo el actual sea el siguiente
	Insertar(T)

	// Borra el actual y devuelve su valor
	//Si se usa en un iterador que itero todos loselementos, entra en panico con
	//el mensaje "El iterador termino de iterar"
	Borrar() T
}

type Lista[T any] interface {

	// EstaVacia devuelve verdadero si la lista no tiene elementos enlistados, false en caso contrario.
	EstaVacia() bool

	// VerPrimero obtiene el valor del primero de la lista. Si está vacía, entra en pánico con un mensaje
	// "La lista esta vacia".
	VerPrimero() T

	// VerPrimero obtiene el valor del ultimo de la lista. Si está vacía, entra en pánico con un mensaje
	// "La lista esta vacia".
	VerUltimo() T

	// agrega un nuevo elemento a la lista, al inicio de la misma.
	InsertarPrimero(T)

	// agrega un nuevo elemento a la lista, al final de la misma.
	InsertarUltimo(T)

	// Saca el primer elemento de la lista. Si la lista tiene elementos, se quita el primero de la misma,
	// y se devuelve ese valor. Si está vacía, entra en pánico con un mensaje "La lista esta vacia".
	BorrarPrimero() T

	// devuelve el largo de la lista.
	Largo() int

	//iterador externo

	Iterador() IteradorLista[T]

	// iterador interno
	Iterar(funcion func(T) bool)
}
