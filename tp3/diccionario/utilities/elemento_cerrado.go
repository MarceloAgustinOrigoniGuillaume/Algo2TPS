package utilities

type Status = int

const (
	_VACIO   Status = 0
	_BORRADO Status = -1
	_OCUPADO Status = 1
)

type ElementoCerrado[K comparable] interface {
	Key() K
	Status() Status
}

func deberiaSeguir[K comparable](elemento ElementoCerrado[K], clave K) bool {
	return elemento.Status() != _VACIO && elemento.Key() != clave
}

// se podrian haber usado enums en vez de bool para ahorrar un if en guardar.
// tambien se probo devolviendo punteros, y era mas lento.
func BuscarPosicionElementoCerrado[K comparable](elementos []ElementoCerrado[K], clave K) (int, bool) {

	indiceRef := -1

	posInicial := AplicarFuncionDeHash(clave, len(elementos))

	i := posInicial

	for i < len(elementos) && deberiaSeguir(elementos[i], clave) {
		if indiceRef == -1 && elementos[i].Status() == _BORRADO { // se agarra el primer borrado por defecto
			indiceRef = i
		}
		i++
	}
	if i < len(elementos) { // significaria deberiaSeguir fue false, es decir vacio o igual clave.
		return i, elementos[i].Status() == _OCUPADO
	}

	i = 0

	for i < posInicial && deberiaSeguir(elementos[i], clave) {
		if indiceRef == -1 && elementos[i].Status() == _BORRADO { // se agarra el primer borrado por defecto
			indiceRef = i
		}
		i++
	}

	if i < posInicial { // significaria deberiaSeguir fue false, es decir vacio o igual clave.
		return i, elementos[i].Status() == _OCUPADO
	}

	return indiceRef, false // no se encontro, se devuelve la referencia del borrado
}
