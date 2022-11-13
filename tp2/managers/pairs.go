package managers
import "tp2/usuarios"
import "tp2/posts"
import heap "cola_prioridad"



// agrupacion para un caso generico de id y post
type pairIdPost[D usuarios.DatosUsuario, I comparable, P posts.Post[D]] struct{
	identificador I
	post *P
}

func (pair *pairIdPost[D,I,P]) Id() I {
	return pair.identificador
}

func (pair *pairIdPost[D,I,P]) Post() *P {
	return pair.post 
}


func CrearPairIdPost[D usuarios.DatosUsuario, I comparable, P posts.Post[D]](id I, post *P) PairIdPost[D,I,P]{
	return &pairIdPost[D,I,P]{id,post}
}


// agrupacion de Usuario y Posts generico

// Este struct Algogram seria la conexion de usuarios a posts.
// No me agrada del todo el uso de tantos generics, pero para mantener generalidad es necesario
type StructAlgogram[D usuarios.UsuarioAlgogram,I comparable, P posts.PostLikeable[D]] struct{
	usuario *D
	postsAVer heap.ColaPrioridad[PairIdPost[D,I,P]]
}

func (conexion *StructAlgogram[D,I,P]) VerUsuario() *D {
	return conexion.usuario
}


func (conexion *StructAlgogram[D,I,P]) VerPost() PairIdPost[D,I,P] {
	if conexion.postsAVer.EstaVacia(){
		return nil
	}
	return conexion.postsAVer.Desencolar().post
}

func CrearStructAlgogram[D usuarios.UsuarioAlgogram,I comparable, P posts.PostLikeable[D]](usuario *D,
		comparador func(PairIdPost[D,I,P],PairIdPost[D,I,P]) int) PairPostUser[D,I,P]{
	
	return StructAlgogram[I]{usuario,heap.CrearHeap[PairIdPost[D,I,P]](comparador)}
}
