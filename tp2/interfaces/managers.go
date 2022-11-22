package interfaces

// una simple representacion de un sistema que gestionaria ids asociados a elementos.
type IdManager[T comparable, D any] interface {
	Agregar(data D) T // idealmente siempre O(1)
	Existe(id T) bool // idealmente siempre O(1)
	Obtener(id T) D // idealmente siempre O(1)
	NuevoId() T

	Iterar(visitar func(T,D) bool)
}



// Representa una conexion entre dos tipos de datos, usando MapConexiones la idea es permitir una relacion
// one to many de un dato a otro, aunque dependiendo de la implementacion  en especifico, podria ser one to one.
type MapConexiones[U any, D any] interface{
	VerNodo() U // O(1)
	AgregarConexion(element D) // Tiene que ser max O(log p)
	ObtenerConexion() *D // Tiene que ser max O(log p)
}


//Recomendador esta para Controlar las MapConexiones entre usuarios y posts, usando la implementacion que le parezca 
// U = cant usuarios, P = cant posts
type Recomendador[I comparable,D DatosUsuario,P Post[I,D]] interface{
	AgregarUsuario(usuario D) // Una primitiva que no se va a usar, pero si se quisiese agregar nuevos usuarios
	AgregarPost(autor D, post P) // tiene que ser max O(U* log P)
	ObtenerRecomendaciones(usuario D) MapConexiones[D,P] // Tiene que ser max O(1)
}

