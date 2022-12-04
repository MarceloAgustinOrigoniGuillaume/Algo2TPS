package grafos

//import "golang.org/x/exp/constraints"

type Numero interface {
	int | int8 | int32 | int64 | float32 | float64
}

type Grafo[V any, T any] interface {
	EsPesado() bool
	EsDirigido() bool

	//Agrega un vertice, que hace si ya estaba dependera de la implementacion, algunos permitiran repetidos otros no.
	AgregarVertice(vertice V)

	//true si existia el vertice, false si no
	ExisteVertice(vertice V) bool

	//true si existia la arista de desde a hasta, false si no
	ExisteArista(desde, hasta V) bool

	//Borra el vertice
	BorrarVertice(vertice V)

	// Agrega la arista, devuelve si true si agrego una arista, false si ya estaba
	AgregarArista(desde V, hasta V, peso T) bool

	// Borra la arista
	BorrarArista(desde V, hasta V)

	//Obtiene un peso
	ObtenerPeso(desde V, hasta V) T

	// Obtiene todos los vertices
	ObtenerVertices() []V

	// Obtiene todos los adyacentes al vertice
	ObtenerAdyacentes(vertice V) []V

	// Devuelve la cantidad total de vertices
	CantidadVertices() int

	// Devuelve la cantidad total de aristas, en grafos no dirigidos, seria el doble de la cantidad efectiva
	CantidadAristas() int

	// Iteradores
	IterarVertices(func(V) bool)

	IterarAristas(func(desde V, hasta V, peso T) bool)
	IterarAdyacentes(vertice V, visitar func(hasta V, peso T) bool)
}
