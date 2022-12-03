package lib

import grafos "tp3/grafos"

func MSTKruskal[V comparable, T grafos.Numero](grafo grafos.Grafo[V,T],ordenador func([]Arista[V,T]) []Arista[V,T]) []Arista[V,T]{
	aristas := AristasOrdenadas(grafo,ordenador)

	aristasMST := make([]Arista[V,T],grafo.CantidadVertices()-1)// Al ser un arbol tiene E = V-1
	union := CrearUnionFind(grafo.ObtenerVertices())

	i:= 0
	for _,arista := range aristas{
		if union.Unite(arista.Desde(),arista.Hasta()){
			aristasMST[i] = arista
			i++
		}
	}


	return aristasMST

}



// No hay prim para vos
func MSTPrim[V comparable, T grafos.Numero](grafo grafos.Grafo[V,T],ordenador func([]Arista[V,T]) []Arista[V,T]) []Arista[V,T]{
	aristas := AristasOrdenadas(grafo,ordenador)

	aristasMST := make([]Arista[V,T],grafo.CantidadVertices()-1)// Al ser un arbol tiene E = V-1
	union := CrearUnionFind(grafo.ObtenerVertices())

	i:= 0
	for _,arista := range aristas{
		if union.Unite(arista.Desde(),arista.Hasta()){
			aristasMST[i] = arista
			i++
		}
	}


	return aristasMST

}