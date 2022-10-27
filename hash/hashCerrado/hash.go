package diccionario

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

// *2 redimensionar
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

func aplicaFuncionDeHash[K comparable](clave K, maximo int) int { // paso intermedio para hacer mas facil cambios
	return _JenkinsHashFunction(toBytes(clave)) % maximo
}

type elementoCerrado[K comparable, V any] struct {
	clave  K
	valor  V
	estado status
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
	return len(hash.elementos) >= _CAPACIDAD_INICIAL<<1 && 100*hash.cantidad <= len(hash.elementos)*_MINIMA_CARGA
}

func deberiaSeguir[K comparable, V any](elemento *elementoCerrado[K, V], clave K) bool {
	return elemento.estado != _VACIO && elemento.clave != clave
}

// se podrian haber usado enums en vez de bool para ahorrar un if en guardar.
// tambien se probo devolviendo punteros, y era mas lento.
func buscarPosicionElementoCerrado[K comparable, V any](elementos []elementoCerrado[K, V], clave K) (int, bool) {

	indiceRef := -1

	posInicial := aplicaFuncionDeHash(clave, len(elementos))

	i := posInicial

	for i < len(elementos) && deberiaSeguir(&elementos[i], clave) {
		if indiceRef == -1 && elementos[i].estado == _BORRADO { // se agarra el primer borrado por defecto
			indiceRef = i
		}
		i++
	}
	if i < len(elementos) { // significaria deberiaSeguir fue false, es decir vacio o igual clave.
		return i, elementos[i].estado == _OCUPADO
	}

	i = 0

	for i < posInicial && deberiaSeguir(&elementos[i], clave) {
		if indiceRef == -1 && elementos[i].estado == _BORRADO { // se agarra el primer borrado por defecto
			indiceRef = i
		}
		i++
	}

	if i < posInicial { // significaria deberiaSeguir fue false, es decir vacio o igual clave.
		return i, elementos[i].estado == _OCUPADO
	}

	return indiceRef, false // no se encontro, se devuelve la referencia del borrado
}

func (hash *hashCerrado[K, V]) redimensionar(nuevoLargo int) {
	nuevos := crearTabla[K, V](nuevoLargo)
	hash.borrados = 0

	hash.Iterar(func(clave K, valor V) bool { // se podria usar punteros, pero no se vio especial mejora
		indice, _ := buscarPosicionElementoCerrado(nuevos, clave)
		nuevos[indice] = crearElementoCerrado(clave, valor)
		return true
	})

	hash.elementos = nuevos

}

func (hash *hashCerrado[K, V]) Guardar(clave K, valor V) {
	if hash.superoCargaPermitida() {
		hash.redimensionar(len(hash.elementos)<<2)
	}
	indice, estaba := buscarPosicionElementoCerrado(hash.elementos, clave)
	if !estaba {
		hash.cantidad++
		if hash.elementos[indice].estado == _BORRADO {
			hash.borrados--
		}
	}

	hash.elementos[indice] = crearElementoCerrado(clave, valor)
}

// da el elemento si pertenece, sino panic
func (hash *hashCerrado[K, V]) dameElemento(ind int, estaba bool) *elementoCerrado[K, V] {
	if !estaba {
		panic(ERROR_NO_ESTABA)
	}

	return &hash.elementos[ind]
}

func (hash *hashCerrado[K, V]) Pertenece(clave K) bool {
	_, estaba := buscarPosicionElementoCerrado(hash.elementos, clave) // no me gusta tener que guardarlo pero bueno, para reutilizar directamente
	return estaba
}

func (hash *hashCerrado[K, V]) Obtener(clave K) V {
	return hash.dameElemento(buscarPosicionElementoCerrado(hash.elementos, clave)).valor
}

func (hash *hashCerrado[K, V]) Borrar(clave K) V {
	elemento := hash.dameElemento(buscarPosicionElementoCerrado(hash.elementos, clave))
	elem := elemento.valor
	hash.cantidad--
	*elemento = crearElementoCerradoVacio[K, V]()
	elemento.estado = _BORRADO
	hash.borrados++

	if hash.ocupaMuchaMemoria() {
		hash.redimensionar(len(hash.elementos)>>1)
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
