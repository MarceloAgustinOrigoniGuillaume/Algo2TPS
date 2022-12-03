package lib
import "tp3/grafos"

type PairDistancia[V any,T grafos.Numero] struct{
	desde *V
	visitado V
	distancia T
}

func(pair *PairDistancia[V,T]) Desde() *V{
	return pair.desde
}

func(pair *PairDistancia[V,T]) Actual() V{
	return pair.visitado
}

func(pair *PairDistancia[V,T]) Distancia() T{
	return pair.distancia
}

func(pair *PairDistancia[V,T]) Add(added T) {
	pair.distancia += added
}




func CrearPairDistancia[V comparable,T grafos.Numero](desde *V,hasta V, distancia T) PairDistancia[V,T]{
	return PairDistancia[V,T]{desde,hasta,distancia}
}

func comparadorDistancias[V any,T grafos.Numero](dist1,dist2 PairDistancia[V,T]) int{
	if dist1.distancia>dist2.distancia{
		return -1
	}

	if dist1.distancia<dist2.distancia{
		return 1
	}
	return 0
}