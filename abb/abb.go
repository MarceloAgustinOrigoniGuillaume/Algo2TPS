package diccionario

import hash "hash/interface" // esto se deberia quitar al entregar...

type nodoABB[K comparable, V any] struct{
	clave K
	valor V
	izq *nodoABB[K,V]
	der *nodoABB[K,V]
}

type abbStruct[K comparable, V any] struct{
	raiz *nodoABB[K,V]
	cantidad int
	cmp func(K,K) int
}

func CrearABB[K comparable, V any](comparador func(K,K) int) DiccionarioOrdenado[K,V]{
	res := new(abbStruct[K,V])
	res.cmp = comparador

	return res
}


// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func (abb *abbStruct[K,V]) Guardar(clave K, dato V){

}


// Pertenece determina si una clave ya se encuentra en el diccionario, o no
func (abb *abbStruct[K,V]) Pertenece(clave K) bool{
	return false
}



// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
// 'La clave no pertenece al diccionario'
func (abb *abbStruct[K,V]) Obtener(clave K) V{
	return abb.raiz.valor
}



// Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
// pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
func (abb *abbStruct[K,V]) Borrar(clave K) V{
	return abb.raiz.valor
}


// Cantidad devuelve la cantidad de elementos dentro del diccionario
func (abb *abbStruct[K,V]) Cantidad() int{
	return abb.cantidad
}


// ITERADORES


// Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del
// mismo
func (abb *abbStruct[K,V]) Iterar(func(clave K, dato V) bool){

}




// IterarRango itera sólo incluyendo a los elementos que se encuentren comprendidos en el rango indicado,
// incluyéndolos en caso de encontrarse
func (abb *abbStruct[K,V]) IterarRango(desde *K, hasta *K, visitar func(clave K, dato V) bool){

}


// IteradorRango crea un IterDiccionario que sólo itere por las claves que se encuentren en el rango indicado
func (abb *abbStruct[K,V]) IteradorRango(desde *K, hasta *K) hash.IterDiccionario[K, V]{
	return nil
}

// Iterador devuelve un IterDiccionario para este Diccionario
func (abb *abbStruct[K,V]) Iterador() hash.IterDiccionario[K, V]{
	return nil
}
