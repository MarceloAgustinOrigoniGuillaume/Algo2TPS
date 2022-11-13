package usuario
import "tp2/usuarios"
import "tp2/posts"
import hash "hash/hashCerrado"

const _CAPACIDAD_INICIAL = 10


type userManagerAlgogram[D usuarios.UsuarioAlgogram,I comparable] struct {
	usuarios hash.Diccionario[string, D]
	instanciarCmpPosts func(*D) func(PairIdPost[D,I,P],PairIdPost[D,I,P]) int
}



// devuelve < 0 el usuario1 es mas afin con la opcion2. > 0 si el usuario1 es mas afin con la opcion1. = 0 si son igual de afines
func (manager *UserManagerAlgogram[D,I,P]) masAfinCon(usuario1,opcion1,opcion2 D) int {
	indice_target := usuario1.Indice()

	diff1 := indice_target - opcion1.Indice()
	if diff1 < 0 {
		diff1 = -diff1
	}

	diff2 := indice_target - opcion2.Indice()
	
	if diff2 < 0 {
		diff2 = -diff2
	}

	return diff2 - diff1

}



func CrearUserManagerAlgogram[D usuarios.UsuarioAlgogram,I comparable,P posts.PostLikeable[D]](recomendador func(I,I) bool) *UserManager[string,D,I,P]{
	manager := new(userManagerAlgogram[D,I,P])

	manager.usuarios = hash.CrearHash[string, D]()
	
	manager.instanciarCmpPosts = func(usuarioAsociado *D) func(PairIdPost[D,I,P],PairIdPost[D,I,P]) int{ 
		return func (post1 PairIdPost[D,I,P],post2 PairIdPost[D,I,P]) int{
			res:= manager.masAfinCon(usuarioAsociado,post1.Post().Autor(),post2.Post().Autor())

			if res == 0{
				res = -1
				if manager.cmpPostIds(post1.Id(),post2.Id()){
					res = 1
				}
			}

			return res
		}
	}

	return manager
}


func (manager *UserManagerAlgogram[D,I,P]) Agregar(usuario D) string {
	manager.usuarios.Guardar(usuario.Nombre(),manager.crearDatosAfinidad(&usuario))

	return usuario.Nombre()
}

func (manager *UserManagerAlgogram[D,I,P]) Existe(nombre string) bool {
	return manager.usuarios.Pertenece(nombre)
}

func (manager *UserManagerAlgogram[D,I,P]) Obtener(nombre string) *D {
	return manager.usuarios.Obtener(nombre).usuario
}


func (manager *UserManagerAlgogram[D,I,P]) AgregarPost(pair PairIdPost[D,I,P]){
	nombreAutor := post.Autor().Nombre()
	manager.usuarios.Iterar(func(nombre string,elem structAlgogram) bool { 
		if  nombreAutor != nombre {
			elem.postsAVer.Encolar(pair)
		}
		return true
	})
}
