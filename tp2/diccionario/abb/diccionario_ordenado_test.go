package abb_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
	TDAABB "tp2/diccionario/abb"
)

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
func funcionCompararBasicaIntsS(elemento1 int, elemento2 int) int {
	fmt.Printf("%d - %d = %d\n", elemento1, elemento2, elemento1-elemento2)
	return elemento1 - elemento2
}

// TEST DE LA CATEDRA PARA DICCIONARIO ADAPTADOS

func TestABBVacio(t *testing.T) {
	t.Log("Comprueba que ABB vacio no tiene claves")
	dic := TDAABB.CrearABB[string, string](funcionCompararBasicaStrings)
	require.EqualValues(t, 0, dic.Cantidad())
	require.False(t, dic.Pertenece("A"))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Obtener("A") })
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar("A") })
}

func TestUnElement(t *testing.T) {
	t.Log("Comprueba que ABB con un elemento tiene esa Clave, unicamente")
	dic := TDAABB.CrearABB[string, int](funcionCompararBasicaStrings)
	dic.Guardar("A", 10)
	require.EqualValues(t, 1, dic.Cantidad())
	require.True(t, dic.Pertenece("A"))
	require.False(t, dic.Pertenece("B"))
	require.EqualValues(t, 10, dic.Obtener("A"))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Obtener("B") })
}

func TestABBGuardar(t *testing.T) {
	t.Log("Guarda algunos pocos elementos en el ABB, y se comprueba que en todo momento funciona acorde")
	clave1 := "Gato"
	clave2 := "Perro"
	clave3 := "Vaca"
	valor1 := "miau"
	valor2 := "guau"
	valor3 := "moo"
	claves := []string{clave1, clave2, clave3}
	valores := []string{valor1, valor2, valor3}

	dic := TDAABB.CrearABB[string, string](funcionCompararBasicaStrings)
	require.False(t, dic.Pertenece(claves[0]))
	require.False(t, dic.Pertenece(claves[0]))
	dic.Guardar(claves[0], valores[0])
	require.EqualValues(t, 1, dic.Cantidad())
	require.True(t, dic.Pertenece(claves[0]))
	require.True(t, dic.Pertenece(claves[0]))
	require.EqualValues(t, valores[0], dic.Obtener(claves[0]))
	require.EqualValues(t, valores[0], dic.Obtener(claves[0]))

	require.False(t, dic.Pertenece(claves[1]))
	require.False(t, dic.Pertenece(claves[2]))
	dic.Guardar(claves[1], valores[1])
	require.True(t, dic.Pertenece(claves[0]))
	require.True(t, dic.Pertenece(claves[1]))
	require.EqualValues(t, 2, dic.Cantidad())
	require.EqualValues(t, valores[0], dic.Obtener(claves[0]))
	require.EqualValues(t, valores[1], dic.Obtener(claves[1]))

	require.False(t, dic.Pertenece(claves[2]))
	dic.Guardar(claves[2], valores[2])
	require.True(t, dic.Pertenece(claves[0]))
	require.True(t, dic.Pertenece(claves[1]))
	require.True(t, dic.Pertenece(claves[2]))
	require.EqualValues(t, 3, dic.Cantidad())
	require.EqualValues(t, valores[0], dic.Obtener(claves[0]))
	require.EqualValues(t, valores[1], dic.Obtener(claves[1]))
	require.EqualValues(t, valores[2], dic.Obtener(claves[2]))
}

func TestClaveVacia(t *testing.T) {
	t.Log("Guardamos una clave vacía (i.e. \"\") y deberia funcionar sin problemas")
	dic := TDAABB.CrearABB[string, string](funcionCompararBasicaStrings)
	clave := ""
	dic.Guardar(clave, clave)
	require.True(t, dic.Pertenece(clave))
	require.EqualValues(t, 1, dic.Cantidad())
	require.EqualValues(t, clave, dic.Obtener(clave))
}

func TestReemplazoDato(t *testing.T) {
	t.Log("Guarda un par de claves, y luego vuelve a guardar, buscando que el dato se haya reemplazado")
	clave := "Gato"
	clave2 := "Perro"
	dic := TDAABB.CrearABB[string, string](funcionCompararBasicaStrings)
	dic.Guardar(clave, "miau")
	dic.Guardar(clave2, "guau")
	require.True(t, dic.Pertenece(clave))
	require.True(t, dic.Pertenece(clave2))
	require.EqualValues(t, "miau", dic.Obtener(clave))
	require.EqualValues(t, "guau", dic.Obtener(clave2))
	require.EqualValues(t, 2, dic.Cantidad())

	dic.Guardar(clave, "miu")
	require.EqualValues(t, 2, dic.Cantidad())
	dic.Guardar(clave2, "baubau")
	require.EqualValues(t, 2, dic.Cantidad())
	require.True(t, dic.Pertenece(clave))
	require.True(t, dic.Pertenece(clave2))
	require.EqualValues(t, "miu", dic.Obtener(clave))
	require.EqualValues(t, "baubau", dic.Obtener(clave2))
	require.EqualValues(t, 2, dic.Cantidad())
}

func TestABBBorrar(t *testing.T) {
	t.Log("Guarda algunos pocos elementos en el ABB, y se los borra, revisando que en todo momento " +
		"el ABB se comporte de manera adecuada")
	clave1 := "Gato"
	clave2 := "Perro"
	clave3 := "Vaca"
	valor1 := "miau"
	valor2 := "guau"
	valor3 := "moo"
	claves := []string{clave1, clave2, clave3}
	valores := []string{valor1, valor2, valor3}
	dic := TDAABB.CrearABB[string, string](funcionCompararBasicaStrings)

	require.False(t, dic.Pertenece(claves[0]))
	require.False(t, dic.Pertenece(claves[0]))
	dic.Guardar(claves[0], valores[0])
	dic.Guardar(claves[1], valores[1])
	dic.Guardar(claves[2], valores[2])

	require.True(t, dic.Pertenece(claves[2]))
	require.EqualValues(t, valores[2], dic.Borrar(claves[2]))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar(claves[2]) })
	require.EqualValues(t, 2, dic.Cantidad())
	require.False(t, dic.Pertenece(claves[2]))

	require.True(t, dic.Pertenece(claves[0]))
	require.EqualValues(t, valores[0], dic.Borrar(claves[0]))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar(claves[0]) })
	require.EqualValues(t, 1, dic.Cantidad())
	require.False(t, dic.Pertenece(claves[0]))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Obtener(claves[0]) })

	require.True(t, dic.Pertenece(claves[1]))
	require.EqualValues(t, valores[1], dic.Borrar(claves[1]))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar(claves[1]) })
	require.EqualValues(t, 0, dic.Cantidad())
	require.False(t, dic.Pertenece(claves[1]))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Obtener(claves[1]) })
}

func TestReutlizacionDeBorrados(t *testing.T) {
	t.Log("Prueba de caja blanca: revisa, para el caso que fuere un HashCerrado, que no haya problema " +
		"reinsertando un elemento borrado")
	dic := TDAABB.CrearABB[string, string](funcionCompararBasicaStrings)
	clave := "hola"
	dic.Guardar(clave, "mundo!")
	dic.Borrar(clave)
	require.EqualValues(t, 0, dic.Cantidad())
	require.False(t, dic.Pertenece(clave))
	dic.Guardar(clave, "mundooo!")
	require.True(t, dic.Pertenece(clave))
	require.EqualValues(t, 1, dic.Cantidad())
	require.EqualValues(t, "mundooo!", dic.Obtener(clave))
}

func TestConClavesNumericas(t *testing.T) {
	t.Log("Valida que no solo funcione con strings")
	dic := TDAABB.CrearABB[int, string](funcionCompararBasicaInts)
	clave := 10
	valor := "Gatito"

	dic.Guardar(clave, valor)
	require.EqualValues(t, 1, dic.Cantidad())
	require.True(t, dic.Pertenece(clave))
	require.EqualValues(t, valor, dic.Obtener(clave))
	require.EqualValues(t, valor, dic.Borrar(clave))
	require.False(t, dic.Pertenece(clave))
}

func TestConClavesStructs(t *testing.T) {
	t.Log("Valida que tambien funcione con estructuras mas complejas")
	type basico struct {
		a string
		b int
	}
	type avanzado struct {
		w int
		x basico
		y basico
		z string
	}

	dic := TDAABB.CrearABB[avanzado, int](func(elem1, elem2 avanzado) int {
		return funcionCompararBasicaStrings(elem1.z, elem2.z)
	})

	a1 := avanzado{w: 10, z: "hola", x: basico{a: "mundo", b: 8}, y: basico{a: "!", b: 10}}
	a2 := avanzado{w: 10, z: "aloh", x: basico{a: "odnum", b: 14}, y: basico{a: "!", b: 5}}
	a3 := avanzado{w: 10, z: "hello", x: basico{a: "world", b: 8}, y: basico{a: "!", b: 4}}

	dic.Guardar(a1, 0)
	dic.Guardar(a2, 1)
	dic.Guardar(a3, 2)

	require.True(t, dic.Pertenece(a1))
	require.True(t, dic.Pertenece(a2))
	require.True(t, dic.Pertenece(a3))
	require.EqualValues(t, 0, dic.Obtener(a1))
	require.EqualValues(t, 1, dic.Obtener(a2))
	require.EqualValues(t, 2, dic.Obtener(a3))
	dic.Guardar(a1, 5)
	require.EqualValues(t, 5, dic.Obtener(a1))
	require.EqualValues(t, 2, dic.Obtener(a3))
	require.EqualValues(t, 5, dic.Borrar(a1))
	require.False(t, dic.Pertenece(a1))
	require.EqualValues(t, 2, dic.Obtener(a3))

}

func TestValorNulo(t *testing.T) {
	t.Log("Probamos que el valor puede ser nil sin problemas")
	dic := TDAABB.CrearABB[string, *int](funcionCompararBasicaStrings)
	clave := "Pez"
	dic.Guardar(clave, nil)
	require.True(t, dic.Pertenece(clave))
	require.EqualValues(t, 1, dic.Cantidad())
	require.EqualValues(t, (*int)(nil), dic.Obtener(clave))
	require.EqualValues(t, (*int)(nil), dic.Borrar(clave))
	require.False(t, dic.Pertenece(clave))
}

func TestCadenaLargaParticular(t *testing.T) {
	t.Log("Se han visto casos problematicos al utilizar la funcion de hashing de K&R, por lo que " +
		"se agrega una prueba con dicha funcion de hashing y una cadena muy larga")
	// El caracter '~' es el de mayor valor en ASCII (126).
	claves := make([]string, 10)
	cadena := "%d~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" +
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	dic := TDAABB.CrearABB[string, string](funcionCompararBasicaStrings)
	valores := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	for i := 0; i < 10; i++ {
		claves[i] = fmt.Sprintf(cadena, i)
		dic.Guardar(claves[i], valores[i])
	}
	require.EqualValues(t, 10, dic.Cantidad())

	ok := true
	for i := 0; i < 10 && ok; i++ {
		ok = dic.Obtener(claves[i]) == valores[i]
	}

	require.True(t, ok, "Obtener clave larga funciona")
}

func TestIterarClavesNumericas(t *testing.T) {
	t.Log("Iterar con ints, y funcion comparar sencilla, seguira mientras esten ordenados")
	dic := TDAABB.CrearABB[int, string](funcionCompararBasicaInts)

	clave := 0
	for i := 5; i < 10; i++ {
		clave = rand.Intn(100)
		dic.Guardar(clave, "n"+strconv.Itoa(clave))
	}
	anterior := -1
	iterados := 0
	dic.Iterar(func(clave int, valor string) bool {
		if clave < anterior {
			return false
		}
		anterior = clave
		iterados++
		return true
	})
	require.EqualValues(t, dic.Cantidad(), iterados, " NO itero todos los elementos, no estaban ordenados")

}

func TestIterarClavesNumericasRango(t *testing.T) {
	t.Log("Iterar con ints, y funcion comparar sencilla, seguira mientras esten ordenados y en rango")
	dic := TDAABB.CrearABB[int, string](funcionCompararBasicaInts)

	for i := 2; i < 12; i++ {
		dic.Guardar(i, "n"+strconv.Itoa(i))
	}

	anterior := -1
	iterados := 0
	desde := 3
	hasta := 10
	mensaje := ""
	dic.IterarRango(&desde, &hasta, func(clave int, valor string) bool {
		if clave < desde {
			mensaje = fmt.Sprintf("elemento no esta en rango, %d < %d", clave, desde)
			return false
		}

		if clave > hasta {
			mensaje = fmt.Sprintf("elemento no esta en rango, %d < %d", hasta, clave)
			return false
		}

		if clave < anterior {
			mensaje = fmt.Sprintf("Orden incorrecto no es mayor, %d < %d", anterior, clave)
			return false
		}
		anterior = clave
		iterados++
		return true
	})

	require.EqualValues(t, "", mensaje, fmt.Sprintf(" NO itero todos los elementos, %s", mensaje))

}

func TestIterarClavesNumericasExterno(t *testing.T) {
	t.Log("Iterar externo con ints, y funcion comparar sencilla, seguira mientras esten ordenados")

	dic := TDAABB.CrearABB[int, string](funcionCompararBasicaInts)

	for i := 2; i < 12; i++ {
		dic.Guardar(i, "n"+strconv.Itoa(i))
	}

	anterior := -1
	iterados := 0

	mensaje := ""

	iterador := dic.Iterador()

	for iterador.HaySiguiente() {
		clave, valor := iterador.VerActual()
		if clave < anterior {
			mensaje = fmt.Sprintf("Orden incorrecto no es mayor, %d < %d", anterior, clave)
			break
		}
		if valor != "n"+strconv.Itoa(clave) {
			mensaje = fmt.Sprintf("valor incorrecto , en indice %d ", clave)
			break
		}
		anterior = clave
		iterador.Siguiente()
		iterados++
	}
	if mensaje == "" && iterados < dic.Cantidad() {
		mensaje = "No recorrio todos por alguna magica razon"
	}

	require.EqualValues(t, "", mensaje, fmt.Sprintf(" NO itero todos los elementos, %s", mensaje))

}

func TestIterarClavesNumericasRangoExterno(t *testing.T) {
	t.Log("Iterar externo con ints, y funcion comparar sencilla, seguira mientras esten ordenados")

	dic := TDAABB.CrearABB[int, string](funcionCompararBasicaInts)

	for i := 2; i < 12; i++ {
		dic.Guardar(i, "n"+strconv.Itoa(i))
	}

	anterior := -1
	iterados := 0
	desde := 3
	hasta := 10

	mensaje := ""

	iterador := dic.IteradorRango(&desde, &hasta)

	for iterador.HaySiguiente() {
		clave, valor := iterador.VerActual()

		if clave < desde {
			mensaje = fmt.Sprintf("elemento no esta en rango, %d < %d", clave, desde)
			break
		}

		if clave > hasta {
			mensaje = fmt.Sprintf("elemento no esta en rango, %d < %d", hasta, clave)
			break
		}

		if clave < anterior {
			mensaje = fmt.Sprintf("Orden incorrecto no es mayor, %d < %d", anterior, clave)
			break
		}
		if valor != "n"+strconv.Itoa(clave) {
			mensaje = fmt.Sprintf("valor incorrecto , en indice %d ", clave)
			break
		}
		anterior = clave
		iterador.Siguiente()
		iterados++
	}
	if mensaje == "" && iterados < 7 {
		mensaje = "No recorrio todos por alguna magica razon"
	}

	require.EqualValues(t, "", mensaje, fmt.Sprintf(" NO itero todos los elementos, %s", mensaje))

}

func llamadosBinariosRango(minimo int, maximo int, visitar func(int)) {
	if minimo+1 >= maximo {
		return
	}
	medio := (minimo + maximo) >> 1
	visitar(medio)
	llamadosBinariosRango(minimo, medio, visitar)
	llamadosBinariosRango(medio, maximo, visitar)
}

func TestVolumen(t *testing.T) {
	dic := TDAABB.CrearABB[int, string](funcionCompararBasicaInts)
	tamanio := 200000
	t.Log(fmt.Sprintf("Insertando 200 000 elementos desordenados y verificando hayan sido agregados y se recorran en orden, despues se borra la mitad y se verificar se itere en orden externamente por rango de 100 200 a 200 256, verificando no itere de mas"))
	llamadosBinariosRango(0, tamanio, func(indice int) {
		dic.Guardar(indice, "n"+strconv.Itoa(indice))
	})

	require.EqualValues(t, tamanio-1, dic.Cantidad(), "No agrego todos los elementos")
	i := 1
	dic.IterarRango(nil, nil, func(clave int, valor string) bool {
		require.EqualValues(t, i, clave, "Clave no estuvo en orden")
		require.EqualValues(t, valor, "n"+strconv.Itoa(i), "Valor no fue correcto")
		i++
		return true
	})

	require.EqualValues(t, tamanio, i, "No itero todos los elementos internamente")
	i = 1
	for i < tamanio>>1 {
		dic.Borrar(i)
		i++
	}

	require.EqualValues(t, tamanio-tamanio>>1, dic.Cantidad(), "No borro una cantidad correcta de elementos")

	minimo := tamanio - tamanio>>1 + 200
	maximo := tamanio + 256
	iterador := dic.IteradorRango(&minimo, &maximo)
	i = minimo
	for iterador.HaySiguiente() {
		require.EqualValues(t, i, iterador.Siguiente(), "Valor incorrecto en iterador externo por rango.")
		i++
	}

	require.EqualValues(t, tamanio, i, "No recorrio una cantidad correcta de elementos iterador externo por rango")

}

func TestCasoBorde(t *testing.T) {
	t.Log("Se buscara por nodos fantasmas, repetidos, en el caso borrar cuando se hace el ingresado tal que el abb hace al borrar el swap con el mayor de los menores, o menor de los mayores. Y ademas ese hijo mayor de los menores tiene un menor")
	dic := TDAABB.CrearABB[int, string](funcionCompararBasicaInts)
	dic.Guardar(0, "U")
	dic.Guardar(1, "U")
	dic.Guardar(4, "U")
	dic.Guardar(5, "U")
	dic.Guardar(3, "U")
	dic.Guardar(2, "U")

	dic.Borrar(4)
	require.False(t, dic.Pertenece(4))
	require.True(t, dic.Pertenece(2))

	repetidos := make([]bool, 6)
	strRepetidos := " "
	dic.Iterar(func(clave int, valor string) bool {
		if repetidos[clave] {
			strRepetidos += strconv.Itoa(clave) + " "
		}
		repetidos[clave] = true
		return true
	})
	require.EqualValues(t, 5, dic.Cantidad(), fmt.Sprintf("No borro correctamente"))
	require.EqualValues(t, " ", strRepetidos, "Se encontro esos repetidos... ")

}
