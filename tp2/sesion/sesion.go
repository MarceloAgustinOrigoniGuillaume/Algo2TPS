package sesion
import "tp2/posts"
import "tp2/usuario"
import hash "hash/hashCerrado"
import "strconv"
import "errors"


const (
	ERROR_YA_LOGEO = "Error: Ya habia un usuario loggeado"
	ERROR_USUARIO_INVALIDO = "Error: usuario no existente"
	ERROR_NO_LOGEO = "Error: no habia usuario loggeado"
	ERROR_VER_POST = "Usuario no loggeado o no hay mas posts para ver"
	ERROR_LIKE_POST = "Error: Usuario no loggeado o Post inexistente"
	ERROR_MOSTRAR_LIKES_POST = "Error: Post inexistente o sin likes"
)

type Sesion interface {
	
	Login(string) error
	Logout() error
	Publicar(string) error
	VerSiguientePost() (string,error)
	Likear(string) error
	MostrarLikes(string) (string,error)
}

type sesionStruct struct{
	usuarios hash.Diccionario[string,usuario.Usuario]
	postManager posts.PostManager
	loggeado usuario.Usuario
}


func CrearSesion(archivo_usuarios string) (Sesion,error){
	sesion := new(sesionStruct)
	sesion.usuarios = hash.CrearHash[string,usuario.Usuario]()
	sesion.postManager = posts.CrearPostManager()
	indice := 0
	err := LeerArchivo(archivo_usuarios,func(nombre string) bool{
		sesion.usuarios.Guardar(nombre,usuario.CrearUsuario(nombre,indice, sesion.postManager.RecomendacionPara))
		indice++
		return true
	})

	if err != nil{
		return nil,err
	}

	return sesion,nil
}




func (sesion *sesionStruct) Login(nombre string) error{
	if sesion.loggeado!= nil{
		return errors.New(ERROR_YA_LOGEO)
	}

	if !sesion.usuarios.Pertenece(nombre){
		return errors.New(ERROR_USUARIO_INVALIDO)		
	}

	sesion.loggeado = sesion.usuarios.Obtener(nombre)
	return nil
}
func (sesion *sesionStruct) Logout() error{
	if sesion.loggeado== nil{
		return errors.New(ERROR_NO_LOGEO)
	}

	sesion.loggeado = nil
	return nil

}
func (sesion *sesionStruct) Publicar(contenido string) error{
	if sesion.loggeado== nil{
		return errors.New(ERROR_NO_LOGEO)
	}

	id := sesion.postManager.AgregarPost(sesion.loggeado,contenido).Id()

	sesion.usuarios.Iterar(func (nombre string, user usuario.Usuario) bool {
		if nombre != sesion.loggeado.Nombre(){
			user.Feed().Encolar(id)
		}
		return true
	})

	return nil
}

func (sesion *sesionStruct) VerSiguientePost() (string,error){
	if sesion.loggeado== nil || sesion.loggeado.Feed().EstaVacia(){
		return "",errors.New(ERROR_VER_POST)
	}

	return sesion.postManager.ObtenerPost(sesion.loggeado.Feed().Desencolar()).String(),nil

}
func (sesion *sesionStruct) Likear(idStr string) error{

	id,err := strconv.Atoi(idStr)

	if sesion.loggeado== nil || err != nil{
		return errors.New(ERROR_LIKE_POST)
	}
	post:= sesion.postManager.ObtenerPost(id)
	if post == nil{		
		return errors.New(ERROR_LIKE_POST)
	}

	post.AgregarLike(sesion.loggeado)
	return nil
}
func (sesion *sesionStruct) MostrarLikes(idStr string) (string,error){
	id,err := strconv.Atoi(idStr)
	if err != nil{
		return "",errors.New(ERROR_MOSTRAR_LIKES_POST)
	}

	post:= sesion.postManager.ObtenerPost(id)
	if post == nil || post.CantidadLikes() == 0{		
		return "",errors.New(ERROR_MOSTRAR_LIKES_POST)
	}

	return post.MostrarLikes(),nil
}
