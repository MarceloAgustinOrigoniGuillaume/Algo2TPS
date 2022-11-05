package posts
import "tp2/usuario"
//import "fmt"
type PostManager interface {
	AgregarPost(usuario.Usuario,string) Post
	ObtenerPost(int) Post
	RecomendacionPara(usuario.Usuario,int,int) int
}

const _CAPACIDAD_INICIAL = 10


type postManagerArreglo struct{
	postsArr []Post
	ultimoId int
}

func (manager *postManagerArreglo) redimensionar(nuevoLargo int) {
	postsNew := make([]Post, nuevoLargo)

	copy(postsNew, manager.postsArr)
	manager.postsArr = postsNew
}



func CrearPostManager() PostManager{
	return &postManagerArreglo{make([]Post, _CAPACIDAD_INICIAL),-1}
}


func (manager *postManagerArreglo) AgregarPost(autor usuario.Usuario, contenido string) Post{
	manager.ultimoId++

	if manager.ultimoId >= len(manager.postsArr){
		manager.redimensionar(2*manager.ultimoId)
	}

	manager.postsArr[manager.ultimoId]=  CrearPost(manager.ultimoId,autor,contenido)

	return manager.postsArr[manager.ultimoId]
}

func (manager *postManagerArreglo) ObtenerPost(id int) Post{

	if id <0 || id> manager.ultimoId{
		return nil
	}

	return manager.postsArr[id]
}

func (manager *postManagerArreglo) RecomendacionPara(user usuario.Usuario, id1,id2 int) int{
		
	res := usuario.FuncionAfinidad(user,manager.postsArr[id1].Autor(),manager.postsArr[id2].Autor())

	//fmt.Printf("RES AFINIDAD %d??  de %s con %s y %s\n",res,user.Nombre(),manager.postsArr[id1].Autor().Nombre(),manager.postsArr[id2].Autor().Nombre())

	if res == 0{
		return id2-id1 // significa id2 es mayor, fue creado despues, id1 va antes.
	}

	return res

}
