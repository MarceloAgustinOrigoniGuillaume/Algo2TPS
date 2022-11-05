package posts
import "tp2/usuario"
import abb "diccionario"
import "fmt"
type Post interface {
	Id() int
	Autor() usuario.Usuario

	AgregarLike(usuario.Usuario)
	CantidadLikes() int
	MostrarLikes() string

	AntesQue(Post) bool
	String() string
}

type postStruct struct {
	id int
	autor usuario.Usuario
	contenido string
	likes abb.DiccionarioOrdenado[string,string]
}


func CrearPost(id int, autor usuario.Usuario, contenido string) Post{
	return &postStruct{id,autor,contenido, abb.CrearABB[string,string](usuario.Compare)}
}


func (post *postStruct) Id() int{
	return post.id
}
func (post *postStruct) Autor() usuario.Usuario{
	return post.autor
}

func (post *postStruct)	CantidadLikes() int{
	return post.likes.Cantidad()
}

func (post *postStruct) AgregarLike(liker usuario.Usuario){
	post.likes.Guardar(liker.Nombre(),liker.Nombre())
}
func (post *postStruct) MostrarLikes() string {
	res := fmt.Sprintf("El post tiene %d likes:",post.CantidadLikes())
	post.likes.Iterar(func (nombre string, _ string) bool{
		res += "\n\t"+nombre
		return true
	})

	return res
}

func (post *postStruct) AntesQue(otro Post) bool{
	return (post.id- otro.Id()) < 0
}

func (post *postStruct) String() string{
	return fmt.Sprintf("Post ID %d\n%s dijo: %s\nLikes: %d",post.id,post.Autor().Nombre(),post.contenido,post.CantidadLikes())
}
