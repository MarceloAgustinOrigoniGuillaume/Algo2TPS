package cola_test

import (
	TDACola "cola"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func verificarColaEstaVacia[T any](t *testing.T, cola TDACola.Cola[T]) {
	require.NotNil(t, cola)
	require.PanicsWithValue(t, "La cola esta vacia", func() { cola.VerPrimero() })
	require.PanicsWithValue(t, "La cola esta vacia", func() { cola.Desencolar() })
	require.True(t, cola.EstaVacia())
}

func testAll[T any](t *testing.T, cola TDACola.Cola[T], testArreglo []T) {
	t.Log(fmt.Sprintf("Hacemos pruebas encolando %d elementos", len(testArreglo)))

	for _, valor := range testArreglo {
		cola.Encolar(valor)
	}

	t.Log(fmt.Sprintf("Hacemos pruebas viendo y desencolando %d elementos", len(testArreglo)))

	for i := 0; i < len(testArreglo); i++ {
		require.EqualValues(t, testArreglo[i], cola.VerPrimero()) // para saber que te desencola, ademas de que de el primero
		require.EqualValues(t, testArreglo[i], cola.Desencolar())
	}

}

func testEncolado[T any](t *testing.T, cola TDACola.Cola[T], testArreglo []T) {
	t.Log(fmt.Sprintf("Hacemos pruebas encolando %d elementos", len(testArreglo)))

	for _, valor := range testArreglo {
		cola.Encolar(valor)
	}

	t.Log(fmt.Sprintf("Hacemos pruebas desencolando %d elementos", len(testArreglo)))

	for i := 0; i < len(testArreglo); i++ {
		require.EqualValues(t, testArreglo[i], cola.Desencolar())
	}

}

func testEncoladoParcial[T any](t *testing.T, cola TDACola.Cola[T], testArreglo []T, aDesencolar int) []T {
	t.Log(fmt.Sprintf("Hacemos pruebas encolando %d elementos", len(testArreglo)))

	for _, valor := range testArreglo {
		cola.Encolar(valor)
	}

	if aDesencolar >= len(testArreglo) {
		aDesencolar = len(testArreglo) - 1
	}

	t.Log(fmt.Sprintf("Hacemos pruebas desencolando %d elementos", aDesencolar))

	for i := 0; i <= aDesencolar; i++ {
		require.EqualValues(t, testArreglo[i], cola.Desencolar())
	}

	return testArreglo[aDesencolar+1 : len(testArreglo)]

}

func testDesencolado[T any](t *testing.T, cola TDACola.Cola[T], testArreglo []T) {
	t.Log(fmt.Sprintf("Hacemos pruebas desencolando restantes de la cola, %d elementos", len(testArreglo)))
	for i := 0; i < len(testArreglo); i++ {
		require.EqualValues(t, testArreglo[i], cola.Desencolar())
	}
}

func testEncoladoYDesencoladoParcial[T any](t *testing.T, cola TDACola.Cola[T], testRemainder []T, testArreglo []T, aDesencolar int) []T {
	t.Log(fmt.Sprintf("Hacemos pruebas encolando %d elementos", len(testArreglo)))

	for _, valor := range testArreglo {
		cola.Encolar(valor)
	}

	testDesencolado[T](t, cola, testRemainder)

	if aDesencolar >= len(testArreglo) {
		aDesencolar = len(testArreglo) - 1
	}

	t.Log(fmt.Sprintf("Hacemos pruebas desencolando %d elementos", aDesencolar))

	for i := 0; i <= aDesencolar; i++ {
		require.EqualValues(t, testArreglo[i], cola.Desencolar())
	}

	return testArreglo[aDesencolar+1 : len(testArreglo)]

}

func TestVacia(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()

	t.Log("Hacemos pruebas con una cola Vacia")
	verificarColaEstaVacia[int](t, cola)
}

func TestEncoladoPocosElementos(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()

	testEncolado[int](t, cola, []int{2, 4, 8, 1, 3, 6, 5, 7, 9, 0})
	testEncolado[int](t, cola, []int{4, 8, 1, 3, 6, 5, 7, 0, 4, 32, 2, 1})

	t.Log("Se verifica que al final la cola esta vacia")
	verificarColaEstaVacia[int](t, cola)

}

func TestAllPocosElementos(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()

	testAll[int](t, cola, []int{2, 4, 8, 1, 3, 6, 5, 7, 9, 0})
	testAll[int](t, cola, []int{4, 8, 1, 3, 6, 5, 7, 0, 4, 32, 2, 1})

	t.Log("Se verifica que al final la cola esta vacia")
	verificarColaEstaVacia[int](t, cola)

}

func TestEncoladoParcialFinal(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()

	restante := testEncoladoParcial[int](t, cola, []int{2, 4, 8, 1, 3, 6, 5, 7, 9, 0}, 5)
	restante = testEncoladoYDesencoladoParcial[int](t, cola, restante, []int{4, 8, 1, 3, 6, 5, 7, 0, 4, 32, 2, 1}, 3)

	testDesencolado(t, cola, restante)

	t.Log("Se verifica que al final la cola este vacia")
	verificarColaEstaVacia[int](t, cola) // Esto es para saber si despues del trabajo de Encolador se volvio a estar vacia

}

func TestEncoladoParcialFinalStrings(t *testing.T) {
	t.Log("Test con strings para variar")

	cola := TDACola.CrearColaEnlazada[string]()

	restante := testEncoladoParcial[string](t, cola, []string{"UNO", "DOS", "TRES", "CUATRO", "CINCO"}, 2)
	restante = testEncoladoYDesencoladoParcial[string](t, cola, restante, []string{"OCHO", "DOS", "TRES", "CUATRO", "CINCO", "SEIS"}, 3)

	testDesencolado(t, cola, restante)

	t.Log("Se verifica que al final la cola este vacia")
	verificarColaEstaVacia[string](t, cola) // Esto es para saber si despues del trabajo de Encolador se volvio a estar vacia

}

func TestEncoladoParcialFinalVolumen(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()

	largo := 10000
	t.Log(fmt.Sprintf("Test encolando %d elementos", largo))

	for i := 0; i < largo; i++ {
		cola.Encolar(i)
	}

	for i := largo; i > 0; i-- {
		cola.Desencolar()
	}

	t.Log("Se verifica que al final la cola este vacia")
	verificarColaEstaVacia[int](t, cola) // Esto es para saber si despues del trabajo de Encolador se volvio a estar vacia

}
