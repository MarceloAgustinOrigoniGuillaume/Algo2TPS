package lib

import "tp3/grafos"

import hash "tp3/diccionario"
import set "tp3/diccionario/set"

import "tp3/utils"
import "errors"

const ERROR_RECORRIDO = "No se encontro recorrido"

func ErrorRecorrido() error {
	return errors.New(ERROR_RECORRIDO)
}

func CrearErrorGrafo(err string) error {
	return errors.New(err)
}

// O(V + E)
func GradosDeSalida[V comparable, T any](grafo grafos.Grafo[V, T]) hash.Diccionario[V, int] {
	res := hash.CrearHash[V, int]() //make([]int, len(vertices))

	grafo.IterarVertices(func(vert V) bool {
		res.Guardar(vert, len(grafo.ObtenerAdyacentes(vert)))
		return true
	})

	return res
}

func GradosDeEntrada[V comparable, T any](grafo grafos.Grafo[V, T]) hash.Diccionario[V, int] {
	if !grafo.EsDirigido() {
		return GradosDeSalida(grafo)
	}

	res := hash.CrearHash[V, int]()

	grafo.IterarVertices(func(vert V) bool {
		res.Guardar(vert, 0)
		return true
	})

	grafo.IterarAristas(func(desde V, hasta V, peso T) bool {
		res.Guardar(hasta, res.Obtener(hasta)+1)
		return true
	})

	return res
}

// O(V + E)
func AristasDeSalida[V comparable, T any](grafo grafos.Grafo[V, T]) hash.Diccionario[V, hash.Diccionario[V, T]] {
	res := hash.CrearHash[V, hash.Diccionario[V, T]]()

	grafo.IterarVertices(func(vert V) bool {
		res.Guardar(vert, hash.CrearHash[V, T]())
		return true
	})

	grafo.IterarAristas(func(desde V, hasta V, peso T) bool {
		res.Obtener(desde).Guardar(hasta, peso)
		return true
	})

	return res
}

func AristasDeSalidaYGradosEntrada[V comparable, T any](grafo grafos.Grafo[V, T]) (hash.Diccionario[V, hash.Diccionario[V, T]], hash.Diccionario[V, int]) {
	res := hash.CrearHash[V, hash.Diccionario[V, T]]()
	grados := hash.CrearHash[V, int]()

	grafo.IterarVertices(func(vert V) bool {
		res.Guardar(vert, hash.CrearHash[V, T]())
		grados.Guardar(vert, 0)
		return true
	})

	grafo.IterarAristas(func(desde V, hasta V, peso T) bool {
		res.Obtener(desde).Guardar(hasta, peso)
		grados.Guardar(hasta, grados.Obtener(hasta)+1)
		return true
	})

	return res, grados
}

func AristasNoDirigido[V comparable, T grafos.Numero](grafo grafos.Grafo[V, T]) []Arista[V, T] {
	aristas := make([]Arista[V, T], grafo.CantidadAristas()/2)
	visitados := set.CrearSetWith[V](grafo.CantidadVertices())
	i := 0

	grafo.IterarVertices(func(vert V) bool {
		grafo.IterarAdyacentes(vert, func(ady V, peso T) bool {
			if !visitados.Pertenece(ady) {
				aristas[i] = CrearAristaSimple(vert, ady, peso)
				i++
			}

			return true
		})

		visitados.Guardar(vert) // se guarda despues por si hay lazos
		return true
	})

	return aristas

}

func AristasDirigido[V comparable, T grafos.Numero](grafo grafos.Grafo[V, T]) []Arista[V, T] {
	aristas := make([]Arista[V, T], grafo.CantidadAristas())
	i := 0
	grafo.IterarAristas(func(desde V, hasta V, peso T) bool {
		aristas[i] = CrearAristaSimple(desde, hasta, peso)
		i++
		return true
	})

	return aristas
}

func Aristas[V comparable, T grafos.Numero](grafo grafos.Grafo[V, T]) []Arista[V, T] {
	var aristas []Arista[V, T]
	if grafo.EsDirigido() {
		aristas = AristasDirigido(grafo)
	} else {
		aristas = AristasNoDirigido(grafo)
	}

	return aristas
}

func AristasSortedQuicksort[V comparable, T grafos.Numero](grafo grafos.Grafo[V, T]) []Arista[V, T] {
	return AristasOrdenadas(grafo, QuickSortAristas[V, T])
}

func AristasSortedRadix[V comparable](grafo grafos.Grafo[V, int]) []Arista[V, int] {

	if !grafo.EsPesado() {
		return Aristas(grafo)
	}

	aristas := Aristas(grafo)

	return utils.QuickSort(aristas, 0, len(aristas)-1, func(arista1, arista2 Arista[V, int]) bool {
		return arista2.Compare(arista1) > 0 // arista1 < arista2
	})
}

func AristasOrdenadas[V comparable, T grafos.Numero](grafo grafos.Grafo[V, T], ordenador func([]Arista[V, T]) []Arista[V, T]) []Arista[V, T] {

	if !grafo.EsPesado() {
		return Aristas(grafo)
	}

	return ordenador(Aristas(grafo))
}

func QuickSortAristas[V comparable, T grafos.Numero](aristas []Arista[V, T]) []Arista[V, T] {
	return utils.QuickSort(aristas, 0, len(aristas)-1, func(arista1, arista2 Arista[V, T]) bool {
		return arista2.Compare(arista1) > 0 // arista1 < arista2
	})
}
