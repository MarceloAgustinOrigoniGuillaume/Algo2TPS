package cola_prioridad

const (
	_CAPACIDAD_INICIAL = 16
	_ESTA_VACIA        = "La cola esta vacia"
)

func swap[T comparable](elem1, elem2 *T) {
	*elem1, *elem2 = *elem2, *elem1
}

func downHeap[T comparable](elementos []T, cmp func(T, T) int, indice, ultimoIndice int) {
	swapInd := indice // asume el indice esta en rango

	candidato := (swapInd << 1) + 1 // izq
	if candidato <= ultimoIndice && cmp(elementos[swapInd], elementos[candidato]) < 0 {
		swapInd = candidato
	}

	candidato++ // der

	if candidato <= ultimoIndice && cmp(elementos[swapInd], elementos[candidato]) < 0 {
		swapInd = candidato
	}

	if swapInd != indice {
		swap(&elementos[indice], &elementos[swapInd])
		downHeap(elementos, cmp, swapInd, ultimoIndice)
	}
}

func inicializarHeapOrder[T comparable](elementos []T, cmp func(T, T) int, ultimoIndice int) {
	for i := (ultimoIndice >> 1) + 1; i >= 0; i-- { // se inicializa el pseudo orden, la segunda mitad se puede ignorar
		downHeap(elementos, cmp, i, ultimoIndice)
	}
}

func HeapSort[T comparable](arreglo []T, cmp func(T, T) int) {

	ultimoIndice := len(arreglo) - 1
	if ultimoIndice <= 0 { // si hay uno o esta vacio
		return
	}

	inicializarHeapOrder(arreglo, cmp, ultimoIndice) // pseudo ordenado

	for ultimoIndice > 0 {
		swap(&arreglo[0], &arreglo[ultimoIndice]) // garantizado el primero es el maximo.

		ultimoIndice--
		downHeap(arreglo, cmp, 0, ultimoIndice) // ya habiendo actualizado el indice se hace down heap desde la "raiz"
	}
}

type heapDinamico[T comparable] struct {
	elementos    []T
	ultimoIndice int
	cmp          func(T, T) int
}

func CrearHeap[T comparable](cmp func(T, T) int) ColaPrioridad[T] {
	heap := new(heapDinamico[T])
	heap.elementos = make([]T, _CAPACIDAD_INICIAL)
	heap.ultimoIndice = -1
	heap.cmp = cmp
	return heap
}

func CrearHeapArr[T comparable](arreglo []T, cmp func(T, T) int) ColaPrioridad[T] {
	heap := new(heapDinamico[T])
	heap.cmp = cmp
	heap.elementos = arreglo
	heap.ultimoIndice = len(arreglo) - 1

	if len(arreglo) < _CAPACIDAD_INICIAL { // se asegura la capacidad inicial
		heap.redimensionar(_CAPACIDAD_INICIAL)
	} else {
		heap.redimensionar(len(arreglo)) // se asegura que sea una copia
	}

	inicializarHeapOrder(heap.elementos, heap.cmp, heap.ultimoIndice)

	return heap
}

func (heap *heapDinamico[T]) redimensionar(nuevoLargo int) {
	elementosNew := make([]T, nuevoLargo)

	copy(elementosNew, heap.elementos)
	heap.elementos = elementosNew
}

func (heap *heapDinamico[T]) upHeap(indice int) {
	if indice == 0 {
		return
	}
	padre := (indice - 1) >> 1

	if heap.cmp(heap.elementos[padre], heap.elementos[indice]) < 0 {
		swap(&heap.elementos[padre], &heap.elementos[indice])
		heap.upHeap(padre)
	}
}

// EstaVacia devuelve true si la la cola se encuentra vacía, false en caso contrario.
func (heap *heapDinamico[T]) EstaVacia() bool {
	return heap.ultimoIndice == -1
}

// Encolar Agrega un elemento al heap.
func (heap *heapDinamico[T]) Encolar(elemento T) {
	heap.ultimoIndice++
	if heap.ultimoIndice >= len(heap.elementos) {
		heap.redimensionar(heap.ultimoIndice << 1)
	}

	heap.elementos[heap.ultimoIndice] = elemento

	heap.upHeap(heap.ultimoIndice)
}

func (heap *heapDinamico[T]) dameRaiz() *T {
	if heap.EstaVacia() {
		panic(_ESTA_VACIA)
	}

	return &heap.elementos[0]
}

// VerMax devuelve el elemento con máxima prioridad. Si está vacía, entra en pánico con un mensaje
// "La cola esta vacia".
func (heap *heapDinamico[T]) VerMax() T {

	return *heap.dameRaiz()
}

// Desencolar elimina el elemento con máxima prioridad, y lo devuelve. Si está vacía, entra en pánico con un
// mensaje "La cola esta vacia"
func (heap *heapDinamico[T]) Desencolar() T {

	swap(heap.dameRaiz(), &heap.elementos[heap.ultimoIndice])

	valor := heap.elementos[heap.ultimoIndice]
	// ultimo indice pasaria a ser la cantidad actual, antes de que se le reste.
	if len(heap.elementos) >= _CAPACIDAD_INICIAL<<1 && heap.ultimoIndice<<2 <= len(heap.elementos) {
		heap.redimensionar(heap.ultimoIndice << 1)
	}

	heap.ultimoIndice--
	downHeap(heap.elementos, heap.cmp, 0, heap.ultimoIndice) // ya habiendo actualizado el indice se hace down heap desde la raiz

	return valor
}

// Cantidad devuelve la cantidad de elementos que hay en la cola de prioridad.
func (heap *heapDinamico[T]) Cantidad() int {
	return heap.ultimoIndice + 1
}
