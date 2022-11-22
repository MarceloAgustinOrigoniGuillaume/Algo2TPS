package managers

import heap "cola_prioridad"
import "tp2/interfaces"

// Implementacion de un Mapa de conexiones de usuarios a posts
// En especifico la implementacion requerida para Algogram

type conexionUserPost[U any, D comparable] struct {
	nodo       U
	conexiones heap.ColaPrioridad[D]
}

func (conexion *conexionUserPost[U, D]) VerNodo() U {
	return conexion.nodo
}

func (conexion *conexionUserPost[U, D]) AgregarConexion(element D) {
	conexion.conexiones.Encolar(element)
}

func (conexion *conexionUserPost[U, D]) ObtenerConexion() *D {
	if conexion.conexiones.EstaVacia() {
		return nil
	}
	res := conexion.conexiones.Desencolar()
	return &(res)
}

func crearConexionUserPost[U any, D comparable](nodo U, comparadorMultiple func(D, D) int) interfaces.MapConexiones[U, D] {
	return &conexionUserPost[U, D]{nodo, heap.CrearHeap[D](comparadorMultiple)}
}
