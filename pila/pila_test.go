package pila_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	TDAPila "pila"
	"testing"
)

func verificarPilaEstaVacia[T any](t *testing.T, pila TDAPila.Pila[T]) {
	require.NotNil(t, pila)

	require.PanicsWithValue(t, "La pila esta vacia", func() { pila.VerTope() })
	require.PanicsWithValue(t, "La pila esta vacia", func() { pila.Desapilar() })
	require.True(t, pila.EstaVacia())
}

func testAllDesdeArreglo[T any](t *testing.T, pila TDAPila.Pila[T], testArreglo []T) {
	t.Log(fmt.Sprintf("Hacemos pruebas apilando %d elementos", len(testArreglo)))

	for _, valor := range testArreglo {
		pila.Apilar(valor)
	}

	t.Log(fmt.Sprintf("Hacemos pruebas viendo y desapilando %d elementos", len(testArreglo)))

	for i := len(testArreglo) - 1; i >= 0; i-- {
		require.EqualValues(t, testArreglo[i], pila.VerTope()) // para saber que te no desapila, ademas de que de el ultimo
		require.EqualValues(t, testArreglo[i], pila.Desapilar())
	}

}

func testApiladoDesdeArreglo[T any](t *testing.T, pila TDAPila.Pila[T], testArreglo []T) {
	t.Log(fmt.Sprintf("Hacemos pruebas apilando %d elementos", len(testArreglo)))

	for _, valor := range testArreglo {
		pila.Apilar(valor)
	}

	t.Log(fmt.Sprintf("Hacemos pruebas desapilando %d elementos", len(testArreglo)))

	for i := len(testArreglo) - 1; i >= 0; i-- {
		require.EqualValues(t, testArreglo[i], pila.Desapilar())
	}

}

func testApiladoParcialDesdeArreglo[T any](t *testing.T, pila TDAPila.Pila[T], testArreglo []T, aDespilar int) int {
	t.Log(fmt.Sprintf("Hacemos pruebas apilando %d elementos", len(testArreglo)))

	for _, valor := range testArreglo {
		pila.Apilar(valor)
	}

	restantes := 0
	if aDespilar < len(testArreglo) {
		restantes = len(testArreglo) - aDespilar
	}

	t.Log(fmt.Sprintf("Hacemos pruebas desapilando %d elementos", aDespilar))

	for i := len(testArreglo) - 1; i >= restantes; i-- {
		require.EqualValues(t, testArreglo[i], pila.Desapilar())
	}

	return restantes

}

func TestPilaVacia(t *testing.T) {

	t.Log("Hacemos pruebas con pila Vacia")
	pila := TDAPila.CrearPilaDinamica[int]()
	verificarPilaEstaVacia[int](t, pila)

}

func TestAllPocosElementos(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	testAllDesdeArreglo[int](t, pila, []int{2, 4, 8, 1, 3, 6, 5, 7, 9, 0})
	testAllDesdeArreglo[int](t, pila, []int{4, 8, 1, 3, 6, 5, 7, 0, 4, 32, 2, 1})

	t.Log("Se verifica que al final la pila esta vacia")
	require.True(t, pila.EstaVacia()) // Esto es para saber si despues del trabajo de apilador se volvio a estar vacia

}

func TestApiladoPocosElementos(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	testApiladoDesdeArreglo[int](t, pila, []int{2, 4, 8, 1, 3, 6, 5, 7, 9, 0})
	testApiladoDesdeArreglo[int](t, pila, []int{4, 8, 1, 3, 6, 5, 7, 0, 4, 32, 2, 1})

	t.Log("Se verifica que al final la pila esta vacia")
	require.True(t, pila.EstaVacia()) // Esto es para saber si despues del trabajo de apilador se volvio a estar vacia

}

func TestApiladoParcialPocosElementos(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	restante := testApiladoParcialDesdeArreglo[int](t, pila, []int{2, 4, 8, 1, 3, 6, 5, 7, 9, 0}, 5)
	restante += testApiladoParcialDesdeArreglo[int](t, pila, []int{4, 8, 1, 3, 6, 5, 7, 0, 4, 32, 2, 1}, 3)

	t.Log(fmt.Sprintf("Desapilando restantes %d elementos", restante))

	for restante > 0 {
		pila.Desapilar()
		restante--
	}

	t.Log("Se verifica que al final la pila este vacia")
	verificarPilaEstaVacia[int](t, pila) // Esto es para saber si despues del trabajo de apilador se volvio a estar vacia

}

func TestApiladoParcialFinal(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	restante := testApiladoParcialDesdeArreglo[int](t, pila, []int{2, 4, 8, 1, 3, 6, 5, 7, 9, 0}, 5)
	restante += testApiladoParcialDesdeArreglo[int](t, pila, []int{4, 8, 1, 3, 6, 5, 7, 0, 4, 32, 2, 1}, 3)

	t.Log(fmt.Sprintf("Viendo y Desapilando restantes %d elementos", restante))

	for restante > 0 {
		pila.VerTope()
		pila.Desapilar()
		restante--
	}

	t.Log("Se verifica que al final la pila este vacia")
	verificarPilaEstaVacia[int](t, pila) // Esto es para saber si despues del trabajo de apilador se volvio a estar vacia

}

func TestApiladoParcialFinalStrings(t *testing.T) {
	t.Log("Test con strings para variar")

	pila := TDAPila.CrearPilaDinamica[string]()

	restante := testApiladoParcialDesdeArreglo[string](t, pila, []string{"UNO", "DOS", "TRES", "CUATRO", "CINCO"}, 2)
	restante += testApiladoParcialDesdeArreglo[string](t, pila, []string{"UNO", "DOS", "TRES", "CUATRO", "CINCO", "SEIS"}, 3)

	t.Log(fmt.Sprintf("Viendo y Desapilando restantes %d elementos", restante))

	for restante > 0 {
		pila.VerTope()
		pila.Desapilar()
		restante--
	}

	t.Log("Se verifica que al final la pila este vacia")
	verificarPilaEstaVacia[string](t, pila) // Esto es para saber si despues del trabajo de apilador se volvio a estar vacia

}

func TestApiladoParcialFinalVolumen(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	largo := 10000
	t.Log(fmt.Sprintf("Test apilando %d elementos", largo))

	for i := 0; i < largo; i++ {
		pila.Apilar(i)
	}

	for i := largo; i > 0; i-- {
		pila.Desapilar()
	}

	t.Log("Se verifica que al final la pila este vacia")
	verificarPilaEstaVacia[int](t, pila) // Esto es para saber si despues del trabajo de apilador se volvio a estar vacia

}
