package hash_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	HashAbierto "hash/hashCerrado"
	Hash "hash/hashCerrado2"
	HashCuckoo "hash/hashCuckoo2"
	"testing"
	"time"
)

type Diccionario[K comparable, V any] interface {

	// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
	Guardar(clave K, dato V)

	// Pertenece determina si una clave ya se encuentra en el diccionario, o no
	Pertenece(clave K) bool

	// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
	// 'La clave no pertenece al diccionario'
	Obtener(clave K) V

	// Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
	// pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
	Borrar(clave K) V

	// Cantidad devuelve la cantidad de elementos dentro del diccionario
	Cantidad() int

	// Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del
	// mismo
	Iterar(func(clave K, dato V) bool)
}

const ERROR_FUNCION_HASH = "Error: mala funcion de hash, redimension requeridad demasiadas veces seguidas"
const ERROR_NO_ESTABA = "La clave no pertenece al diccionario"
const ERROR_ITERADOR_TERMINO = "El iterador termino de iterar"

func ejecutarPruebaVolumen(b *testing.T, n int) {
	init := time.Now()
	dic := Hash.CrearHash[string, int]()

	claves := make([]string, n)
	valores := make([]int, n)

	/* Inserta 'n' parejas en el hash */
	for i := 0; i < n; i++ {
		valores[i] = i
		claves[i] = fmt.Sprintf("%20d", i)
		dic.Guardar(claves[i], valores[i])
	}

	//require.EqualValues(b, n, dic.Cantidad(), "La cantidad de elementos es incorrecta")

	/* Verifica que devuelva los valores correctos */
	ok := true
	for i := 0; i < n; i++ {
		ok = dic.Pertenece(claves[i])
		if !ok {
			break
		}
		ok = dic.Obtener(claves[i]) == valores[i]
		if !ok {
			break
		}
	}

	//require.True(b, ok, "Pertenece y Obtener con muchos elementos no funciona correctamente")
	//require.EqualValues(b, n, dic.Cantidad(), "La cantidad de elementos es incorrecta")

	/* Verifica que borre y devuelva los valores correctos */
	for i := 0; i < n; i++ {
		ok = dic.Borrar(claves[i]) == valores[i]
		if !ok {
			break
		}
	}

	b.Log(fmt.Sprintf("Pruebas con %d elementos Tomo %s", n, time.Since(init)))
	//require.True(b, ok, "Borrar muchos elementos no funciona correctamente")
	//require.EqualValues(b, 0, dic.Cantidad())
}

func ejecutarPruebaVolumenGenerico[T Diccionario[string, int]](b *testing.T, dic T, n int) {
	defer func() {
		if r := recover(); r != nil {
			b.Log("ERROR :: ", r)
		}
	}()

	claves := make([]string, n)
	valores := make([]int, n)

	/* Inserta 'n' parejas en el hash */
	for i := 0; i < n; i++ {
		valores[i] = i
		claves[i] = fmt.Sprintf("%08d", i)
		dic.Guardar(claves[i], valores[i])
	}

	//require.EqualValues(b, n, dic.Cantidad(), "La cantidad de elementos es incorrecta")

	/* Verifica que devuelva los valores correctos */
	ok := true
	for i := 0; i < n; i++ {
		ok = dic.Pertenece(claves[i])
		if !ok {
			break
		}
		ok = dic.Obtener(claves[i]) == valores[i]
		if !ok {
			break
		}
	}

	//require.True(b, ok, "Pertenece y Obtener con muchos elementos no funciona correctamente")
	//require.EqualValues(b, n, dic.Cantidad(), "La cantidad de elementos es incorrecta")

	/* Verifica que borre y devuelva los valores correctos */
	for i := 0; i < n; i++ {
		ok = dic.Borrar(claves[i]) == valores[i]
		if !ok {
			break
		}
	}

	//require.True(b, ok, "Borrar muchos elementos no funciona correctamente")
	//require.EqualValues(b, 0, dic.Cantidad())
}

func ejecutarPruebasVolumenIterador(b *testing.T, n int) {
	init := time.Now()
	dic := Hash.CrearHash[string, *int]()

	claves := make([]string, n)
	valores := make([]int, n)

	/* Inserta 'n' parejas en el hash */
	for i := 0; i < n; i++ {
		claves[i] = fmt.Sprintf("%08d", i)
		valores[i] = i
		dic.Guardar(claves[i], &valores[i])
	}

	// Prueba de iteración sobre las claves almacenadas.
	iter := dic.Iterador()
	require.True(b, iter.HaySiguiente())

	ok := true
	var i int
	var clave string
	var valor *int

	for i = 0; i < n; i++ {
		if !iter.HaySiguiente() {
			ok = false
			break
		}
		c1, v1 := iter.VerActual()
		clave = c1
		if clave == "" {
			ok = false
			break
		}
		valor = v1
		if valor == nil {
			ok = false
			break
		}
		*valor = n
		iter.Siguiente()
	}
	require.True(b, ok, "Iteracion en volumen no funciona correctamente")
	require.EqualValues(b, n, i, "No se recorrió todo el largo")
	require.False(b, iter.HaySiguiente(), "El iterador debe estar al final luego de recorrer")

	ok = true
	for i = 0; i < n; i++ {
		if valores[i] != n {
			ok = false
			break
		}
	}
	b.Log(fmt.Sprintf("Pruebas iteracion con %d elementos Tomo %s", n, time.Since(init)))
	require.True(b, ok, "No se cambiaron todos los elementos")
}

func TestAgregado(t *testing.T) {
	hash := Hash.CrearHash[string, int]()

	require.NotNil(t, hash, "Crear devolvio nil")

	hash.Guardar("Uno", 1)
	hash.Guardar("Dos", 2)
	hash.Guardar("Tres", 3)

	require.EqualValues(t, 3, hash.Cantidad(), "Cantidad incorrecta")
}

func TestPertenece(t *testing.T) {
	hash := Hash.CrearHash[string, int]()

	require.NotNil(t, hash)

	hash.Guardar("Uno", 1)
	hash.Guardar("Dos", 2)
	hash.Guardar("Tres", 3)

	require.EqualValues(t, 3, hash.Cantidad(), "Cantidad incorrecta")

	require.False(t, hash.Pertenece("O"), "Se invento una clave 'O'")
	require.False(t, hash.Pertenece(""), "Se invento una clave ''")
	require.False(t, hash.Pertenece("Makise"), "Se invento una clave 'Makise'")

	require.True(t, hash.Pertenece("Uno"), "Se olvido de la clave 'Uno'")
	require.True(t, hash.Pertenece("Dos"), "Se olvido de la clave 'Dos'")
	require.True(t, hash.Pertenece("Tres"), "Se olvido de la clave 'Tres'")
}

func TestObtener(t *testing.T) {
	hash := Hash.CrearHash[string, int]()

	require.NotNil(t, hash)

	hash.Guardar("Uno", 1)
	hash.Guardar("Dos", 2)
	hash.Guardar("Tres", 3)

	require.EqualValues(t, 3, hash.Cantidad(), "Cantidad incorrecta")

	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Obtener("Makise") })
	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Obtener("Miu") })
	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Obtener("Inaba") })
	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Obtener("Zero Two") })

	require.EqualValues(t, 1, hash.Obtener("Uno"))
	require.EqualValues(t, 2, hash.Obtener("Dos"))
	require.EqualValues(t, 3, hash.Obtener("Tres"))
}

func TestBorrar(t *testing.T) {
	hash := Hash.CrearHash[string, int]()

	require.NotNil(t, hash)

	hash.Guardar("Uno", 1)
	hash.Guardar("Dos", 2)
	hash.Guardar("Tres", 3)
	hash.Guardar("Makise", 10)
	hash.Guardar("Inaba", 9)

	require.EqualValues(t, 5, hash.Cantidad(), "Cantidad incorrecta")

	require.EqualValues(t, 10, hash.Borrar("Makise"))
	require.EqualValues(t, 1, hash.Borrar("Uno"))

	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Borrar("Makise") })
	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Borrar("Misaka") })
	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Obtener("Makise") })
	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Obtener("Uno") })

	require.EqualValues(t, 9, hash.Obtener("Inaba"))
	require.EqualValues(t, 2, hash.Obtener("Dos"))
	require.EqualValues(t, 3, hash.Obtener("Tres"))
	require.EqualValues(t, 3, hash.Cantidad(), "Cantidad incorrecta tras borrar")
}

func TestIteradorExterno(t *testing.T) {
	hash := Hash.CrearHash[string, int]()

	require.NotNil(t, hash)

	hash.Guardar("Uno", 1)
	hash.Guardar("Misaka", 2)
	hash.Guardar("Miu", 3)
	hash.Guardar("Makise", 10)
	hash.Guardar("Inaba", 9)

	require.EqualValues(t, 5, hash.Cantidad(), "Cantidad incorrecta")

	require.EqualValues(t, 10, hash.Borrar("Makise"))
	require.EqualValues(t, 1, hash.Borrar("Uno"))

	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Borrar("Makise") })
	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Obtener("Makise") })
	require.PanicsWithValue(t, ERROR_NO_ESTABA, func() { hash.Obtener("Uno") })

	require.EqualValues(t, 9, hash.Obtener("Inaba"))
	require.EqualValues(t, 3, hash.Obtener("Miu"))
	require.EqualValues(t, 2, hash.Obtener("Misaka"))
	require.EqualValues(t, 3, hash.Cantidad(), "Cantidad incorrecta tras borrar")

	iterador := hash.Iterador()

	i := 0
	for iterador.HaySiguiente() && i < 10 {
		i++
		t.Log("ITERADOR VIO == " + iterador.Siguiente())
	}

	require.EqualValues(t, 3, i, "No mostro la cantidad correcta el iterador externo?")

}

func TestVolumen(t *testing.T) {
	n := 400000
	var iteraciones int64 = 20
	var milisecond int64 = 1000000 //000
	init := time.Now()
	var i int64 = 0
	for i < iteraciones {
		ejecutarPruebaVolumenGenerico(t, Hash.CrearHash[string, int](), n)
		i++
	}
	t.Log(fmt.Sprintf("Tomo en promedio %dms, pruebas con hash cerrado 2, %d elementos %d veces", (int64(time.Since(init)) / (milisecond * iteraciones)), n, iteraciones))

	init = time.Now()
	i = 0
	for i < iteraciones {
		ejecutarPruebaVolumenGenerico(t, HashAbierto.CrearHash[string, int](), n)
		i++
	}
	t.Log(fmt.Sprintf("Tomo en promedio %dms, pruebas con hash cerrado , %d elementos %d veces", (int64(time.Since(init)) / (milisecond * iteraciones)), n, iteraciones))

	init = time.Now()
	ejecutarPruebaVolumenGenerico(t, HashCuckoo.CrearHash[string, int](), n)
	t.Log(fmt.Sprintf("Pruebas con hash cuckoo, %d elementos Tomo %s", n, time.Since(init)))

	//ejecutarPruebasVolumenIterador(t,n)
}

func BenchmarkIterador(b *testing.B) {
	b.Log("ESTA HACIENDO EL TEST?")
	n := 12500

	b.Run(fmt.Sprintf("Prueba %d elementos", n+8), func(b *testing.B) {
		hash := Hash.CrearHash[string, int]()
		hash.Guardar("Uno", 1)
		hash.Guardar("Dos", 2)
		hash.Guardar("Tres", 3)
		hash.Guardar("Makise", 10)
		hash.Guardar("Inaba", 9)
		hash.Guardar("Misaka", 9)
		hash.Guardar("Miu", 9)
		hash.Guardar("Marcelo", 9)
		i := 0
		keys := make([]string, n)
		values := make([]int, n)
		defer func() {
			if r := recover(); r != nil {
				b.Log(fmt.Sprintf("Err:%s ... key = %s, value = %d", r, keys[i], values[i]))
			}
		}()

		for i < n {
			keys[i] = fmt.Sprintf("%08d", i)
			values[i] = i + 1
			hash.Guardar(keys[i], i)
			hash.Guardar(keys[i], i+1)
			i++
		}

		i = 0
		res := 0
		for i < n {
			if !hash.Pertenece(keys[i]) {
				panic("NO PERTENECIA???")
			}
			res = hash.Obtener(keys[i])
			if res != values[i] {
				panic(fmt.Sprintf("Que te inventas???? no es %d", res))
			}

			i++
		}

	})
}
