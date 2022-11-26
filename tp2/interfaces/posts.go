package interfaces

// Lo basico y escencial para ser considerado un post
// El tipo del autor se deja a criterio, normalmente dependera de la implementacion de Post
type Post[I comparable, T DatosUsuario] interface {
	comparable
	Id() I
	Autor() T
	String() string
}

// Un post con funcionalidad de likes

type PostLikeable[I comparable, T DatosUsuario] interface {
	Post[I, T]

	AgregarLike(T)
	CantidadLikes() int
	MostrarLikes() string
}

type PostAlgogram[U UsuarioAlgogram] interface {
	PostLikeable[int, U]
}
