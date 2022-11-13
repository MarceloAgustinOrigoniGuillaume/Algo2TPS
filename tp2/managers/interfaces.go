package managers

import "tp2/usuario"
import "tp2/posts"

// una simple representacion de un sistema que gestionaria ids asociados a elementos.
type IdManager[T comparable, D any] interface {
	Agregar(D) T 
	Obtener(T) *D
	Existe(T) bool

}

// Basado en lo que se necesitaria para AlgoGram inicialmente.
// Un post manager, es un gestor de Ids, agrega la funcionalidad tener un orden intrinsico desde los ids
// Se agrega un tercer tipo generico para poder dejar a la implementacion elegir el tipo de Post
type PostManager[T comparable, D usuarios.DatosUsuario, P posts.Post[D]] interface {
	IdManager[T,P] 
//	RecomendacionPara(D, T, T) int
	AntesQue(T,T) bool
}

// Una interfaz de getters esencialmente
type PairIdPost[D usuarios.DatosUsuario,I comparable, P posts.Post[D]] interface{
	Id() I
	Post() *P
}

// Una forma de generalizar las posibles interacciones entre usuarios y posts.
// En caso de Algogram VerPost equivaldria a ver SiguientePost ya que siempre mostraria uno nuevo ,etc.
type PairPostUser[D usuario.DatosUsuario, I comparable,P posts.Post[D]] interface {
	VerUsuario() *D
	VerPost() PairIdPost[D,I,P]
}

// Basado en lo que se necesitaria para AlgoGram inicialmente.
// El gestor de usuarios de algogram , es un gestor de Ids, y a la vez agrega la funcionalidad de afinidad y 
// la conexion a un heap de Posts... Este UserManager agrega eso mismo, I representaria el tipo de identificador del post
// No me agrada del todo usar 4 generics, pero si se quiere mantener generalidad se debe.
type UserManager[T comparable, K PairPostUser[D,I,P]] interface {
	IdManager[T,K] 
	AgregarPost(PairIdPost[D,I,P])
}
