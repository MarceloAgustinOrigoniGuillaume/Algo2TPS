package pila

const capacidadInicial = 10

type pilaDinamica[T any] struct {
	datos        []T
	ultimoIndice int // se decidio usar ultimoIndice en vez de cantidad, eso te ahorra los operaciones constantes -1
}

func (pila *pilaDinamica[T]) redimensionar(nuevoLargo int) {
	datosNew := make([]T, nuevoLargo)

	copy(datosNew, pila.datos)
	pila.datos = datosNew
}

func CrearPilaDinamica[T any]() Pila[T] {
	pila := new(pilaDinamica[T])
	pila.datos = make([]T, capacidadInicial)
	pila.ultimoIndice = -1

	return pila

}

func (pila *pilaDinamica[T]) panicEstaVacia() { // aunque no se cambie nada se usa un puntero para no copiar la struct

	if pila.EstaVacia() { // pila no puede ser nil
		panic("La pila esta vacia")
	}
}

func (pila *pilaDinamica[T]) EstaVacia() bool {
	return pila.ultimoIndice == -1
}

func (pila *pilaDinamica[T]) VerTope() T {
	pila.panicEstaVacia()

	return pila.datos[pila.ultimoIndice]
}

func (pila *pilaDinamica[T]) Apilar(elemento T) {

	pila.ultimoIndice++

	if pila.ultimoIndice >= len(pila.datos) {
		pila.redimensionar(2 * pila.ultimoIndice)
	}

	pila.datos[pila.ultimoIndice] = elemento

}

func (pila *pilaDinamica[T]) Desapilar() T {
	pila.panicEstaVacia()

	var elemento T = pila.datos[pila.ultimoIndice]

	if len(pila.datos) > 2*capacidadInicial && len(pila.datos) >= 4*pila.ultimoIndice {
		pila.redimensionar(2 * pila.ultimoIndice)
	}

	pila.ultimoIndice--

	return elemento
}
