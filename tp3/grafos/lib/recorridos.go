package lib

import grafos "tp3/grafos"
import cola "tp3/cola"
import set "tp3/diccionario/set"

func _dfs[V comparable, T any](grafo grafos.Grafo[V, T], origen V, visitar func(visitado V), visitados set.Set[V]) {
	visitados.Guardar(origen)
	visitar(origen)

	grafo.IterarAdyacentes(origen, func(ady V, _ T) bool {
		if !visitados.Pertenece(ady) {
			_dfs(grafo, ady, visitar, visitados)
		}

		return true
	})
}

func DFS[V comparable, T any](grafo grafos.Grafo[V, T], origen V, visitar func(visitado V)) {
	visitados := set.CrearSet[V]()
	_dfs(grafo, origen, visitar, visitados)
}

func DFS_ALL[V comparable, T any](grafo grafos.Grafo[V, T], visitar func(visitado V), terminoComponente func()) {
	visitados := set.CrearSet[V]()

	grafo.IterarVertices(func(vert V) bool {

		if !visitados.Pertenece(vert) {
			_dfs(grafo, vert, visitar, visitados)
			terminoComponente()
		}

		return true
	})

}

func BFS[V comparable, T any](grafo grafos.Grafo[V, T], origen V, visitar func(visitado V)) {
	visitados := set.CrearSet[V]()
	aVisitar := cola.CrearColaEnlazada[V]()

	visitados.Guardar(origen)
	aVisitar.Encolar(origen)

	visitar(origen)
	var visitado V

	for !aVisitar.EstaVacia() {
		visitado = aVisitar.Desencolar()

		grafo.IterarAdyacentes(visitado, func(ady V, _ T) bool {
			if !visitados.Pertenece(ady) {
				aVisitar.Encolar(ady)
				visitados.Guardar(ady)
				visitar(ady) // al ser bfs ya se sabe se va a visitar primero, ya que se usa una cola
			}

			return true
		})
	}
}
