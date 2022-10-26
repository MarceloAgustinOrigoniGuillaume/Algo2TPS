package diccionario

import "fmt"
import "reflect"
import TDALista "lista"

const _CAPACIDAD_INICIAL = 1000
const _MAXIMA_CARGA = 150 // esta constante tendria unidad de 10%, osea 9 = 70%
const _MINIMA_CARGA = 10  // esta constante tendria unidad de 1%, osea 10 = 10%
const ERROR_NO_ESTABA = "La clave no pertenece al diccionario"
const ERROR_ITERADOR_TERMINO = "El iterador termino de iterar"

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

func aplicarFuncionHash[K comparable](clave K, maximo int) int {
	return _JenkinsHashFunction(toBytes(clave)) % maximo
}

type elementoAbierto[K comparable, V any] struct {
	clave K
	valor V
}

func crearElementoAbierto[K comparable, V any](clave K, valor V) *elementoAbierto[K, V] {

	return &elementoAbierto[K, V]{clave, valor}
}

func crearTabla[K comparable, V any](largo int) []TDALista.Lista[*elementoAbierto[K, V]] {
	tabla := make([]TDALista.Lista[*elementoAbierto[K, V]], largo)
	for i := range tabla {
		tabla[i] = TDALista.CrearListaEnlazada[*elementoAbierto[K, V]]()
	}
	return tabla
}

type hashAbierto[K comparable, V any] struct {
	elementos []TDALista.Lista[*elementoAbierto[K, V]]
	cantidad  int
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	hash := new(hashAbierto[K, V])
	hash.elementos = crearTabla[K, V](_CAPACIDAD_INICIAL)
	return hash
}

func (hash *hashAbierto[K, V]) superoCargaPermitida() bool {
	return 100*(hash.cantidad) >= len(hash.elementos)*_MAXIMA_CARGA
}

func (hash *hashAbierto[K, V]) ocupaMuchaMemoria() bool {
	return len(hash.elementos) >= 2*_CAPACIDAD_INICIAL && 100*hash.cantidad <= len(hash.elementos)*_MINIMA_CARGA
}

func (hash *hashAbierto[K, V]) buscarElemento(listaPosicion TDALista.Lista[*elementoAbierto[K, V]], clave K) *elementoAbierto[K, V] {

	var res *elementoAbierto[K, V] = nil
	listaPosicion.Iterar(func(elemento *elementoAbierto[K, V]) bool {
		if elemento.clave == clave {
			res = elemento
			return false
		}
		return true
	})

	return res
}

func (hash *hashAbierto[K, V]) iterarInterno(visitar func(*elementoAbierto[K, V]) bool) {
	i := 0
	seguir := true
	for seguir && i < len(hash.elementos) {
		hash.elementos[i].Iterar(func(elemento *elementoAbierto[K, V]) bool {
			seguir = visitar(elemento)
			return seguir
		})

		i++
	}
}

func (hash *hashAbierto[K, V]) redimensionar(nuevoLargo int) {
	nuevasListas := crearTabla[K, V](nuevoLargo)

	hash.iterarInterno(func(elemento *elementoAbierto[K, V]) bool {
		nuevasListas[aplicarFuncionHash(elemento.clave, len(nuevasListas))].InsertarUltimo(elemento)
		return true
	})

	hash.elementos = nuevasListas

}

func (hash *hashAbierto[K, V]) dameLista(clave K) TDALista.Lista[*elementoAbierto[K, V]] {
	return hash.elementos[aplicarFuncionHash(clave, len(hash.elementos))]
}

func (hash *hashAbierto[K, V]) Guardar(clave K, valor V) {
	listaCorrespondiente := hash.dameLista(clave)
	existente := hash.buscarElemento(listaCorrespondiente, clave)
	if existente != nil {
		existente.valor = valor
		return
	}
	hash.cantidad++
	if hash.superoCargaPermitida() {
		hash.redimensionar(2 * len(hash.elementos))
		listaCorrespondiente = hash.dameLista(clave)
	}

	listaCorrespondiente.InsertarUltimo(crearElementoAbierto(clave, valor))
}

func (hash *hashAbierto[K, V]) Pertenece(clave K) bool {
	return hash.buscarElemento(hash.dameLista(clave), clave) != nil
}

func (hash *hashAbierto[K, V]) Obtener(clave K) V {
	elemento := hash.buscarElemento(hash.dameLista(clave), clave)
	if elemento == nil {
		panic(ERROR_NO_ESTABA)
	}
	return elemento.valor
}

func (hash *hashAbierto[K, V]) Borrar(clave K) V {




	iterador := hash.dameLista(clave).Iterador()

	for iterador.HaySiguiente() && iterador.VerActual().clave != clave {
		iterador.Siguiente()
	}

	if !iterador.HaySiguiente() {
		panic(ERROR_NO_ESTABA)
	}

	elem := iterador.Borrar()


	if hash.ocupaMuchaMemoria() {
		hash.redimensionar(len(hash.elementos) / 2)
	}

	hash.cantidad--
	return elem.valor
}

func (hash *hashAbierto[K, V]) Cantidad() int {
	return hash.cantidad
}

func (hash *hashAbierto[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	i := 0
	seguir := true
	//
	for seguir && i < len(hash.elementos) {
		hash.elementos[i].Iterar(func(elemento *elementoAbierto[K, V]) bool {
			seguir = visitar(elemento.clave, elemento.valor)
			return seguir
		})

		i++
	}
}

// Iterador externo

func (hash *hashAbierto[K, V]) Iterador() IterDiccionario[K, V] {
	return creariteradorCerrado(hash)
}

type iteradorAbierto[K comparable, V any] struct {
	referencia     *hashAbierto[K, V]
	posActual      int
	iteradorActual TDALista.IteradorLista[*elementoAbierto[K, V]]
}

func creariteradorCerrado[K comparable, V any](referencia *hashAbierto[K, V]) IterDiccionario[K, V] {
	iterador := new(iteradorAbierto[K, V])

	iterador.referencia = referencia
	iterador.iteradorActual = referencia.elementos[0].Iterador()
	iterador.ajustarIteradorActual()
	return iterador
}

func (iterador *iteradorAbierto[K, V]) ajustarIteradorActual() {
	if iterador.iteradorActual == nil || iterador.iteradorActual.HaySiguiente() {
		return
	}

	aProbar := iterador.posActual + 1
	for aProbar < len(iterador.referencia.elementos) && !iterador.iteradorActual.HaySiguiente() {
		iterador.iteradorActual = iterador.referencia.elementos[aProbar].Iterador() // se repite una vez de forma innecesaria
		aProbar++
	}

	if !iterador.iteradorActual.HaySiguiente() {
		iterador.iteradorActual = nil
	}
	iterador.posActual = aProbar - 1
}

func (iterador *iteradorAbierto[K, V]) panicTermino() {
	if !iterador.HaySiguiente() {
		panic(ERROR_ITERADOR_TERMINO)
	}
}

func (iterador *iteradorAbierto[K, V]) HaySiguiente() bool {
	iterador.ajustarIteradorActual()
	return iterador.iteradorActual != nil
}

func (iterador *iteradorAbierto[K, V]) VerActual() (K, V) {
	iterador.panicTermino()
	elemento := iterador.iteradorActual.VerActual()
	return elemento.clave, elemento.valor
}

func (iterador *iteradorAbierto[K, V]) Siguiente() K {
	iterador.panicTermino()

	return iterador.iteradorActual.Siguiente().clave

}
