package hash
const _CAPACIDAD_INICIAL = 128
const ERROR_FUNCION_HASH = "Error: mala funcion de hash, redimension dos veces seguidas no permitida"
const ERROR_NO_ESTABA = "Error: No estaba el elemento key"
func funcionHashing1() int{

}

func funcionHashing2() int{

}

func funcionHashing3() int{

}

type elementoHash[K comparable, V any]{
	clave K
	valor V
	indiceFuncion int
}

func crearElementoHash[K comparable, V any](clave K, valor V) elementoHash[K,V]{
	elemento:= new(elementoHash[K,V])

	elemento.clave = clave
	elemento.valor = valor

	elemento.funcionesHash = []func() int {funcionHashing1,funcionHashing2,funcionHashing3}

}

type hashCuckoo[K comparable, V any] struct{
	elementos []*elementoHash
	cantidad int
	funcionesHash []func() int
}


func CrearHashCuckoo[K comparable, V any]() Hash[K,V]{
	hash := new(hashCuckoo)

	elementos = make([]elementoHash,_CAPACIDAD_INICIAL)

}
func (hash *hashCuckoo) reemplazoCuckoo(nuevoElemento elementoHash,pos int) elementoHash{
	aGuardar:= hash.elementos[pos]
	aGuardar.indiceFuncion++
	if aGuardar.indiceFuncion == len(hash.funcionesHash){
		aGuardar.indiceFuncion = 0
	}

	hash.elementos[pos] = nuevoElemento
	return hash.funcionesHash[elementoHash.indiceFuncion](clave) % len(hash.elementos)
}

func (hash *hashCuckoo) insertarCuckoo(nuevoElemento elementoHash) bool{
	
	posicionando := hash.reemplazoCuckoo(nuevoElemento)

	for posicionando != nuevoElemento && posicionando != nil {
		posicionando = hash.reemplazoCuckoo(posicionando)
	}

	return posicionando == nil
}

func (hash *hashCuckoo) buscarPosicionCuckoo(elementos []elementoHash, clave K) int{
	for _,funcionHash := range hash.funcionesHash{
		indice := funcionHash(clave) % len(elementos)
		if(elementos[indice].clave == clave){
			return indice
		}
	}

	return -1
}

func (hash *hashCuckoo) redimensionar(nuevoLargo int){
	elementosNew := make([]elementoHash, nuevoLargo)
	hash.Iterar(func (clave K, valor V) {
		if(!hash.insertarCuckoo(elementosNew)){
			panic(ERROR_FUNCION_HASH) // no deberia pasar
		}
	})

	hash.elementos = elementosNew
}



func (hash *hashCuckoo) Guardar(clave K, valor V){
	indice := hash.buscarPosicionCuckoo(nuevo)


	if(indice == -1){
		nuevo:= crearElementoHash(clave,valor)

		if(!hash.insertarCuckoo(nuevo)){
			hash.redimensionar(2 * len(hash.elementos))

			if(!hash.insertarCuckoo(nuevo)){
				panic(ERROR_FUNCION_HASH) // no deberia pasar
			}
		}

		hash.cantidad++
	} else{
		hash.elementos[i].valor = valor
	}

}


func (hash *hashCuckoo) Pertenece(clave K) bool{
	return buscarPosicionCuckoo(hash.elementos,clave) != -1
}

func (hash *hashCuckoo) Obtener(clave K) V{
	i:= buscarPosicionCuckoo(hash.elementos,clave)
	if(i == -1){
		panic(ERROR_NO_ESTABA)
	}
	return hash.elementos[i]
	
}

func (hash *hashCuckoo) Borrar(clave K) V{
	i:= buscarPosicionCuckoo(hash.elementos,clave)
	
	if(i != -1){
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
