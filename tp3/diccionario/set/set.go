package set

type Set[K comparable] interface {

	// Guardar guarda el par clave-dato en el Set. Si la clave ya se encontraba, se actualiza el dato asociado
	Guardar(dato K) bool

	// Pertenece determina si una clave ya se encuentra en el Set, o no
	Pertenece(dato K) bool

	// Borrar borra del Set la clave indicada. Si la clave no
	// pertenece al Set, debe entrar en pánico con un mensaje 'La clave no pertenece al Set'
	Borrar(dato K) bool

	// Cantidad devuelve la cantidad de elementos dentro del Set
	Cantidad() int

	// Iterar itera internamente el Set, aplicando la función pasada por parámetro a todos los elementos del
	// mismo
	Iterar(func(dato K) bool)

	// Iterador devuelve un IterSet para este Set
	Iterador() IterSet[K]
}

type IterSet[K comparable] interface {

	// HaySiguiente devuelve si hay más datos para ver. Esto es, si en el lugar donde se encuentra parado
	// el iterador hay un elemento.
	HaySiguiente() bool

	// VerActual devuelve la clave y el dato del elemento actual en el que se encuentra posicionado el iterador.
	// Si no HaySiguiente, debe entrar en pánico con el mensaje 'El iterador termino de iterar'
	VerActual() K

	// Siguiente si HaySiguiente, devuelve la clave actual (equivalente a VerActual, pero únicamente la clave), y
	// además avanza al siguiente elemento en el Set. Si no HaySiguiente, entonces debe entrar en pánico con
	// mensaje 'El iterador termino de iterar'
	Siguiente() K
}
