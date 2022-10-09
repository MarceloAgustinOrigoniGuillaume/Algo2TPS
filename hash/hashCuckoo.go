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

func crearElementoHash[K comparable, V any](clave K, valor V) elementoHash[K,V]{
	elemento:= new(elementoHash[K,V])

	elemento.clave = clave
	elemento.valor = valor

	return elemento
}

type hashCuckoo[K comparable, V any] struct{
	elementos []*elementoHash
	cantidad int
}


func CrearHashCuckoo[K comparable, V any]() Hash[K,V]{
	hash := new(hashCuckoo)

	elementos = make([]elementoHash,_CAPACIDAD_INICIAL)

	return hash
}

func buscarPosicionCuckoo(elementos []elementoHash, clave K) int{
	secuencia := []func() int {funcionHashing1,funcionHashing2,funcionHashing3}
	for _,funcionHash := range secuencia{
		
		//recordar agregar chequeos de borrado
		valor := funcionHash(clave) % len(elementos)
		if(elementos[valor] == nil || elementos[valor].clave == clave){
			return valor 
		}
	}

	return -1
}
func posicionarElementoCuckoo(elementos []elementoHash, elemento elementoHash) bool{
	pos := buscarPosicionCuckoo(elementos,elemento.clave)
	if(pos == -1){
		return false
	}

	elementos[pos] = elemento

	return true
}

func (hash *hashCuckoo) redimensionar(nuevoLargo int){
	elementosNew := make([]elementoHash, nuevoLargo)
	hash.Iterar(func (clave K, valor V) {
		posicionarElementoCuckoo(elementosNew,CrearHashCuckoo(clave,valor))
	})

	hash.elementos = elementosNew
}



func (hash *hashCuckoo) Guardar(clave K, valor V){
	nuevo:= crearElementoHash(clave,valor)
	if(!posicionarElementoCuckoo(hash.elementos,nuevo)){
		hash.redimensionar(2 * len(hash.elementos))
		
		if(!posicionarElementoCuckoo(hash.elementos,nuevo)){
			return
		}
	}

	hash.cantidad++
}


func (hash *hashCuckoo) Pertenece(clave K) bool{
	i:= buscarPosicionCuckoo(hash.elementos,clave)
	return i != -1 && hash.elementos[i] != nil
}

func (hash *hashCuckoo) Obtener(clave K) V{
	i:= buscarPosicionCuckoo(hash.elementos,clave)
	if(i != -1 && hash.elementos[i] != nil){
		return hash.elementos[i].valor
	}
	return nil
	
}

func (hash *hashCuckoo) Borrar(clave K) V{
	i:= buscarPosicionCuckoo(hash.elementos,clave)
	
	if(i != -1 && hash.elementos[i] != nil){
		elem := hash.elementos[i].valor
		hash.elementos[i] = nil
		hash.cantidad--
		return elem
	}
}

func (hash *hashCuckoo) Cantidad() int{
	return hash.cantidad
}
	
func (hash *hashCuckoo) Iterar(visitar func(clave K, dato V) bool){
	i:= 0
	for (i<len(hash.elementos) && hash.elementos[i] == nil && visitar(hash.elementos[i].clave, hash.elementos[i].valor)){
		i++
	}	
}

func (hash *hashCuckoo) Iterador() IterDiccionario[K,V]{

}
