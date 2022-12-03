package lib
import "tp3/grafos"
type Arista[V any, T grafos.Numero] interface{
	Desde() V
	Hasta() V
	Peso() T
	Sum(Arista[V,T])
	Compare(Arista[V,T]) int
}



type AristaSimple[V any, T grafos.Numero] struct{
	desde V
	hasta V
	peso T
}

func (arista *AristaSimple[V,T]) Desde() V{
	return arista.desde
}

func (arista *AristaSimple[V,T]) Hasta() V{
	return arista.hasta
}

func (arista *AristaSimple[V,T]) Peso() T{
	return arista.peso
}

func (arista *AristaSimple[V,T]) Sum(otra Arista[V,T]) {
	arista.peso+= otra.Peso()
}

func (arista *AristaSimple[V,T]) Compare(otra Arista[V,T]) int{
	peso_otra := otra.Peso()
	if arista.peso == peso_otra{
		return 0
	}

	diff := (int)(arista.peso - peso_otra)

	if diff != 0{
		return diff
	}

	if arista.peso > peso_otra{
		return 1
	}

	return -1
}

func CrearAristaSimple[V any,T grafos.Numero](desde V,hasta V,peso T) Arista[V,T]{
	return &AristaSimple[V,T]{desde,hasta,peso}
}