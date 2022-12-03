package generators

import "tp3/grafos"
import "tp3/pila"

//import "math/rand"

// Lower suffix elements should be before
func BuildGrafoTopologico[T grafos.Numero](dirigido bool, labels []string, vert_quantity int) grafos.Grafo[string, T] {
	var last_sufix []string = nil
	curr_sufix := make([]string, len(labels))
	ind := 0

	aristas := pila.CrearPilaDinamica[[2]string]()

	grafo := BuildVertices[T](dirigido, labels, vert_quantity,
		func(vert string) {
			curr_sufix[ind] = vert
			ind++
		},
		func(_ string) {
			ind = 0
			if last_sufix != nil {
				for _, antes := range last_sufix {
					for _, despues := range curr_sufix {
						aristas.Apilar([2]string{antes, despues})
					}
				}
				copy(last_sufix, curr_sufix)
			} else {
				last_sufix = curr_sufix
			}

			curr_sufix = make([]string, len(labels))
		})

	// Checking for not completed suffix
	if last_sufix != nil {
		for _, antes := range last_sufix {
			for i := 0; i < ind; i++ {
				aristas.Apilar([2]string{antes, curr_sufix[i]})
			}
		}
	}

	var curr [2]string
	for !aristas.EstaVacia() { // adding aristas
		curr = aristas.Desapilar()
		grafo.AgregarArista(curr[0], curr[1], 1)
	}

	return grafo
}
