package sesion

import hash "hash/hashCerrado"
import "strconv"
import "errors"

import "tp2/utilities"
import "tp2/posts"
import "tp2/usuarios"
import "tp2/managers"

const (
	ERROR_YA_LOGEO           = "Error: Ya habia un usuario loggeado"
	ERROR_USUARIO_INVALIDO   = "Error: usuario no existente"
	ERROR_NO_LOGEO           = "Error: no habia usuario loggeado"
	ERROR_VER_POST           = "Usuario no loggeado o no hay mas posts para ver"
	ERROR_LIKE_POST          = "Error: Usuario no loggeado o Post inexistente"
	ERROR_MOSTRAR_LIKES_POST = "Error: Post inexistente o sin likes"
)

type Sesion interface {
	Login(string) error
	Logout() error
	Publicar(string) error
	VerSiguientePost() (string, error)
	Likear(string) error
	MostrarLikes(string) (string, error)
}

type sesionStruct struct {
	postManager managers.PostManager[int,usuarios.UsuarioAlgogram,posts.PostAlgogram]
	userManager managers.UserManager[string,usuarios.UsuarioAlgogram,int,posts.PostAlgogram]
	loggeado    *usuario.UsuarioAlgogram
}


func CrearSesion(archivo_usuarios string) (Sesion, error) {
	sesion := new(sesionStruct)
	sesion.postManager = managers.CrearPostManager[usuarios.UsuarioAlgogram,posts.PostAlgogram]()
	sesion.userManager = managers.CrearUserManagerAlgogram[usuarios.UsuarioAlgogram,int,posts.PostAlgogram](sesion.postManager.AntesQue)
	
	ind:= 0

	err := utilities.LeerArchivo(archivo_usuarios, func(nombre string) bool {
		sesion.userManager.Agregar(usuario.CrearUsuarioAlgogram(nombre, ind))
		ind++
		return true
	})

	
	if err != nil {
		return nil, err
	}

	return sesion, nil
}

func (sesion *sesionStruct) Login(nombre string) error {
	if sesion.loggeado != nil {
		return errors.New(ERROR_YA_LOGEO)
	}
	if !sesion.userManager.Existe(nombre) {
		return errors.New(ERROR_USUARIO_INVALIDO)
	}

	sesion.loggeado = sesion.userManager.Obtener(nombre)


	return nil
}
func (sesion *sesionStruct) Logout() error {
	if sesion.loggeado == nil {
		return errors.New(ERROR_NO_LOGEO)
	}

	sesion.loggeado = nil
	return nil

}
func (sesion *sesionStruct) Publicar(contenido string) error {
	if sesion.loggeado == nil {
		return errors.New(ERROR_NO_LOGEO)
	}

	post:= posts.CrearPostAlgogram(sesion.loggeado, contenido)
	id := sesion.postManager.Agregar(post)
	sesion.userManager.AgregarPost(id,post)

	return nil
}

func (sesion *sesionStruct) VerSiguientePost() (string, error) {
	if sesion.loggeado == nil || sesion.loggeado.Feed().EstaVacia() {
		return "", errors.New(ERROR_VER_POST)
	}

	return sesion.postManager.ObtenerPost(sesion.loggeado.Feed().Desencolar()).String(), nil

}
func (sesion *sesionStruct) Likear(idStr string) error {

	id, err := strconv.Atoi(idStr)

	if sesion.loggeado == nil || err != nil || !sesion.postManager.Existe(id){
		return errors.New(ERROR_LIKE_POST)
	}
	
	sesion.postManager.ObtenerPost(id).AgregarLike(sesion.loggeado)
	
	return nil
}
func (sesion *sesionStruct) MostrarLikes(idStr string) (string, error) {
	id, err := strconv.Atoi(idStr)
	
	if err != nil || !sesion.postManager.Existe(id){
		return "", errors.New(ERROR_MOSTRAR_LIKES_POST)
	}

	post := sesion.postManager.ObtenerPost(id)
	if post.CantidadLikes() == 0 {
		return "", errors.New(ERROR_MOSTRAR_LIKES_POST)
	}

	return post.MostrarLikes(), nil
}
