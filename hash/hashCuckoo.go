package hash

func funcionHashing1() int{

}

func funcionHashing2() int{

}

func funcionHashing3() int{

}

const _CAPACIDAD_INICIAL = 128
type elementoHash[K comparable, V any]{
	clave K
	valor V
}

type hashCuckoo[K comparable, V any] struct{
	elementos []elementoHash
}


func CrearHashCuckoo[K comparable, V any]() Hash[K,V]{
	hash := new(hashCuckoo)

	elementos = make([]elementoHash,_CAPACIDAD_INICIAL)

}

func (hash *hashCuckoo) 

func (hash *hashCuckoo) redimensionar(nuevoLargo int){
	elementosNew := make([]elementoHash, nuevoLargo)

	copy(elementosNew, hash.elementos)
	hash.elementos = elementosNew
}



func (hash *hashCuckoo) Guardar(clave K, valor V){
	
}


func (hash *hashCuckoo) Pertenece(clave K) bool{
	
}

func (hash *hashCuckoo) Obtener(clave K) V{
	
}

func (hash *hashCuckoo) Borrar(clave K) V{
	
}

func (hash *hashCuckoo) Cantidad() int{
	
}
	
func (hash *hashCuckoo) Iterar(visitar func(clave K, dato V) bool){
	
}

func (hash *hashCuckoo) Iterador() IterDiccionario[K,V]{

}
