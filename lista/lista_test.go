package lista_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	TDALista "lista"
	"testing"
)

func verificarListaEstaVacia[T any](t *testing.T, lista TDALista.Lista[T]) {
	require.NotNil(t, lista)
	require.PanicsWithValue(t, "La lista esta vacia", func() { lista.VerPrimero() })
	require.PanicsWithValue(t, "La lista esta vacia", func() { lista.VerUltimo() })
	require.PanicsWithValue(t, "La lista esta vacia", func() { lista.BorrarPrimero() })

	require.EqualValues(t, 0, lista.Largo())
	require.True(t, lista.EstaVacia())

	verificarNoHaySiguiente[T](t, lista.Iterador())

}
func verificarNoHaySiguiente[T any](t *testing.T, iterador TDALista.IteradorLista[T]) {
	require.NotNil(t, iterador)
	require.False(t, iterador.HaySiguiente())
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iterador.Borrar() })
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iterador.VerActual() })
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iterador.Siguiente() })
}

func testInsertarUltimo[T any](t *testing.T, lista TDALista.Lista[T], testArreglo []T) {
	t.Log(fmt.Sprintf("Hacemos pruebas insertando al final %d elementos", len(testArreglo)))

	largo := lista.Largo()
	for i, valor := range testArreglo {
		lista.InsertarUltimo(valor)
		largo++
		require.EqualValues(t, testArreglo[i], lista.VerUltimo())
		require.EqualValues(t, largo, lista.Largo())
	}
}

func testInsertarPrimero[T any](t *testing.T, lista TDALista.Lista[T], testArreglo []T) {
	t.Log(fmt.Sprintf("Hacemos pruebas insertando al inicio %d elementos", len(testArreglo)))
	largo := lista.Largo()
	for i, valor := range testArreglo {
		lista.InsertarPrimero(valor)
		largo++
		require.EqualValues(t, testArreglo[i], lista.VerPrimero())
		require.EqualValues(t, largo, lista.Largo())
	}
}

func insertarPrimeroOrdenadamente[T any](t *testing.T, lista TDALista.Lista[T], testArreglo []T) {
	iterador := lista.Iterador()

	for _, valor := range testArreglo {
		iterador.Insertar(valor)
		iterador.Siguiente()
	}
}

func testBorradoPrimeroDirecto[T any](t *testing.T, lista TDALista.Lista[T], testArreglo []T) {
	t.Log(fmt.Sprintf("Hacemos pruebas viendo y borrando el primero, de %d elementos", len(testArreglo)))

	for i := 0; i < len(testArreglo); i++ {
		require.EqualValues(t, testArreglo[i], lista.VerPrimero()) // para saber que no te borra, ademas de que de el primero
		require.EqualValues(t, testArreglo[i], lista.BorrarPrimero())
	}

}

func testBorradoPrimeroInverso[T any](t *testing.T, lista TDALista.Lista[T], testArreglo []T) {
	t.Log(fmt.Sprintf("Hacemos pruebas viendo y borrando el primero, de %d elementos", len(testArreglo)))

	for i := len(testArreglo) - 1; i >= 0; i-- {
		require.EqualValues(t, testArreglo[i], lista.VerPrimero()) // para saber que no te borra, ademas de que de el primero
		require.EqualValues(t, testArreglo[i], lista.BorrarPrimero())
	}

}

func testInsertarYBorrarElementos[T any](t *testing.T, arregloCola []T, arregloPila []T) {
	lista := TDALista.CrearListaEnlazada[T]()

	verificarListaEstaVacia[T](t, lista)

	testInsertarPrimero[T](t, lista, arregloCola)
	testInsertarUltimo[T](t, lista, arregloPila)

	testBorradoPrimeroInverso[T](t, lista, arregloCola)
	testBorradoPrimeroDirecto[T](t, lista, arregloPila)

	t.Log("Se verifica que al final la lista esta vacia")
	verificarListaEstaVacia[T](t, lista)

}

func aplicarEnIteradorDesde[T any](t *testing.T, lista TDALista.Lista[T], testArreglo []T, indice int,
	insertarOBorrar func(TDALista.IteradorLista[T], T, *int)) {
	iterador := lista.Iterador()

	i := 0
	for iterador.HaySiguiente() && i < indice { // Se recorre hasta el final
		iterador.Siguiente()
	}

	largo := lista.Largo()

	for _, valor := range testArreglo { // Se recorre el arreglo
		insertarOBorrar(iterador, valor, &largo)
	}

}

func crearIteradorBorrar[T any](t *testing.T, lista TDALista.Lista[T], chequeo func(T)) func(TDALista.IteradorLista[T], T, *int) {

	//func(*testing.T, TDALista.Lista[T],TDALista.IteradorLista[T],T, *int)
	return func(iterador TDALista.IteradorLista[T], valor T, largo *int) {

		chequeo(valor)
		require.EqualValues(t, valor, iterador.VerActual())
		require.EqualValues(t, valor, iterador.Borrar())
		(*largo)--
		require.EqualValues(t, (*largo), lista.Largo())

	}
}

func crearIteradorInsertar[T any](t *testing.T, lista TDALista.Lista[T], chequeo func(T, TDALista.IteradorLista[T])) func(TDALista.IteradorLista[T], T, *int) {

	return func(iterador TDALista.IteradorLista[T], valor T, largo *int) {
		iterador.Insertar(valor)
		(*largo)++
		require.EqualValues(t, (*largo), lista.Largo())
		require.EqualValues(t, valor, iterador.VerActual())
		chequeo(valor, iterador)
	}

}

func testInsertadoYBorradoExterno[T any](t *testing.T, arreglo []T, invArreglo []T, valorTest T) {

	t.Log("Hacemos pruebas insertando y borrando elementos mediante el iterador externo en la primer posicion")

	lista := TDALista.CrearListaEnlazada[T]()

	aplicarEnIteradorDesde(t, lista, arreglo, 0,
		crearIteradorInsertar[T](t, lista, func(valor T, iterador TDALista.IteradorLista[T]) {
			require.EqualValues(t, valor, lista.VerPrimero())
		}))

	aplicarEnIteradorDesde(t, lista, invArreglo, 0,
		crearIteradorBorrar[T](t, lista, func(valor T) {
			require.EqualValues(t, valor, lista.VerPrimero())
		}))

	verificarListaEstaVacia(t, lista)

	t.Log("Hacemos pruebas insertando mediante el iterador externo en la ultima posicion")

	lista.InsertarPrimero(valorTest)

	aplicarEnIteradorDesde(t, lista, arreglo, lista.Largo(),
		crearIteradorInsertar[T](t, lista, func(valor T, iterador TDALista.IteradorLista[T]) {
			require.EqualValues(t, valor, lista.VerUltimo())
			iterador.Siguiente()
			verificarNoHaySiguiente(t, iterador)
		}))

	require.EqualValues(t, valorTest, lista.VerPrimero())

	t.Log("Insertamos un elemento en la segunda posicion y borramos el de la tercera. Despues se borra el ultimo y se verifica que se hayan borrando")

	iterador := lista.Iterador()
	iterador.Siguiente() // ie indice 0 del arreglo

	iterador.Insertar(valorTest) // inserto valor en la posicion del indice 0 del arreglo

	iterador.Siguiente() // ie al indice 0 del arreglo

	require.EqualValues(t, arreglo[0], iterador.VerActual())
	require.EqualValues(t, arreglo[0], iterador.Borrar()) // borro el indice 0 y fue al indice 1

	require.EqualValues(t, arreglo[1], iterador.VerActual())

	iterador2 := lista.Iterador()
	//iterador := lista.Iterador()
	iterador2.Siguiente()
	require.EqualValues(t, iterador2.VerActual(), valorTest)
	iterador2.Siguiente()

	require.EqualValues(t, iterador2.VerActual(), arreglo[1])

	for i := 3; i < lista.Largo(); i++ { // i = 3 porque ya nos movimos de posicion 3 veces, uno al crear y dos con Siguiente
		iterador2.Siguiente()

	}
	require.EqualValues(t, lista.VerUltimo(), iterador2.Borrar())
	require.EqualValues(t, lista.VerUltimo(), arreglo[len(arreglo)-2])
	verificarNoHaySiguiente(t, iterador2)

}

func agregarSumando(t *testing.T, lista TDALista.Lista[int], arreglo []int, sumaActual int) int {

	for _, valor := range arreglo {
		sumaActual += valor
		lista.InsertarPrimero(valor)
	}
	t.Log(fmt.Sprintf("Se agrego %d elementos con una suma de %d", len(arreglo), sumaActual))

	return sumaActual
}

func testIterarSuma(t *testing.T, lista TDALista.Lista[int], valorSuma int) {
	t.Log(fmt.Sprintf("Se va a testear ir sumando los elementos hasta encontrar el primero 0 y se requerira sea %d", valorSuma))

	sum := 0
	lista.Iterar(func(num int) bool {
		sum += num
		return num != 0
	})

	require.EqualValues(t, valorSuma, sum)

}

func TestVacia(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	t.Log("Hacemos pruebas con una lista Vacia")
	verificarListaEstaVacia[int](t, lista)
}

func TestInsertarPocosEnteros(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	testInsertarPrimero[int](t, lista, []int{2, 4, 8, 1, 3, 6, 5, 7, 9, 0})
	testInsertarUltimo[int](t, lista, []int{4, 8, 1, 3, 6, 5, 7, 0, 4, 32, 2, 1})

	t.Log("Se verifica que al final la lista no esta vacia")
	require.False(t, lista.EstaVacia())

}

func TestPocosElementos(t *testing.T) {

	t.Log("Test insertando viendo y borrando pocos elementos ints")

	testInsertarYBorrarElementos[int](t, []int{2, 4, 8, 1, 3, 6, 5, 7, 9, 0},
		[]int{4, 8, 1, 3, 6, 5, 7, 0, 4, 32, 2, 1})

	t.Log("Test insertando viendo y borrando pocos elementos floats")

	testInsertarYBorrarElementos[float32](t, []float32{1.2, 1.3, 0.4}, []float32{1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7})
	t.Log("Test insertando pocos elementos strings")

	testInsertarYBorrarElementos[string](t, []string{"a", "b", "c"},
		[]string{"c", "d", "f", "e", "g", "h", "u", "v", "w", "x", "y", "z"})
}

func TestVolumen(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	largo := 100000
	t.Log(fmt.Sprintf("Test insertando %d elementos", largo))

	for i := 0; i < largo; i++ {
		lista.InsertarUltimo(i)
		require.EqualValues(t, i, lista.VerUltimo())
	}

	require.EqualValues(t, largo, lista.Largo())
	for i := 0; i < largo; i++ {
		require.EqualValues(t, i, lista.BorrarPrimero())
	}

	t.Log("Se verifica que al final la lista este vacia")
	verificarListaEstaVacia[int](t, lista) // Esto es para saber si despues del trabajo de insertador se volvio a estar vacia

}

func TestInsertadoYBorradoExternoLlamado(t *testing.T) {
	t.Log("Test externo ints")

	testInsertadoYBorradoExterno[int](t, []int{1, 2, 3, 4, 5, 6, 8, 9}, []int{9, 8, 6, 5, 4, 3, 2, 1}, 69)
	t.Log("Test externo string")

	testInsertadoYBorradoExterno[string](t, []string{"a", "b", "c", "d", "e", "f"}, []string{"f", "e", "d", "c", "b", "a"}, "j")

	t.Log("Test externo floats")

	testInsertadoYBorradoExterno[float32](t, []float32{1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7}, []float32{1.7, 1.6, 1.5, 1.4, 1.3, 1.2, 1.1}, 5.5)

}

func TestIterarInternoSuma(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[int]()

	testInsertarUltimo[int](t, lista, []int{0, 4, 32, 2, 1})
	suma := agregarSumando(t, lista, []int{4, 8, 1, 3, 6, 5, 7}, 0)

	testInsertarUltimo[int](t, lista, []int{4, 1, 2, 3, 4})

	testIterarSuma(t, lista, suma)
}

func TestIterarInternoBytes(t *testing.T) {
	lista := TDALista.CrearListaEnlazada[string]()

	for _, valor := range []string{"!", " ", "M", "a", "k", "i", "s", "e"} {
		lista.InsertarUltimo(valor)
	}
	insertarPrimeroOrdenadamente[string](t, lista, []string{"T", "E", "S", "T"})

	t.Log(fmt.Sprintf("Se va a agarrar de la lista de caracteres hasta el primer espacio y se requerira sea TEST!"))

	res := ""
	lista.Iterar(func(caracter string) bool {
		if caracter == " " {
			return false
		}
		res += caracter
		return true
	})

	t.Log(fmt.Sprintf("Se va a agarrar de la lista de caracteres hasta el final y se verificara sea TEST! Makise"))
	require.EqualValues(t, "TEST!", res)

	res = ""
	lista.Iterar(func(caracter string) bool {
		res += string(caracter)
		return true
	})
	require.EqualValues(t, "TEST! Makise", string(res))
}
