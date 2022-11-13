package posts

import "tp2/usuario"

// Lo basico y escencial para ser considerado un post
// El tipo del autor se deja a criterio, normalmente dependera de la implementacion de Post
type Post[T usuario.DatosUsuario] interface { 
	Autor() T
	String() string
}


// Un post con funcionalidad de likes

type PostLikeable[T usuario.DatosUsuario] interface {
	Post[T]

	AgregarLike(T)
	CantidadLikes() int
	MostrarLikes() string

}