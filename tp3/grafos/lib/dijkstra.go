package lib

import grafos "tp3/grafos"

import pila "tp3/pila"
import heap "tp3/heap"
import hash "tp3/diccionario"
import set "tp3/diccionario/set"

import "fmt"

func Dijkstra[V comparable, T grafos.Numero](grafo grafos.Grafo[V, T], origen V, visitar func(PairDistancia[V, T]) bool) {
	distancias := hash.CrearHash[V, T]()
	//distanciasMasCortas := hash.CrearHash[V,T]()
	visitados := set.CrearSetWith[V](grafo.CantidadVertices())

	//padres := hash.CrearHash[V,V]()
	aVisitar := heap.CrearHeap[PairDistancia[V, T]](comparadorDistancias[V, T])

	//visitados.Guardar(visitar(origen))
	distancias.Guardar(origen, 0)
	aVisitar.Encolar(CrearPairDistancia[V, T](nil, origen, 0))

	var pairVisitado PairDistancia[V, T]
	var distanciaActual T
	for !aVisitar.EstaVacia() {
		pairVisitado = aVisitar.Desencolar()
		if !visitados.Guardar(pairVisitado.Actual()) { // no se guardo, ya se analizo sus hijos

			/*
				// si esta distancia es menor, habia uno negativo, por lo que se podria tirar panic
				// se asume nadie usaria Dijkstra con pesos negativos
				distanciaActual = distancias.Obtener(pairVisitado.visitado)
			*/

			if pairVisitado.distancia < distancias.Obtener(pairVisitado.Actual()) { // fue mejorado? hubo pesos negativos, mal, no uses dijkstra
				fmt.Printf("\nHubo habia negativos? se quiso mejorar una segunda vez %v < %v\n", pairVisitado.distancia, distancias.Obtener(pairVisitado.Actual()))
				return //panic("Se uso Dijkstra con pesos negativos, Dijkstra no funciona con pesos negativos")
			}

			continue
		}

		if !visitar(pairVisitado) {
			return
		}

		visitadoCopy := pairVisitado.visitado
		visitado := &visitadoCopy

		grafo.IterarAdyacentes(visitadoCopy, func(ady V, peso T) bool {

			distanciaActual = pairVisitado.Distancia() + peso //grafo.ObtenerArista(pairVisitado.visitado,ady).Peso()
			if !distancias.Pertenece(ady) || distanciaActual < distancias.Obtener(ady) {
				distancias.Guardar(ady, distanciaActual)
				aVisitar.Encolar(CrearPairDistancia(visitado, ady, distanciaActual))
			}

			return true
		})
	}
}

func caminoDesdePadres[V comparable](padres hash.Diccionario[V, *V], origen V, dest V) []V {
	camino := pila.CrearPilaDinamica[V]()
	camino.Apilar(dest)
	length := 1
	for dest != origen {
		dest = *(padres.Obtener(dest))
		camino.Apilar(dest)
		length++
	}

	res := make([]V, length)

	for ind := 0; ind < length; ind++ {
		res[ind] = camino.Desapilar()
	}

	return res
}

func CaminosMinimosDijkstra[V comparable, T grafos.Numero](grafo grafos.Grafo[V, T], origen V) hash.Diccionario[V, *V] {
	padres := hash.CrearHash[V, *V]()

	Dijkstra[V](grafo, origen, func(pair PairDistancia[V, T]) bool {
		padres.Guardar(pair.Actual(), pair.Desde())
		return true
	})

	return padres
}

func CaminoMinimoDijkstraHasta[V comparable, T grafos.Numero](grafo grafos.Grafo[V, T], origen V, dest V) ([]V, error) {

	padres := hash.CrearHash[V, *V]()

	encontrado := false
	Dijkstra[V](grafo, origen, func(pair PairDistancia[V, T]) bool {
		padres.Guardar(pair.Actual(), pair.Desde())
		encontrado = pair.Actual() == dest
		return !encontrado
	})

	if !encontrado {
		return nil, CrearErrorGrafo("No existe camino entre los vertices, o el grafo tenia pesos negativos")
	}

	return caminoDesdePadres(padres, origen, dest), nil
}
