package cola_prioridad_test

import (
	TDAHeap "cola_prioridad"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func llamadosBinariosRango(minimo int, maximo int, visitar func(int)) {
	if minimo+1 >= maximo {
		return
	}
	medio := (minimo + maximo) >> 1
	visitar(medio)
	llamadosBinariosRango(minimo, medio, visitar)
	llamadosBinariosRango(medio, maximo, visitar)
}

func crearArregloNumeros(minimo int, maximo int) []int {
	arr := make([]int, maximo-minimo)
	i := 0
	llamadosBinariosRango(minimo, maximo+1, func(elem int) {
		arr[i] = elem
		i++
	})
	return arr
}

const ERROR_VACIO = "La cola esta vacia"

func funcionCompararBasicaStrings(elemento1 string, elemento2 string) int {
	res := 0
	resaux := 0
	for i, c := range elemento1 {
		res += (i + 1) * int(c)
		resaux += int(c)
	}

	for i, c := range elemento2 {
		res -= (i + 1) * int(c)
		resaux -= int(c)
	}
	if res == 0 && resaux != 0 {
		return resaux
	}
	return res
}

func funcionCompararBasicaInts(elemento1 int, elemento2 int) int {
	return elemento1 - elemento2
}

func testHeapVacio[T comparable](t *testing.T, heap TDAHeap.ColaPrioridad[T]) {
	require.Zero(t, heap.Cantidad(), "La cantidad no era 0")
	require.True(t, heap.EstaVacia(), "Esta vacia era false")
	require.PanicsWithValue(t, ERROR_VACIO, func() { heap.Desencolar() })
	require.PanicsWithValue(t, ERROR_VACIO, func() { heap.VerMax() })
}

func TestHeapVacio(t *testing.T) {
	t.Log("Comprueba que el Heap vacio no tiene elementos y los panics")
	testHeapVacio(t, TDAHeap.CrearHeap[string](funcionCompararBasicaStrings))
}
func TestHeapAgregadoArregloVacio(t *testing.T) {
	t.Log("Se testeara crear heap arr con un arreglo vacio y se verificara lo este")
	testHeapVacio(t, TDAHeap.CrearHeapArr[int]([]int{}, funcionCompararBasicaInts))
}
func TestHeapAgregadoArregloSinModificar(t *testing.T) {
	t.Log("Crear un heap con un arreglo con los numeros del 1 al 32 de forma desordenada, se verificara no se modifique el arreglo original")
	arr1 := crearArregloNumeros(0, 32)
	arr2 := crearArregloNumeros(0, 32)
	TDAHeap.CrearHeapArr[int](arr1, funcionCompararBasicaInts)

	for i, valor := range arr1 {
		require.EqualValues(t, valor, arr2[i], fmt.Sprintf("No eran iguales en posicion %d", i))
	}
}

func TestHeapAgregadoBasico(t *testing.T) {
	t.Log("Agrega los numeros del 1 al 32 de forma desordenada y verifica que la cantidad,VerMax y  desencolar den el valor correcto")
	heap := TDAHeap.CrearHeap[int](funcionCompararBasicaInts)

	llamadosBinariosRango(0, 33, heap.Encolar)

	require.EqualValues(t, 32, heap.Cantidad(), "No tuvo la cantidad correcta")
	require.EqualValues(t, 32, heap.VerMax(), "No mostro el maximo correctamente")
	require.EqualValues(t, 32, heap.Desencolar(), "No dio el maximo correctamente")
}

func TestHeapAgregadoArregloBasico(t *testing.T) {
	t.Log("Agrega un arreglo con los numeros del 1 al 32 de forma desordenada y despues el 10 y verifica que la cantidad, VerMax y  desencolar den el valor correcto")

	heap := TDAHeap.CrearHeapArr[int](crearArregloNumeros(0, 32), funcionCompararBasicaInts)

	heap.Encolar(10)
	require.EqualValues(t, 33, heap.Cantidad(), "No tuvo la cantidad correcta")
	require.EqualValues(t, 32, heap.VerMax(), "No mostro el maximo correctamente")
	require.EqualValues(t, 32, heap.Desencolar(), "No dio el maximo correctamente")

}

func TestHeapSort(t *testing.T) {
	arr := crearArregloNumeros(0, 32)

	TDAHeap.HeapSort(arr, funcionCompararBasicaInts) // como comparar devuelve si es mayor heapsort ordenaria de menor a mayor

	for i := 0; i < len(arr); i++ {
		require.EqualValues(t, i+1, arr[i], "No ordeno correctamente el heap")
	}
}

func TestHeapDesencoladoArregloBasico(t *testing.T) {
	t.Log("Agrega un arreglo con los numeros del 1 al 32 de forma desordenada y despues el 0 y verifica que la cantidad, VerMax y  desencolar den el valor correcto")

	heap := TDAHeap.CrearHeapArr[int](crearArregloNumeros(0, 32), funcionCompararBasicaInts)

	heap.Encolar(0)
	require.EqualValues(t, 33, heap.Cantidad(), "No tuvo la cantidad correcta despues de agregar")

	for i := 32; i >= 0; i-- {
		require.EqualValues(t, i, heap.VerMax(), "No mostro el maximo correctamente")
		require.EqualValues(t, i, heap.Desencolar(), "No dio el maximo correctamente")
	}

	testHeapVacio(t, heap)
}

func testHeapDesencolado(t *testing.T, minimo, maximo int) {
	t.Log(fmt.Sprintf("Agrega los numeros del %d al %d de forma desordenada y verifica que la cantidad VerMax y  desencolar den el valor correcto, hasta que se vacia el heap, se verificara que despues este vacia", minimo, maximo))
	heap := TDAHeap.CrearHeap[int](funcionCompararBasicaInts)

	llamadosBinariosRango(minimo-1, maximo+1, heap.Encolar)

	require.EqualValues(t, maximo-minimo+1, heap.Cantidad(), "No tuvo la cantidad correcta despues de agregar")

	for i := maximo; i > 0; i-- {
		require.EqualValues(t, i, heap.VerMax(), "No mostro el maximo correctamente")
		require.EqualValues(t, i, heap.Desencolar(), "No dio el maximo correctamente")
	}

	testHeapVacio(t, heap)
}

func TestHeapDesencoladoBasico(t *testing.T) {
	testHeapDesencolado(t, 1, 32)
}

func TestHeapVolumen(t *testing.T) {
	testHeapDesencolado(t, 1, 200000)
}
