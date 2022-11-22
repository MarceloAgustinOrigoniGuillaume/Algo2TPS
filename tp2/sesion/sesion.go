package sesion

import "strconv"
import "errors"

import "tp2/utilities"
import "tp2/interfaces"
import "tp2/managers"

//import "tp2/posts"
//import "tp2/usuarios"

const (
	ERROR_YA_LOGEO           = "Error: Ya habia un usuario loggeado"
	ERROR_USUARIO_INVALIDO   = "Error: usuario no existente"
	ERROR_NO_LOGEO           = "Error: no habia usuario loggeado"
	ERROR_VER_POST           = "Usuario no loggeado o no hay mas posts para ver"
	ERROR_LIKE_POST          = "Error: Usuario no loggeado o Post inexistente"
	ERROR_MOSTRAR_LIKES_POST = "Error: Post inexistente o sin likes"
)

type tUsuario = *usuarioAlgogram
type tPost = *postAlgogram[tUsuario]

type Sesion interface {
	Login(string) error
	Logout() error
	Publicar(string) error
	VerSiguientePost() (string, error)
	Likear(string) error
	MostrarLikes(string) (string, error)
}

type sesionStruct struct {
	postManager        interfaces.IdManager[int, tPost]
	userManager        interfaces.IdManager[string, tUsuario] //managers.UserManager[string,tUsuario,int,tPost]
	recomendador       interfaces.Recomendador[int, tUsuario, tPost]
	conexionesLoggeado interfaces.MapConexiones[tUsuario, tPost]
}

func CrearSesion(archivo_usuarios string) (Sesion, error) {
	sesion := new(sesionStruct)
	sesion.postManager = managers.CrearNumericalIdManager[tPost]()
	sesion.userManager = managers.CrearUserManagerAlgogram[tUsuario]()
	sesion.recomendador = managers.CrearEmptyRecomendadorAlgogram[tUsuario, tPost]()

	ultimoIndice := 0

	err := utilities.LeerArchivo(archivo_usuarios, func(nombre string) bool {
		usuario := crearUsuarioAlgogram(nombre, ultimoIndice)
		sesion.userManager.Agregar(usuario)
		sesion.recomendador.AgregarUsuario(usuario)
		ultimoIndice++
		return true
	})

	if err != nil {
		return nil, err
	}

	return sesion, nil
}

func (sesion *sesionStruct) Login(nombre string) error {
	if sesion.conexionesLoggeado != nil {
		return errors.New(ERROR_YA_LOGEO)
	}
	if !sesion.userManager.Existe(nombre) {
		return errors.New(ERROR_USUARIO_INVALIDO)
	}

	sesion.conexionesLoggeado = sesion.recomendador.ObtenerRecomendaciones(sesion.userManager.Obtener(nombre))

	return nil
}
func (sesion *sesionStruct) Logout() error {
	if sesion.conexionesLoggeado == nil {
		return errors.New(ERROR_NO_LOGEO)
	}

	sesion.conexionesLoggeado = nil
	return nil

}
func (sesion *sesionStruct) Publicar(contenido string) error {
	if sesion.conexionesLoggeado == nil {
		return errors.New(ERROR_NO_LOGEO)
	}

	post := crearPostAlgogram(sesion.postManager.NuevoId(), sesion.conexionesLoggeado.VerNodo(), contenido)

	sesion.postManager.Agregar(post)
	sesion.recomendador.AgregarPost(sesion.conexionesLoggeado.VerNodo(), post)

	return nil
}

func (sesion *sesionStruct) VerSiguientePost() (string, error) {
	if sesion.conexionesLoggeado == nil {
		return "", errors.New(ERROR_VER_POST)
	}

	post := sesion.conexionesLoggeado.ObtenerConexion()
	if post == nil {
		return "", errors.New(ERROR_VER_POST)
	}

	return (*post).String(), nil

}
func (sesion *sesionStruct) Likear(idStr string) error {

	id, err := strconv.Atoi(idStr)

	if sesion.conexionesLoggeado == nil || err != nil || !sesion.postManager.Existe(id) {
		return errors.New(ERROR_LIKE_POST)
	}

	sesion.postManager.Obtener(id).AgregarLike(sesion.conexionesLoggeado.VerNodo())

	return nil
}
func (sesion *sesionStruct) MostrarLikes(idStr string) (string, error) {
	id, err := strconv.Atoi(idStr)

	if err != nil || !sesion.postManager.Existe(id) {
		return "", errors.New(ERROR_MOSTRAR_LIKES_POST)
	}

	post := sesion.postManager.Obtener(id)
	if post.CantidadLikes() == 0 {
		return "", errors.New(ERROR_MOSTRAR_LIKES_POST)
	}

	return post.MostrarLikes(), nil
}
