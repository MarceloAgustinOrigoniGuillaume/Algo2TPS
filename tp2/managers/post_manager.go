package managers

import "tp2/usuarios"
import "tp2/posts"

const _CAPACIDAD_INICIAL = 10



// el post manager lo unico que restringe es el tipo de identificador, se podria decir no tiene porque ser especifico
// a manejo de posts, pero en este caso lo es.
type postManagerArreglo[D usuarios.DatosUsuario, P posts.Post[D]] struct {
	postsArr []P
	ultimoId int
}
func CrearPostManagerArreglo[D usuarios.DatosUsuario, P posts.Post[D]]() PostManager[int,D,P] {
	return &postManagerArreglo[D,P]{make([]P, _CAPACIDAD_INICIAL), -1}
}

func (manager *postManagerArreglo[D,P]) redimensionar(nuevoLargo int) {
	postsNew := make([]P, nuevoLargo)

	copy(postsNew, manager.postsArr)
	manager.postsArr = postsNew
}


func (manager *postManagerArreglo[D,P]) Agregar(post P) int {
	manager.ultimoId++

	if manager.ultimoId >= len(manager.postsArr) {
		manager.redimensionar(2 * manager.ultimoId)
	}

	manager.postsArr[manager.ultimoId] = post

	return manager.ultimoId
}

// no se verifica nada en obtener para mayor rapidez, te saltara el panic de indices si lo usas mal. Problema de Barbara.
func (manager *postManagerArreglo[D,P]) Obtener(id int) P {
	return manager.postsArr[id]
}


func (manager *postManagerArreglo[D,P]) Existe(id int) bool {
	return id >= 0 && id <= manager.ultimoId // si borrar fuese posible habria que chequear que sea valido el valor.
}


func (manager *postManagerArreglo[D,P]) AntesQue(post1,post2 int) bool {
	return post1 < post2
}
