package diccionario

import "reflect"
import "fmt"

// constants
type status = int

const (
	_CAPACIDAD_INICIAL            = 128
	_MAXIMA_CARGA                 = 80 // esta constante tendria unidad de 1%, osea 80 = 80%
	_MINIMA_CARGA                 = 5  // esta constante tendria unidad de 1%, osea 5 = 5%
	ERROR_NO_ESTABA               = "La clave no pertenece al diccionario"
	ERROR_ITERADOR_TERMINO        = "El iterador termino de iterar"
	_VACIO                 status = 0
	_BORRADO               status = -1
	_OCUPADO               status = 1
)

// utilities

func toBytes(objeto interface{}) []byte {
	switch objeto.(type) {
	case string: // se chequea el tipo para saber cuando se puede usar una forma mas rapida
		return []byte(reflect.ValueOf(objeto).String())
	default:
		return []byte(fmt.Sprintf("%v", objeto)) // lento pero justo
	}
}

func _JenkinsHashFunction(bytes []byte) int {
	res := 0
	for i := 0; i < len(bytes); i++ {
		res += int(bytes[i])
		res += res << 10
		res ^= res >> 6
	}

	return res
}

func aplicaFuncionDeHash[K comparable](clave K, maximo int) int { // paso intermedio para hacer mas facil cambios
	return _JenkinsHashFunction(toBytes(clave)) % maximo
}

type elementoCerrado[K comparable, V any] struct {
	clave  K
	valor  V
	estado status
}

func (elemento *elementoCerrado[K, V]) modificar(clave K, valor V, estado status) {
	elemento.clave = clave
	elemento.valor = valor
	elemento.estado = estado
}

func crearElementoCerrado[K comparable, V any](clave K, valor V) elementoCerrado[K, V] {
	return elementoCerrado[K, V]{clave, valor, 1}
}

func crearElementoCerradoVacio[K comparable, V any]() elementoCerrado[K, V] {
	return elementoCerrado[K, V]{}
}

func crearTabla[K comparable, V any](largo int) []elementoCerrado[K, V] {
	tabla := make([]elementoCerrado[K, V], largo)
	for i := range tabla {
		tabla[i] = crearElementoCerradoVacio[K, V]()
	}
	return tabla
}

type hashCerrado[K comparable, V any] struct {
	elementos []elementoCerrado[K, V]
	cantidad  int
	borrados  int
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	hash := new(hashCerrado[K, V])
	hash.elementos = crearTabla[K, V](_CAPACIDAD_INICIAL)
	return hash
}

// las unidades serian de 1%, aca lo que se haria es comparar el 100% de la cantidad con el (_MAXIMA_CARGA)% de la longitud(80%)
func (hash *hashCerrado[K, V]) superoCargaPermitida() bool {
	return 100*(hash.cantidad+hash.borrados) >= len(hash.elementos)*_MAXIMA_CARGA
}

func (hash *hashCerrado[K, V]) ocupaMuchaMemoria() bool {
	return len(hash.elementos) >= 2*_CAPACIDAD_INICIAL && 100*hash.cantidad <= len(hash.elementos)*_MINIMA_CARGA
}

func iterarPosicionCerrado[K comparable, V any](elementos []elementoCerrado[K, V], clave K, continuar func(int) bool) {
	posInicial := aplicaFuncionDeHash(clave, len(elementos))

	i := posInicial
	for i < len(elementos) && continuar(i) {
		i++
	}
	if i < len(elementos) {
		return
	}

	i = 0
	for i < posInicial && continuar(i) {
		i++
	}

}

func buscarElementoAModificar[K comparable, V any](elementos []elementoCerrado[K, V], clave K) int {
	indiceRes := -1
	iterarPosicionCerrado(elementos, clave, func(indice int) bool {
		if elementos[indice].estado != _OCUPADO || elementos[indice].clave == clave {
			indiceRes = indice
			return false
		}
		return true
	})
	return indiceRes
}

func (hash *hashCerrado[K, V]) buscarPosicion(clave K) int {
	indiceRes := -1
	iterarPosicionCerrado(hash.elementos, clave, func(indice int) bool {
		if hash.elementos[indice].clave == clave {
			indiceRes = indice
			return false
		}
		return hash.elementos[indice].estado != _VACIO
	})
	return indiceRes
}

func (hash *hashCerrado[K, V]) redimensionar(nuevoLargo int) {
	nuevos := crearTabla[K, V](nuevoLargo)
	hash.borrados = 0

	hash.Iterar(func(clave K, valor V) bool {
		nuevos[buscarElementoAModificar(nuevos, clave)] = crearElementoCerrado(clave, valor)
		return true
	})

	hash.elementos = nuevos

}

func (hash *hashCerrado[K, V]) Guardar(clave K, valor V) {
	if hash.superoCargaPermitida() {
		hash.redimensionar(4 * len(hash.elementos))
	}
	indice := buscarElementoAModificar(hash.elementos, clave)

	if hash.elementos[indice].estado != _OCUPADO {
		hash.cantidad++
		if hash.elementos[indice].estado == _BORRADO {
			hash.borrados--
		}
	}

	hash.elementos[indice] = crearElementoCerrado(clave, valor)
}

func (hash *hashCerrado[K, V]) Pertenece(clave K) bool {
	return hash.buscarPosicion(clave) != -1
}

func (hash *hashCerrado[K, V]) Obtener(clave K) V {
	indice := hash.buscarPosicion(clave)
	if indice == -1 {
		panic(ERROR_NO_ESTABA)
	}
	return hash.elementos[indice].valor
}

func (hash *hashCerrado[K, V]) Borrar(clave K) V {
	indice := hash.buscarPosicion(clave)

	if indice == -1 {
		panic(ERROR_NO_ESTABA)
	}

	elem := hash.elementos[indice].valor
	hash.cantidad--
	hash.elementos[indice] = crearElementoCerradoVacio[K, V]()
	hash.elementos[indice].estado = _BORRADO
	hash.borrados++

	if hash.ocupaMuchaMemoria() {
		hash.redimensionar(len(hash.elementos) / 2)
	}

	return elem
}

func (hash *hashCerrado[K, V]) Cantidad() int {
	return hash.cantidad
}

func (hash *hashCerrado[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	i := 0
	for i < len(hash.elementos) && (hash.elementos[i].estado != _OCUPADO || visitar(hash.elementos[i].clave, hash.elementos[i].valor)) {
		i++
	}
}

// Iterador externo

func (hash *hashCerrado[K, V]) Iterador() IterDiccionario[K, V] {
	return creariteradorCerrado(hash)
}

type iteradorCerrado[K comparable, V any] struct {
	referencia *hashCerrado[K, V]
	posActual  int
}

func creariteradorCerrado[K comparable, V any](referencia *hashCerrado[K, V]) IterDiccionario[K, V] {
	iterador := new(iteradorCerrado[K, V])

	iterador.referencia = referencia
	iterador.posActual = -1
	iterador.iterarSiguiente()
	return iterador
}

func (iterador *iteradorCerrado[K, V]) iterarSiguiente() {
	iterador.posActual++
	for iterador.posActual < len(iterador.referencia.elementos) && iterador.referencia.elementos[iterador.posActual].estado != _OCUPADO {
		iterador.posActual++
	}
}
func (iterador *iteradorCerrado[K, V]) panicTermino() {
	if !iterador.HaySiguiente() {
		panic(ERROR_ITERADOR_TERMINO)
	}
}

func (iterador *iteradorCerrado[K, V]) HaySiguiente() bool {
	return iterador.posActual < len(iterador.referencia.elementos)
}

func (iterador *iteradorCerrado[K, V]) VerActual() (K, V) {
	iterador.panicTermino()
	elemento := iterador.referencia.elementos[iterador.posActual]
	return elemento.clave, elemento.valor
}

func (iterador *iteradorCerrado[K, V]) Siguiente() K {
	iterador.panicTermino()
	claveActual := iterador.referencia.elementos[iterador.posActual].clave
	iterador.iterarSiguiente()
	return claveActual

}
