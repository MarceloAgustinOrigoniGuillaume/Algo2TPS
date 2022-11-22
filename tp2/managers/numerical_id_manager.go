package managers
import "tp2/interfaces"

// Implementacion de un id manager con identaficadores integer, usando un arreglo para rapidez

const _CAPACIDAD_INICIAL = 10


type numericalIdManager[T any] struct {
	elements []T
	ultimoId int
}

func CrearNumericalIdManager[T any]() interfaces.IdManager[int, T] {
	return &numericalIdManager[T]{make([]T, _CAPACIDAD_INICIAL), -1}
}

func initializedNumericalIdManager[T any](cantidadInicial int, initializer func(int) T) interfaces.IdManager[int, T] {
	if(cantidadInicial <= 0){
		return CrearNumericalIdManager[T]()
	}

	idManager := new(numericalIdManager[T])
	idManager.elements = make([]T, cantidadInicial)

	for i:=0;i<cantidadInicial;i++{
		idManager.elements[i] = initializer(i)
	}

	idManager.ultimoId = cantidadInicial

	return idManager
}

func (manager *numericalIdManager[T]) redimensionar(nuevoLargo int) {
	elementsNew := make([]T, nuevoLargo)

	copy(elementsNew, manager.elements)
	manager.elements = elementsNew
}


func (manager *numericalIdManager[T]) Agregar(element T) int {
	manager.ultimoId++

	if manager.ultimoId >= len(manager.elements) {
		manager.redimensionar(2 * manager.ultimoId)
	}

	manager.elements[manager.ultimoId] = element

	return manager.ultimoId
}

// no se verifica nada en obtener para mayor rapidez, te saltara el panic de indices si lo usas mal. Problema de Barbara.
// Debio usar el Existe antes.
func (manager *numericalIdManager[T]) Obtener(id int) T {
	return manager.elements[id]
}


func (manager *numericalIdManager[T]) Existe(id int) bool {
	return id >= 0 && id <= manager.ultimoId // si borrar fuese posible habria que chequear que sea valido el valor.
}


func (manager *numericalIdManager[T]) NuevoId() int{
	return manager.ultimoId+1
}



func (manager *numericalIdManager[T]) Iterar(visitar func(int,T) bool){
	i:= 0
	for i<= manager.ultimoId && visitar(i,manager.elements[i]){
		i++
	}
}




