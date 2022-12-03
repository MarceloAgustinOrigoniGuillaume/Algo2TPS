package grafos
/*
import hash "tp3/diccionario"
import "fmt"
import "golang.org/x/exp/constraints"

//T
type grafoGeneral[V comparable, T any] struct{
	hashConexiones hash.Diccionario[V,hash.Diccionario[V,T]]

	dirigido bool
	cantidadAristas int
}


// No tiene sentido hacer un grafo general no pesado. Si quiere algo de esa funcionalidad
// Tema de barbara
func NewGrafoGeneral[V comparable,T any](dirigido bool) Grafo[V,T]{
	grafo := new(grafoGeneral[V,T])
	grafo.dirigido = dirigido
	grafo.hashConexiones = hash.CrearHash[V,hash.Diccionario[V,T]]()

	return grafo
}

func (grafo *grafoGeneral[V,T]) EsPesado() bool{
	return true
}

func (grafo *grafoGeneral[V,T]) CantidadVertices() int{
	return grafo.hashConexiones.Cantidad()
}

func (grafo *grafoGeneral[V,T]) CantidadAristas() int{
	return grafo.cantidadAristas
}

func (grafo *grafoGeneral[V,T]) EsDirigido() bool{
	return grafo.dirigido
}


func (grafo *grafoGeneral[V,T]) AgregarVertice(vertice V) {
	if grafo.hashConexiones.Pertenece(vertice){
		fmt.Printf("\nWarning: se intento agregar un vertice mas de una vez: %v\n",vertice)
		//panic("Se intento agregar un vertice que ya existia")
		return
	}
	grafo.hashConexiones.Guardar(vertice,hash.CrearHash[V,T]())
}




func (grafo *grafoGeneral[V,T]) BorrarVertice(vertice V) {
	if !grafo.hashConexiones.Pertenece(vertice){
		fmt.Printf("\nWarning: se intento borrar un vertice inexistente : %v\n",vertice)
		//panic("Se intento agregar un vertice que ya existia")
		return
	}
	grafo.hashConexiones.Borrar(vertice)

	grafo.hashConexiones.Iterar(func (_ V,conn hash.Diccionario[V,T]) bool{
		if conn.Pertenece(vertice){
			conn.Borrar(vertice)
		}

		return true
	})
}



func (grafo *grafoGeneral[V,T]) BorrarArista(desde V,hasta V) {
	if (!grafo.hashConexiones.Pertenece(desde) || !grafo.hashConexiones.Pertenece(hasta)){
		panic(fmt.Sprintf("Se quiso borrar una arista de algun vertice inexistente: %v->%v",desde,hasta))
	}

	ady := grafo.hashConexiones.Obtener(desde)

	if(!ady.Pertenece(hasta)){
		fmt.Printf("\nWarning: se intento borrar una arista inexistente : %v->%v\n",desde,hasta)
		return
	}

	grafo.cantidadAristas--

	ady.Borrar(hasta)

	if !grafo.dirigido{
		grafo.cantidadAristas--
		grafo.hashConexiones.Obtener(hasta).Borrar(desde)
	}
}


func (grafo *grafoGeneral[V,T]) AgregarArista(desde V,hasta V, peso T) {

	if (!grafo.hashConexiones.Pertenece(desde) || !grafo.hashConexiones.Pertenece(hasta)){
		panic(fmt.Sprintf("Se quiso agregar una arista de algun vertice inexistente: %v->%v",desde,hasta))
	}

	aristasDesde := grafo.hashConexiones.Obtener(desde)
	diffCantidad := -aristasDesde.Cantidad()
	aristasDesde.Guardar(hasta,peso)
	diffCantidad += aristasDesde.Cantidad()

	if !grafo.dirigido{
		diffCantidad *= 2
		grafo.hashConexiones.Obtener(hasta).Guardar(desde,peso)
	}

	grafo.cantidadAristas+= diffCantidad
}

func (grafo *grafoGeneral[V,T]) ObtenerAdyacentes(vertice V) []V {
	if (!grafo.hashConexiones.Pertenece(vertice) ){
		panic("Se quiso obtener adyacentes de un vertice inexistente")
	}

	conn := grafo.hashConexiones.Obtener(vertice)
	res:= make([]V,conn.Cantidad())
	i:= 0
	conn.Iterar(func (hasta V,_ T) bool{
		res[i] = hasta
		i++
		return true
	})

	return res

}

func (grafo *grafoGeneral[V,T]) ObtenerAristas(vertice V) []Arista[V,T] {
	if (!grafo.hashConexiones.Pertenece(vertice) ){
		panic("Se quiso obtener adyacentes de un vertice inexistente")
	}

	conn := grafo.hashConexiones.Obtener(vertice)
	res:= make([]Arista[V,T],conn.Cantidad())
	i:= 0
	conn.Iterar(func (hasta V,cantidad T) bool{
		res[i] = CrearArista(vertice,hasta,cantidad)
		i++
		return true
	})

	return res

}






func (grafo *grafoGeneral[V,T]) ObtenerVertices() []V {
	vertices := make([]V, grafo.hashConexiones.Cantidad())

	i:= 0
	grafo.IterarVertices(func(vert V) bool{
		vertices[i] = vert
		i++
		return true
	})

	return vertices
}




func (grafo *grafoGeneral[V,T]) IterarVertices(visitar func(vert V) bool) {
	grafo.hashConexiones.Iterar(func (vertice V,_ hash.Diccionario[V,T]) bool{
		return visitar(vertice)
	})
}

func (grafo *grafoGeneral[V,T]) IterarAristas(visitar func(desde V, hasta V, peso T) bool) {
	keepIterating := true
	grafo.hashConexiones.Iterar(func (desde V,aristas hash.Diccionario[V,T]) bool{
		aristas.Iterar(func (hasta V,peso T) bool{
			keepIterating = visitar(desde,hasta,peso)
			return keepIterating
		})
		return keepIterating
	})
}








func (grafo *grafoGeneral[V,T]) MostrarTest(connString string) string {
	res := ""

	grafo.hashConexiones.Iterar(func (vertice V,conn hash.Diccionario[V,T]) bool{
		res += fmt.Sprintf("------->%v\n",vertice)

		conn.Iterar(func (hasta V,cantidad T) bool{
			res += fmt.Sprintf(connString,hasta,cantidad)+"\n"
			return true
		})

		return true
	})


	if res == ""{
		res = "Sin vertices....\n"
	}
	return res
}

*/