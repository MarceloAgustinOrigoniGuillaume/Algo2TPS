package posts

import "tp2/usuarios"
import "tp2/utilities"
import abb "diccionario"

import "fmt"

/*
	La implementacion de un Post especifica a Algogram implementaria un PostLikeable y solo requeriria
	un autor basico, que tenga nombre, pero para su uso mas simple desde otros componentes requerira el usuario Algogram
	Guardaria los likes manteniendo un orden guiado de forma lexicografica.
	Y al mostrar los likes solo muestra el nombre del usuario. Para esto el hecho de guardar el Usuario en si en el abb
	es innecesario, pero para casos donde se quiera hacer una muestra algo mas compleja podria servir.
*/

type PostAlgogram struct {
	autor     usuarios.UsuarioAlgogram
	contenido string
	likes     abb.DiccionarioOrdenado[string, usuarios.UsuarioAlgogram]
}

func CrearPostAlgogram(autor usuarios.UsuarioAlgogram, contenido string) PostLikeable[usuarios.UsuarioAlgogram] {
	return &PostAlgogram[T]{autor, contenido, abb.CrearABB[string, usuarios.UsuarioAlgogram](utilities.CompareLexico)}
}

func (post *PostAlgogram) Autor() usuarios.UsuarioAlgogram {
	return post.autor
}

func (post *PostAlgogram) CantidadLikes() int {
	return post.likes.Cantidad()
}

func (post *PostAlgogram) AgregarLike(liker usuarios.UsuarioAlgogram) {
	post.likes.Guardar(liker.Nombre(), liker)
}
func (post *PostAlgogram) MostrarLikes() string {
	res := fmt.Sprintf("El post tiene %d likes:", post.CantidadLikes())
	post.likes.Iterar(func(nombre string, _ string) bool {
		res += "\n\t" + nombre
		return true
	})

	return res
}

func (post *PostAlgogram) String() string {
	return fmt.Sprintf("Post ID %d\n%s dijo: %s\nLikes: %d", post.id, post.Autor().Nombre(), post.contenido, post.CantidadLikes())
}
