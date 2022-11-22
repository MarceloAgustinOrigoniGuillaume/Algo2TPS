package managers
/*
import "tp2/interfaces"
//posts.Post[]

const _CAPACIDAD_INICIAL = 10



// el post manager lo unico que restringe es el tipo de identificador, se podria decir no tiene porque ser especifico
// a manejo de posts, pero en este caso lo es.
type postManagerAlgogram[D interfaces.DatosUsuario,P interfaces.Post[int,D]] struct {
	posts numericalIdManager[P]
}

//CrearNumericalIdManager[P]
func CrearPostManagerAlgogram[D interfaces.DatosUsuario, P interfaces.Post[int,D]]() PostManager[int,D,P] {
	return &postManagerAlgogram[D,P]{CrearNumericalIdManager[P](), -1}
}

func (manager *postManagerAlgogram[D,P]) Agregar(post P) int {
	return manager.posts.Agregar(post)
}

// no se verifica nada en obtener para mayor rapidez, te saltara el panic de indices si lo usas mal. Problema de Barbara.
func (manager *postManagerAlgogram[D,P]) Obtener(id int) P {
	return manager.posts.Obtener(id)//manager.postsArr[id]
}


func (manager *postManagerAlgogram[D,P]) Existe(id int) bool {
	return manager.posts.Existe(id)//id >= 0 && id <= manager.ultimoId // si borrar fuese posible habria que chequear que sea valido el valor.
}

*/