package hash_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	HashAbierto "hash/hashAbierto"
	aEntregar "hash/hashCerrado"
	HashCuckoo "hash/hashCuckoo"

	Hash "hash/aEntregar"
	HashCerrado3 "hash/hashCerrado3"
	"hash/xxh3"
	"reflect"
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

func ejecutarPruebaVolumenGenericoStepped[T Diccionario[string, int]](b *testing.T, dic T, n int) (guardar, buscar, borrar int64) {
	defer func() {
		if r := recover(); r != nil {
			b.Log("ERROR :: ", r)
		}
	}()
	claves := make([]string, n)
	valores := make([]int, n)

	for i := 0; i < n; i++ {
		valores[i] = i
		claves[i] = fmt.Sprintf("%08d", i)
	}

	init := time.Now()
	/* Inserta 'n' parejas en el hash */

	for i := 0; i < n; i++ {
		dic.Guardar(claves[i], valores[i])
	}

	guardar = int64(time.Since(init))
	//b.Log(fmt.Sprintf("AFTER ADDING CANTIDAD = %d",dic.Cantidad()))
	//require.EqualValues(b, n, dic.Cantidad(), "La cantidad de elementos es incorrecta")
	init = time.Now()

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
	buscar = int64(time.Since(init))

	//require.True(b, ok, "Pertenece y Obtener con muchos elementos no funciona correctamente")
	//require.EqualValues(b, n, dic.Cantidad(), "La cantidad de elementos es incorrecta")

	init = time.Now()
	/* Verifica que borre y devuelva los valores correctos */
	for i := 0; i < n; i++ {
		ok = dic.Borrar(claves[i]) == valores[i]
		if !ok {
			break
		}
	}
	borrar = int64(time.Since(init))
	//b.Log(fmt.Sprintf("AFTER DELETING CANTIDAD = %d",dic.Cantidad()))

	//require.True(b, ok, "Borrar muchos elementos no funciona correctamente")
	//require.EqualValues(b, 0, dic.Cantidad())
	return
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

// /*
func TestAgregado(t *testing.T) {
	hash := Hash.CrearHash[string, int]()

	require.NotNil(t, hash, "Crear devolvio nil")

	hash.Guardar("Uno", 1)
	hash.Guardar("Dos", 2)
	hash.Guardar("Tres", 3)
	hash.Guardar("Cuatro", 3)
	hash.Guardar("Cinco", 3)
	hash.Guardar("Cinco", 6)

	require.EqualValues(t, 5, hash.Cantidad(), "Cantidad incorrecta")
}

func TestPertenece(t *testing.T) {
	hash := Hash.CrearHash[string, int]()

	require.NotNil(t, hash)

	hash.Guardar("Uno", 1)
	hash.Guardar("Dos", 2)
	hash.Guardar("Tres", 3)

	require.EqualValues(t, 3, hash.Cantidad(), "Cantidad incorrecta")

	require.False(t, hash.Pertenece("O"), "Se invento una clave 'O'")
	//require.False(t, hash.Pertenece(""), "Se invento una clave ''")
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

//*/

func _JenkinsHashFunction(bytes []byte) uint64 {
	var res uint64 = 0
	for i := 0; i < len(bytes); i++ {
		res += uint64(bytes[i])
		res += res << 10
		res ^= res >> 6
	}

	return res
}
func TestHashFunctions(t *testing.T) {
	testHashingFunctions(t, "xxh3 hash reflect", func(clave string) uint64 { return xxh3.Hash(toBytes2(clave)) })
	testHashingFunctions(t, "xxh3 hash fmt pure", func(clave string) uint64 { return xxh3.Hash(toBytes3(clave)) })
	testHashingFunctions(t, "xxh3 hash max speed", func(clave string) uint64 { return xxh3.Hash([]byte(clave)) })
	testHashingFunctions(t, "xxh3 hash fmt pure", func(clave string) uint64 { return _JenkinsHashFunction(toBytes3(clave)) })
	testHashingFunctions(t, "Jenkins max speed", func(clave string) uint64 { return _JenkinsHashFunction([]byte(clave)) })
	testHashingFunctions(t, "Jenkins reflect", func(clave string) uint64 { return _JenkinsHashFunction(toBytes2(clave)) })
	testHashingFunctions(t, "creatividad directo", func(clave string) uint64 { return creatividad2([]byte(clave)) })
	testHashingFunctions(t, "creatividad reflect", func(clave string) uint64 { return creatividad2(toBytes2(clave)) })
}

func toBytes3(objeto interface{}) []byte {
	switch objeto.(type) {
	case string: // se chequea el tipo para saber cuando se puede usar una forma mas rapida
		return []byte(fmt.Sprintf("%s", objeto))
	default:
		return []byte(fmt.Sprintf("%v", objeto)) // lento pero justo
	}
}
func toBytes2(objeto interface{}) []byte {
	switch objeto.(type) {
	case string: // se chequea el tipo para saber cuando se puede usar una forma mas rapida
		return []byte(reflect.ValueOf(objeto).String())
	default:
		return []byte(fmt.Sprintf("%v", objeto)) // lento pero justo
	}
}

func creatividad2(bytes []byte) uint64 {
	if len(bytes) == 0 {
		return 0
	}
	i := 0
	i2 := len(bytes) - 1
	var res uint64 = (uint64(bytes[0]<<1|bytes[i2]>>1) + uint64(bytes[i2]>>1))
	res ^= res << 6
	res ^= res >> 3
	i++
	i2--
	for i < 3 && i < i2 {
		res += (uint64(bytes[i]<<1|bytes[i2]>>1) + uint64(bytes[i2]>>1))
		res = (res ^ (res << 6)) ^ (res >> 3)
		i++
		i2--
	}

	return res
}

func testHashingFunctions(t *testing.T, label string, hashFunc func(string) uint64) {
	const millisecond int64 = 1000000
	const iteraciones int64 = 80000 * 5
	const maximo uint64 = 800000
	posiciones := make([]bool, 800000)
	colisiones := 0
	started := time.Now()
	var indice uint64 = 0
	var i int64 = 0
	for i = 0; i < iteraciones; i++ {
		indice = hashFunc(fmt.Sprintf("%08d", i)) % maximo
		if posiciones[indice] {
			colisiones++
		} else {
			posiciones[indice] = true
		}
	}

	ms := (int64(time.Since(started)) / (millisecond))
	t.Log(fmt.Sprintf("Tomo %dms, pruebas con '%s', %d claves, hubo %d colisiones", ms, label, iteraciones, colisiones))
}

func testVolumenPara(t *testing.T, tipo string, provider func() Diccionario[string, int]) {
	n := 400000
	const iteraciones int64 = 2
	const millisecond int64 = 1000000

	init := time.Now()
	var i int64 = 0
	for i < iteraciones {
		ejecutarPruebaVolumenGenerico(t, provider(), n)
		i++
	}

	ms := (int64(time.Since(init)) / (millisecond * iteraciones))
	t.Log(fmt.Sprintf("Tomo en promedio %dms, pruebas con '%s', %d elementos %d veces", ms, tipo, n, iteraciones))

}

func testVolumenSteppedPara(t *testing.T, n int, iteraciones int64, tipo string, provider func() Diccionario[string, int]) {
	const millisecond int64 = 1000000
	total_guardar, total_buscar, total_borrar := ejecutarPruebaVolumenGenericoStepped(t, provider(), n)
	var i int64 = 1
	for i < iteraciones {
		guardar, buscar, borrar := ejecutarPruebaVolumenGenericoStepped(t, provider(), n)
		total_guardar += guardar
		total_buscar += buscar
		total_borrar += borrar
		i++
	}

	total_guardar /= (millisecond * iteraciones)
	total_buscar /= (millisecond * iteraciones)
	total_borrar /= (millisecond * iteraciones)

	t.Log(fmt.Sprintf("%dms %dms %dms promedio con '%s' para guardar,buscar y borrar", total_guardar, total_buscar, total_borrar, tipo))

}

func TestVolumen(t *testing.T) {
	n := 400000
	const iteraciones int64 = 10

	t.Log(fmt.Sprintf("pruebas de %d elementos %d veces", n, iteraciones))

	testVolumenSteppedPara(t, n, iteraciones, "Hash a aEntregar",
		func() Diccionario[string, int] { return Hash.CrearHash[string, int]() })

	testVolumenSteppedPara(t, n, iteraciones, "Hash cerrado punteros",
		func() Diccionario[string, int] { return HashCerrado3.CrearHash[string, int]() })

	testVolumenSteppedPara(t, n, iteraciones, "Hash cerrado ins based",
		func() Diccionario[string, int] { return aEntregar.CrearHash[string, int]() })

	testVolumenSteppedPara(t, n, 2, "Hash abierto",
		func() Diccionario[string, int] { return HashAbierto.CrearHash[string, int]() })

	testVolumenSteppedPara(t, n, 2, "Hash cuckoo",
		func() Diccionario[string, int] { return HashCuckoo.CrearHash[string, int]() })

	//ejecutarPruebasVolumenIterador(t,n)
}
