package lib

import grafos "tp3/grafos"
import cola "tp3/cola"
func OrdenTopologico[V comparable, T any](grafo grafos.Grafo[V,T],visitar func(visitado V) bool){
	if !grafo.EsDirigido(){
		panic("Orden topologico no es posible en un grafo no dirigido...")
	}

	grados_entrada := GradosDeEntrada(grafo)
	aVisitar := cola.CrearColaEnlazada[V]()	
	var visitado V

	grados_entrada.Iterar(func (vert V,grado int) bool {
		if grado == 0{
			aVisitar.Encolar(vert)
		}
		return true
	})

	for !aVisitar.EstaVacia(){
		visitado = aVisitar.Desencolar()
		if !visitar(visitado){
			return
		}


		grafo.IterarAdyacentes(visitado, func (ady V,_ T) bool{

			grados_entrada.Guardar(ady,grados_entrada.Obtener(ady)-1)
			if grados_entrada.Obtener(ady) == 0{
				aVisitar.Encolar(ady)
			} 

			return true
		})



	}	
}




func SecuenciaTopologica[V comparable, T any](grafo grafos.Grafo[V,T]) ([]V,error) {
	secuencia := make([]V ,grafo.CantidadVertices())
	i:= 0

	OrdenTopologico(grafo, func(visitado V) bool{
		secuencia[i] = visitado
		i++
		return true
	})

	if i != len(secuencia){
		return nil, CrearErrorGrafo("Warning: Secuencia topologica no contiene todos los vertices, habia bucles.")
	}

	return secuencia,nil
}


