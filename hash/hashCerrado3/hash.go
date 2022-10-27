package diccionario

import "fmt"
import "hash/xxh3"

// constants
type status = int

const (
	_CAPACIDAD_INICIAL            = 4
	_MAXIMA_CARGA                 = 80 // esta constante tendria unidad de 10%, osea 4 = 40%
	_MINIMA_CARGA                 = 5  // esta constante tendria unidad de 10%, osea 1 = 10%
	ERROR_NO_ESTABA               = "La clave no pertenece al diccionario"
	ERROR_ITERADOR_TERMINO        = "El iterador termino de iterar"
	_VACIO                 status = 0
	_BORRADO               status = -1
	_OCUPADO               status = 1
)

// utilities

func toBytes(objeto interface{}) []byte {
	
	str, esString := objeto.(string)

	if esString{
		return []byte(str)
	}

	return []byte(fmt.Sprintf("%v", objeto))
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

func hasheame(bytes []byte, max int) int {
	return int(xxh3.Hash(bytes) % uint64(max))
}

func aplicaFuncionDeHash[K comparable](clave K, maximo int) int { // paso intermedio para hacer mas facil cambios
	return hasheame(toBytes(clave), maximo)
	//return _JenkinsHashFunction(toBytes(clave)) % maximo
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

func (hash *hashCerrado[K, V]) iterarInterno(visitar func(*elementoCerrado[K, V]) bool) {
	i := 0
	for i < len(hash.elementos) && (hash.elementos[i].estado != _OCUPADO || visitar(&hash.elementos[i])) {
		i++
	}
}

// las unidades serian de 10%, aca lo que se haria es comparar el 100% de la cantidad con el (_MAXIMA_CARGA*10)% de la longitud(40%)
func (hash *hashCerrado[K, V]) superoCargaPermitida() bool {
	return 100*(hash.cantidad+hash.borrados) >= len(hash.elementos)*_MAXIMA_CARGA
}

func (hash *hashCerrado[K, V]) ocupaMuchaMemoria() bool {
	return len(hash.elementos) >= 2*_CAPACIDAD_INICIAL && 100*hash.cantidad <= len(hash.elementos)*_MINIMA_CARGA
}

func deberiaSeguir[K comparable, V any](elemento *elementoCerrado[K, V], clave K) bool {
	return elemento.estado != _VACIO && elemento.clave != clave
}
func buscarElementoCerrado[K comparable, V any](elementos []elementoCerrado[K, V], clave K, haceAlgo func(*elementoCerrado[K, V])) *elementoCerrado[K, V] {

	posInicial := aplicaFuncionDeHash(clave, len(elementos))
	i := posInicial

	for i < len(elementos) && deberiaSeguir(&elementos[i], clave) {
		haceAlgo(&elementos[i])
		i++
	}
	if i < len(elementos) {
		return &elementos[i]
	}

	i = 0
	for i < posInicial && deberiaSeguir(&elementos[i], clave) {
		haceAlgo(&elementos[i])
		i++
	}

	if i == posInicial {
		return nil // significaria recorrio todo sin exito
	}

	return &elementos[i]

}

func buscarElementoAModificar[K comparable, V any](elementos []elementoCerrado[K, V], clave K) *elementoCerrado[K, V] {
	var aDevolver *elementoCerrado[K, V] = nil
	ultimoVisto := buscarElementoCerrado(elementos, clave, func(elemento *elementoCerrado[K, V]) {
		if aDevolver == nil && elemento.estado == _BORRADO {
			aDevolver = elemento
		}
	})

	if ultimoVisto != nil && (aDevolver == nil || ultimoVisto.estado != _VACIO) {
		aDevolver = ultimoVisto
	}

	return aDevolver
}

func (hash *hashCerrado[K, V]) redimensionar(nuevoLargo int) {
	nuevos := crearTabla[K, V](nuevoLargo)
	hash.borrados = 0
	hash.iterarInterno(func(elemento *elementoCerrado[K, V]) bool {
		*buscarElementoAModificar(nuevos, elemento.clave) = *elemento
		return true
	})

	hash.elementos = nuevos

}

func (hash *hashCerrado[K, V]) Guardar(clave K, valor V) {
	if hash.superoCargaPermitida() {
		hash.redimensionar(2 * len(hash.elementos))
	}

	aModificar := buscarElementoAModificar(hash.elementos, clave)
	if aModificar.estado != _OCUPADO {
		hash.cantidad++
		if aModificar.estado == _BORRADO {
			hash.borrados--
		}
	} else {
		fmt.Printf("\nESTABA OCUPADO??? %v::: %v\n ", clave, *aModificar)
	}

	*aModificar = crearElementoCerrado(clave, valor)
}

func (hash *hashCerrado[K, V]) esValido(elem *elementoCerrado[K, V]) bool {
	return elem != nil && elem.estado == _OCUPADO
}

func (hash *hashCerrado[K, V]) panicNoEstaba(elemento *elementoCerrado[K, V]) *elementoCerrado[K, V] {
	if !hash.esValido(elemento) {
		panic(ERROR_NO_ESTABA)
	}

	return elemento
}

func (hash *hashCerrado[K, V]) Pertenece(clave K) bool {
	return hash.esValido(buscarElementoCerrado(hash.elementos, clave, func(elemento *elementoCerrado[K, V]) {}))
}

func (hash *hashCerrado[K, V]) Obtener(clave K) V {
	return hash.panicNoEstaba(buscarElementoCerrado(hash.elementos, clave, func(elemento *elementoCerrado[K, V]) {})).valor
}

func (hash *hashCerrado[K, V]) Borrar(clave K) V {
	elemento := hash.panicNoEstaba(buscarElementoCerrado(hash.elementos, clave, func(elemento *elementoCerrado[K, V]) {}))
	elem := elemento.valor
	hash.cantidad--
	*elemento = crearElementoCerradoVacio[K, V]()
	elemento.estado = _BORRADO
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
	hash.iterarInterno(func(elemento *elementoCerrado[K, V]) bool { return visitar(elemento.clave, elemento.valor) })
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
