package algogram

import abb "diccionario"

import "tp2/interfaces"
import "tp2/utilities"
import "fmt"

// Se puso el generico por si se quiere cambiar
type postAlgogram[D interfaces.UsuarioAlgogram] struct {
	id        int
	autor     D
	contenido string
	likes     abb.DiccionarioOrdenado[string, D]
}

func crearPostAlgogram[D interfaces.UsuarioAlgogram](id int, autor D, contenido string) *postAlgogram[D] { //interfaces.PostLikeable[int,D] {
	return &postAlgogram[D]{id, autor, contenido, abb.CrearABB[string, D](utilities.CompareLexico)}
}

func (post *postAlgogram[D]) Id() int {
	return post.id
}

func (post *postAlgogram[D]) Autor() D {
	return post.autor
}

func (post *postAlgogram[D]) CantidadLikes() int {
	return post.likes.Cantidad()
}

func (post *postAlgogram[D]) AgregarLike(liker D) {
	post.likes.Guardar(liker.Nombre(), liker)
}
func (post *postAlgogram[D]) MostrarLikes() string {
	res := fmt.Sprintf("El post tiene %d likes:", post.CantidadLikes())
	post.likes.Iterar(func(nombre string, _ D) bool {
		res += "\n\t" + nombre
		return true
	})

	return res
}

func (post *postAlgogram[D]) String() string {
	return fmt.Sprintf("Post ID %d\n%s dijo: %s\nLikes: %d", post.id, post.Autor().Nombre(), post.contenido, post.CantidadLikes())
}
