package grafos
/*
import hash "tp3/diccionario"
import "fmt"

const _CAPACIDAD_INICIAL = 10


type grafoAdyacencia[V comparable] struct{
	matrizAdyacencias [][]float32
	vertices []V
	indices hash.Diccionario[V,int]

	dirigido bool
	pesado bool

	indUltimoVertice int
}


func CrearGrafoSinRemocion[V comparable](dirigido bool, pesado bool) Grafo[V]{
	grafo := new(grafoAdyacencia[V])
	grafo.dirigido = dirigido
	grafo.pesado = pesado
	grafo.indices = hash.CrearHash[V,int]()
	grafo.indUltimoVertice= -1

	grafo.matrizAdyacencias = make([][]float32, _CAPACIDAD_INICIAL)
	grafo.matrizAdyacencias = make([][]float32, _CAPACIDAD_INICIAL)

	for i:= 0; i<_CAPACIDAD_INICIAL; i++{
		grafo.matrizAdyacencias[i] = make([]float32,_CAPACIDAD_INICIAL)
	}


	return grafo


}

func appendTimes[T any](slice []T, times int ,instanciador func() T){
	for i:= 0;i<times;i++{
		slice.Append(instanciador())
	}
}

func (grafo *grafoAdyacencia[V]) esPesado() bool{
	return grafo.pesado
}

func (grafo *grafoAdyacencia[V]) esDirigido() bool{
	return grafo.dirigido
}

func (grafo *grafoAdyacencia[V]) chequearIndicesMatriz(){

	if(grafo.indUltimoVertice < len(grafo.matrizAdyacencias)){
		return
	}

	nuevos := grafo.indUltimoVertice-len(grafo.matrizAdyacencias)+1 // deberia ser siempre uno
	nuevaLen := len(grafo.matrizAdyacencias)+nuevos

	for i:= 0;i<len(grafo.matrizAdyacencias); i++{
		appendTimes(grafo.matrizAdyacencias[i],nuevos, func() { return 0 })
	}

	appendTimes(grafo.matrizAdyacencias,nuevos, func() { return make([]float32, nuevaLen) })

}


func (grafo *grafoAdyacencia[V]) AgregarVertice(vertice V) {
	if grafo.indices.Pertenece(vertice){
		fmt.Printf("\nWarning: se intento agregar un vertice mas de una vez: %v\n",vertice)
		//panic("Se intento agregar un vertice que ya existia")
		return
	}



	grafo.indUltimoVertice++
	grafo.chequearIndicesMatriz()

	grafo.indices.Guardar(vertice,grafo.indUltimoVertice)
}

func (grafo *grafoAdyacencia[V]) AgregarArista(arista Arista[V]) {


	if (!grafo.esPesado() && arista.Peso() != 1 && arista.Peso() != 0){
		panic("Se quiso agregar una arista pesada a un grafo no pesado")
	}

	if (!grafo.indices.Pertenece(arista.Desde()) || !grafo.indices.Pertenece(arista.Hasta())){
		panic("Se quiso agregar una arista a vertices inexistentes")
	}


	indice_desde :=grafo.indices.Obtener(arista.Desde())
	indice_hasta :=grafo.indices.Obtener(arista.Hasta())

	grafo.matrizAdyacencias[indice_desde][indice_hasta] = arista.Peso()

	if !grafo.dirigido{
		grafo.matrizAdyacencias[indice_hasta][indice_desde] = arista.Peso()
	}
}

func (grafo *grafoAdyacencia[V]) ObtenerVertices() []V {
	vertices := make([]V, grafo.indices.Cantidad())

	grafos.indices.Iterar(func (vertice V,i int) bool{
		vertices[i] = vertice // deberia estar garantizado i esta en el rango
		return true
	})

	return vertices
}

func (grafo *grafoAdyacencia[V]) ObtenerAdyacentes(vertice V) []Arista[V] {
	if (!grafo.indices.Pertenece(vertice) ){
		panic("Se quiso obtener adyacentes de un vertice inexistente")
	}

	indice_vertice := grafo.indices.Obtener(vertice)


	cantidadInt := 0
	for _,valor:= range matrizAdyacencias[indice_vertice]{
		if(valor != 0){
			cantidadInt++
		}
	}

	res:= make(Arista[V])

}



func (grafo *grafoAdyacencia[V]) MostrarTest(connString string) string {
	return fmt.Printf("\n Matriz ady :: \n%v\n",grafo.matrizAdyacencias)
}

*/