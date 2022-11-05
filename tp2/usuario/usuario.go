package usuario

import heap "cola_prioridad"
import "strings"

//import "tp2/post"

func Compare(user, other string) int {
	nombre := strings.ToLower(user)
	nombreOtro := strings.ToLower(other)
	toReturn := -1

	min := len(nombre)
	if len(nombreOtro) < min {
		toReturn = 1
		min = len(nombreOtro)
	} else if len(nombreOtro) == min {
		toReturn = 0
	}

	for i := 0; i < min; i++ {
		if nombre[i] < nombreOtro[i] {
			return -1
		} else if nombre[i] > nombreOtro[i] {
			return 1
		}
	}

	return toReturn

}

type Usuario interface {
	Feed() heap.ColaPrioridad[int]
	Nombre() string
	AntesQue(Usuario) bool
}

type usuario struct {
	nombre string
	indice int
	feed   heap.ColaPrioridad[int]
}

func FuncionAfinidad(usuario1, usuario2, usuario3 Usuario) int {
	juez := usuario1.(*usuario)
	opcion1 := usuario2.(*usuario)
	opcion2 := usuario3.(*usuario)

	diff1 := juez.indice - opcion1.indice
	if diff1 < 0 {
		diff1 = -diff1
	}
	diff2 := juez.indice - opcion2.indice
	if diff2 < 0 {
		diff2 = -diff2
	}

	return diff2 - diff1
}

func CrearUsuario(nombre string, indice int, recomendador func(Usuario, int, int) int) Usuario {
	user := new(usuario)
	user.nombre = nombre
	user.indice = indice

	user.feed = heap.CrearHeap[int](func(id int, id2 int) int {
		return recomendador(user, id, id2)

	})
	return user
}

func (user *usuario) Nombre() string {
	return user.nombre
}

func (user *usuario) Feed() heap.ColaPrioridad[int] {
	return user.feed
}

func (user *usuario) AntesQue(other Usuario) bool {
	return Compare(user.Nombre(), other.Nombre()) == 1
}
