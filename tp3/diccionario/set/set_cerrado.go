package set

import "tp3/diccionario/utilities"

const (
	_CAPACIDAD_INICIAL            = 128
	_MAXIMA_CARGA                 = 80  // esta constante tendria unidad de 1%, osea 80 = 80%
	_MINIMA_CARGA                 = 5   // esta constante tendria unidad de 1%, osea 5 = 5%
	_100PORCIENTO                 = 100 // representa la totalidad, el 100% de algo en el mismo sistema de unidades
	//_PROPORCION                   = 2 // Seria en realidad 5/4, 1,2 pero se redondea para arriba
	ERROR_NO_ESTABA               = "La clave no pertenece al set"
	ERROR_ITERADOR_TERMINO        = "El iterador termino de iterar"
	_VACIO                 utilities.Status = 0
	_BORRADO               utilities.Status = -1
	_OCUPADO               utilities.Status = 1
)

type elementoCerrado[K comparable] struct {
	clave  K
	estado utilities.Status
}

func (elem elementoCerrado[K]) Key() K{
	return elem.clave
}

func (elem elementoCerrado[K]) Status() utilities.Status{
	return elem.estado
}

func crearElementoCerrado[K comparable](clave K) utilities.ElementoCerrado[K] {
	return elementoCerrado[K]{clave, _OCUPADO}
}

func crearElementoCerradoVacio[K comparable]() utilities.ElementoCerrado[K] {
	return elementoCerrado[K]{}
}

func crearElementoBorrado[K comparable]() utilities.ElementoCerrado[K] {
	elem:= elementoCerrado[K]{}
	elem.estado = _BORRADO
	return elem
}

func crearTabla[K comparable](largo int) []utilities.ElementoCerrado[K] {
	tabla := make([]utilities.ElementoCerrado[K], largo)
	for i := range tabla {
		tabla[i] = crearElementoCerradoVacio[K]()
	}
	return tabla
}

type setCerrado[K comparable] struct {
	elementos []utilities.ElementoCerrado[K]
	cantidad  int
	borrados  int
}

func CrearSet[K comparable]() Set[K] {
	set := new(setCerrado[K])
	set.elementos = crearTabla[K](_CAPACIDAD_INICIAL)
	return set
}

func CrearSetWith[K comparable](capacidad_inicial int) Set[K] {
	set := new(setCerrado[K])

	capacidad_inicial = (int)(capacidad_inicial * _100PORCIENTO/_MAXIMA_CARGA) + 1 // se asegura no se redimensione con esa cap
	if capacidad_inicial < _CAPACIDAD_INICIAL{
		capacidad_inicial = _CAPACIDAD_INICIAL
	}

	set.elementos = crearTabla[K](capacidad_inicial)
	return set
}


/*
	La unidad usada en la comparacion del factor de carga es el 1% para permitir mayor precision.

	Teniendo estas unidades una _MAXICA_CARGA de 80, seria 80% y representando esto como proporcion seria 80/100(0,8)
	Se tomara en cuenta divido por 100 ya que no se piensa usar porcentajes racionales, con coma, floats.

	Por lo que el limite fue superado para una maxima carga seria cantidad >= (80/100) * Capacidad, cantidad supero el 80% escencialmente.
	Para evitar la division, que es un operacion mas lenta, se multiplica miembro a miembro por 100.
	100 * cantidad >= 80 * Capacidad que si reemplazados por el nombre de la variable _MAXIMA_CARGA y constante _100Porciento

	_100Porciento * cantidad >= _MAXIMA_CARGA * Capacidad , obviamente 	cantidad en realidad seria la suma de la cantidad y borrados
	para el caso especifico.

	Analogamente se llega al caso para ocupaMuchaMemoria. Solo que la condicion de ruptura al ser un porcentaje minimo es <= en vez de >=
	_100Porciento * cantidad <= _MINIMA_CARGA * Capacidad , y aca si seria solo la cantidad ya que la suma cantidad+borrados es solo creciente.

	Ademas ocupaMuchaMemoria tiene una primer condicion de capacidad>= 2* _CAPACIDAD_INICIAL para asegurar la capacidad nunca baja de esta.
*/

func (set *setCerrado[K]) superoCargaPermitida() bool {
	return _100PORCIENTO*(set.cantidad+set.borrados) >= _MAXIMA_CARGA*len(set.elementos)
}

func (set *setCerrado[K]) ocupaMuchaMemoria() bool {
	return len(set.elementos) >= 2*_CAPACIDAD_INICIAL && _100PORCIENTO*set.cantidad <= _MINIMA_CARGA*len(set.elementos)
}


func (set *setCerrado[K]) redimensionar(nuevoLargo int) {
	nuevos := crearTabla[K](nuevoLargo)
	set.borrados = 0

	set.Iterar(func(clave K) bool { // se podria usar punteros, pero no se vio especial mejora
		indice, _ := utilities.BuscarPosicionElementoCerrado(nuevos, clave)
		nuevos[indice] = crearElementoCerrado(clave)
		return true
	})

	set.elementos = nuevos

}

func (set *setCerrado[K]) Guardar(clave K) bool{
	if set.superoCargaPermitida() {
		set.redimensionar(4 * len(set.elementos))
	}
	indice, estaba := utilities.BuscarPosicionElementoCerrado(set.elementos, clave)
	if !estaba {
		set.cantidad++
		if set.elementos[indice].Status() == _BORRADO {
			set.borrados--
		}
		set.elementos[indice] = crearElementoCerrado(clave)
		return true
	}

	return false
}

// da el elemento si pertenece, sino panic
func (set *setCerrado[K]) dameElemento(ind int, estaba bool) *utilities.ElementoCerrado[K] {
	if !estaba {
		return nil
	}

	return &set.elementos[ind]
}

func (set *setCerrado[K]) Pertenece(clave K) bool {
	_, estaba := utilities.BuscarPosicionElementoCerrado(set.elementos, clave) // no me gusta tener que guardarlo pero bueno, para reutilizar directamente
	return estaba
}

func (set *setCerrado[K]) Borrar(clave K) bool {
	elemento := set.dameElemento(utilities.BuscarPosicionElementoCerrado(set.elementos, clave))

	if(elemento == nil){
		return false
	}
	
	set.cantidad--
	*elemento = crearElementoBorrado[K]()
	set.borrados++

	if set.ocupaMuchaMemoria() {
		set.redimensionar(len(set.elementos) / 2)
	}

	return true
}

func (set *setCerrado[K]) Cantidad() int {
	return set.cantidad
}

func (set *setCerrado[K]) Iterar(visitar func(clave K) bool) {
	i := 0
	for i < len(set.elementos) && (set.elementos[i].Status() != _OCUPADO || visitar(set.elementos[i].Key())) {
		i++
	}
}

// Iterador externo

func (set *setCerrado[K]) Iterador() IterSet[K] {
	return creariteradorCerrado(set)
}

type iteradorCerrado[K comparable] struct {
	referencia *setCerrado[K]
	posActual  int
}

func creariteradorCerrado[K comparable](referencia *setCerrado[K]) IterSet[K] {
	iterador := new(iteradorCerrado[K])

	iterador.referencia = referencia
	iterador.posActual = -1
	iterador.iterarSiguiente()
	return iterador
}

func (iterador *iteradorCerrado[K]) iterarSiguiente() {
	iterador.posActual++
	for iterador.posActual < len(iterador.referencia.elementos) && iterador.referencia.elementos[iterador.posActual].Status() != _OCUPADO {
		iterador.posActual++
	}
}

func (iterador *iteradorCerrado[K]) panicTermino() {
	if !iterador.HaySiguiente() {
		panic(ERROR_ITERADOR_TERMINO)
	}
}

func (iterador *iteradorCerrado[K]) HaySiguiente() bool {
	return iterador.posActual < len(iterador.referencia.elementos)
}

func (iterador *iteradorCerrado[K]) VerActual() K {
	iterador.panicTermino()
	elemento := iterador.referencia.elementos[iterador.posActual]
	return elemento.Key()
}

func (iterador *iteradorCerrado[K]) Siguiente() K {
	iterador.panicTermino()
	claveActual := iterador.referencia.elementos[iterador.posActual].Key()
	iterador.iterarSiguiente()
	return claveActual

}