package grafos

import hash "tp3/diccionario"
import "fmt"

type grafoNumerico[V comparable, T Numero] struct {
	hashConexiones hash.Diccionario[V, hash.Diccionario[V, T]]

	dirigido        bool
	pesado          bool
	cantidadAristas int
}

// Pesado por default
func GrafoNumericoPesado[V comparable, T Numero](dirigido bool) Grafo[V, T] {
	grafo := new(grafoNumerico[V, T])
	grafo.dirigido = dirigido
	grafo.pesado = true
	grafo.hashConexiones = hash.CrearHash[V, hash.Diccionario[V, T]]()

	return grafo

}

func GrafoNumericoNoPesado[V comparable](dirigido bool) Grafo[V, int] {
	grafo := new(grafoNumerico[V, int])
	grafo.dirigido = dirigido
	grafo.pesado = false
	grafo.hashConexiones = hash.CrearHash[V, hash.Diccionario[V, int]]()

	return grafo
}

func (grafo *grafoNumerico[V, T]) EsPesado() bool {
	return grafo.pesado
}
func (grafo *grafoNumerico[V, T]) EsDirigido() bool {
	return grafo.dirigido
}

func (grafo *grafoNumerico[V, T]) CantidadVertices() int {
	return grafo.hashConexiones.Cantidad()
}

func (grafo *grafoNumerico[V, T]) CantidadAristas() int {
	return grafo.cantidadAristas
}

func (grafo *grafoNumerico[V, T]) AgregarVertice(vertice V) {
	if grafo.hashConexiones.Pertenece(vertice) {
		fmt.Printf("\nWarning: se intento agregar un vertice mas de una vez: %v\n", vertice)
		//Sin panic mejor, capaz es mejor retornar un error?
		return
	}
	grafo.hashConexiones.Guardar(vertice, hash.CrearHash[V, T]())
}

func (grafo *grafoNumerico[V, T]) AgregarArista(desde V, hasta V, peso T) bool {

	// Se deberia chequear que existen, pero le dejamos que le explote a barbara, para algo esta el existe vertice.
	aristasDesde := grafo.hashConexiones.Obtener(desde)
	diffCantidad := 1
	if !aristasDesde.Guardar(hasta, peso) {
		diffCantidad = 0
	}

	if !grafo.dirigido {
		diffCantidad *= 2
		grafo.hashConexiones.Obtener(hasta).Guardar(desde, peso)
	}

	grafo.cantidadAristas += diffCantidad

	return diffCantidad != 0
}

func (grafo *grafoNumerico[V, T]) ExisteVertice(vertice V) bool {
	return grafo.hashConexiones.Pertenece(vertice)
}

func (grafo *grafoNumerico[V, T]) ExisteArista(desde, hasta V) bool {
	return grafo.hashConexiones.Obtener(desde).Pertenece(hasta)
}

func (grafo *grafoNumerico[V, T]) BorrarVertice(vertice V) {
	// Se deberia chequear que existe, pero le dejamos que le explote a barbara, para algo esta el existe vertice.
	grafo.hashConexiones.Borrar(vertice)

	grafo.hashConexiones.Iterar(func(_ V, conn hash.Diccionario[V, T]) bool {
		if conn.Pertenece(vertice) {
			conn.Borrar(vertice)
		}

		return true
	})
}

func (grafo *grafoNumerico[V, T]) BorrarArista(desde V, hasta V) {
	// Se deberia chequear que existen, pero le dejamos que le explote a barbara, para algo esta el existe vertice.

	ady := grafo.hashConexiones.Obtener(desde)

	if !ady.Pertenece(hasta) {
		fmt.Printf("\nWarning: se intento borrar una arista inexistente : %v->%v\n", desde, hasta)
		return
	}

	grafo.cantidadAristas--

	ady.Borrar(hasta)

	if !grafo.dirigido {
		grafo.cantidadAristas--
		grafo.hashConexiones.Obtener(hasta).Borrar(desde)
	}
}

func (grafo *grafoNumerico[V, T]) ObtenerPeso(desde, hasta V) T {
	// Se deberia chequear que existen, pero le dejamos que le explote a barbara, para algo esta el existe vertice.
	ady := grafo.hashConexiones.Obtener(desde)

	if !ady.Pertenece(hasta) {
		fmt.Printf("\nWarning: se intento obtener un peso de una arista inexistente : %v->%v\n", desde, hasta)
		return 0
	}

	return ady.Obtener(hasta)

}

func (grafo *grafoNumerico[V, T]) IterarAdyacentes(vertice V, visitar func(hasta V, peso T) bool) {
	// Se deberia chequear que existe, pero le dejamos que le explote a barbara, para algo esta el existe vertice.
	conn := grafo.hashConexiones.Obtener(vertice)
	conn.Iterar(func(hasta V, peso T) bool {
		return visitar(hasta, peso)
	})
}

func (grafo *grafoNumerico[V, T]) ObtenerAdyacentes(vertice V) []V {
	// Se deberia chequear que existe, pero le dejamos que le explote a barbara, para algo esta el existe vertice.
	conn := grafo.hashConexiones.Obtener(vertice)
	res := make([]V, conn.Cantidad())
	i := 0
	conn.Iterar(func(hasta V, _ T) bool {
		res[i] = hasta
		i++
		return true
	})

	return res

}

func (grafo *grafoNumerico[V, T]) ObtenerVertices() []V {
	vertices := make([]V, grafo.hashConexiones.Cantidad())

	i := 0
	grafo.IterarVertices(func(vert V) bool {
		vertices[i] = vert
		i++
		return true
	})

	return vertices
}

func (grafo *grafoNumerico[V, T]) IterarVertices(visitar func(vert V) bool) {
	grafo.hashConexiones.Iterar(func(vertice V, _ hash.Diccionario[V, T]) bool {
		return visitar(vertice)
	})
}

func (grafo *grafoNumerico[V, T]) IterarAristas(visitar func(desde V, hasta V, peso T) bool) {
	keepIterating := true
	grafo.hashConexiones.Iterar(func(desde V, aristas hash.Diccionario[V, T]) bool {
		aristas.Iterar(func(hasta V, peso T) bool {
			keepIterating = visitar(desde, hasta, peso)
			return keepIterating
		})
		return keepIterating
	})
}
