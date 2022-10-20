package diccionario

import "pila"
import hash "hash/interface" // esto se deberia quitar al entregar...

const(
	ERROR_NO_ESTABA               = "La clave no pertenece al diccionario"
	ERROR_ITERADOR_TERMINO        = "El iterador termino de iterar"
)

type nodoABB[K comparable, V any] struct{
	clave K
	valor V
	izq *nodoABB[K,V]
	der *nodoABB[K,V]
}
func CrearNodoABB[K comparable, V any](clave K, valor V) *nodoABB[K,V]{
    nodo:= new(nodoABB[K,V])
    nodo.clave = clave
    nodo.valor = valor
    return nodo
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



func buscarNodo[K comparable, V any](raiz **nodoABB[K,V],clave K, comparacion func(K,K) int) **nodoABB[K,V] {
    if *raiz == nil{
        return raiz
    }
    
    res:= comparacion(clave,(*raiz).clave)
    if(res == 0){
        return raiz
    }
    
    var nodo **nodoABB[K,V]

    if(res > 0){
        nodo = buscarNodo(&((*raiz).der),clave,comparacion)
    } else{
        nodo = buscarNodo(&((*raiz).izq),clave,comparacion)
    } 

    return nodo
}

// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func (abb *abbStruct[K,V]) Guardar(clave K, dato V){
    aguardar:=buscarNodo(&abb.raiz,clave,abb.cmp)
    if *aguardar==nil{
        abb.cantidad++
        *aguardar = CrearNodoABB[K,V](clave,dato)
    }else{
        (*aguardar).valor = dato
    }
}

func pertenece[K comparable, V any](nodo *nodoABB[K,V]) bool{
    return nodo!=nil
}


// Pertenece determina si una clave ya se encuentra en el diccionario, o no
func (abb *abbStruct[K,V]) Pertenece(clave K) bool{
	return pertenece(*buscarNodo(&abb.raiz,clave,abb.cmp))
}

// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
// 'La clave no pertenece al diccionario'
func (abb *abbStruct[K,V]) panicNoEstaba(nodo **nodoABB[K,V]){
    if(!pertenece(*nodo)){
        panic(ERROR_NO_ESTABA)
    }
}

func (abb *abbStruct[K,V]) Obtener(clave K) V{ 
    nodoEncontrado := buscarNodo(&abb.raiz,clave,abb.cmp)    
    abb.panicNoEstaba(nodoEncontrado)
    return (*nodoEncontrado).valor
}



// Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
// pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'

func (original *nodoABB[K,V]) copy(aCopiar *nodoABB[K,V]){
    original.clave = aCopiar.izq.der.clave
    original.valor = aCopiar.izq.der.valor
}

func swapBorrar[K comparable, V any](borrado **nodoABB[K,V]){
    nodoBorrado := *borrado

    if nodoBorrado == nil{
        return
    }

    if nodoBorrado.izq == nil && nodoBorrado.der == nil{
        (*borrado) = nil
    } else if nodoBorrado.izq == nil{
        (*borrado) = nodoBorrado.der
    } else if nodoBorrado.der == nil {
        (*borrado) = nodoBorrado.izq
    } else{
        var siguienteAborrar *nodoABB[K,V]
        // logica compleja dea, por ahora sin abl
        if nodoBorrado.izq.der != nil {
            siguienteAborrar = nodoBorrado.izq.der
            for siguienteAborrar.der != nil{
                siguienteAborrar = siguienteAborrar.der
            }

        } else {
            siguienteAborrar = nodoBorrado.izq
        }

        nodoBorrado.valor = siguienteAborrar.valor
        swapBorrar(&siguienteAborrar)
    }
}

func (abb *abbStruct[K,V]) Borrar(clave K) V{
    nodoBorrar := buscarNodo(&abb.raiz,clave,abb.cmp)  
    abb.panicNoEstaba(nodoBorrar)    

    res := (*nodoBorrar).valor
    swapBorrar(nodoBorrar)
    abb.cantidad--
    return res
}


// Cantidad devuelve la cantidad de elementos dentro del diccionario
func (abb *abbStruct[K,V]) Cantidad() int{
	return abb.cantidad
}


// ITERADORES


// Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del
// mismo

 
func iterarRango[K comparable, V any](nodo *nodoABB[K,V],desde *K, hasta *K,visitar func(*nodoABB[K,V]) bool,
                                        cmp   func(K,K) int ) bool{
    if nodo == nil{  
        return true
    }

    diffHasta := -1

    if hasta != nil {
        diffHasta = cmp(nodo.clave, *hasta)
    }
    diffDesde := 1
    if desde != nil{
        diffDesde = cmp(nodo.clave, *desde)
    }
     
    
    continuar := true
    if diffDesde >0 {// es mayor que el desde
        continuar = iterarRango(nodo.izq,visitar,desde,hasta,cmp)
    }

    if (continuar && diffDesde >=0 && diffHasta<=0){//esta en el rango
        continuar = visitar(nodo)
    }

    if (continuar && diffHasta<0 ){// es menor que el hasta
        continuar = iterarRango(nodo.der,visitar,desde,hasta,cmp)
    }

    return continuar
}
func (abb *abbStruct[K,V]) Iterar(visitar func(clave K, dato V) bool){

    iterarRango(abb.raiz,nil,nil,visitar, abb.cmp)    

}




// IterarRango itera sólo incluyendo a los elementos que se encuentren comprendidos en el rango indicado,
// incluyéndolos en caso de encontrarse
func (abb *abbStruct[K,V]) IterarRango(desde *K, hasta *K, visitar func(clave K, dato V) bool){
    iterarRango(abb.raiz,desde,hasta,visitar)    
}

// IteradorRango crea un IterDiccionario que sólo itere por las claves que se encuentren en el rango indicado
func (abb *abbStruct[K,V]) IteradorRango(desde *K, hasta *K) hash.IterDiccionario[K, V]{

    

	return crearIteradorABB(abb.raiz,desde,hasta,abb.cmp)
}

// Iterador devuelve un IterDiccionario para este Diccionario
func (abb *abbStruct[K,V]) Iterador() hash.IterDiccionario[K, V]{
	return crearIteradorABB(abb.raiz, abb.cmp, nil, nil)
}

type iteradorABB[K comparable, V any] struct{
    aVisitar pila.Pila[*nodoABB[K,V]]
    cmp func (K,K) int
    desde *K
    hasta *K
}

func (iterador *iteradorABB[K,V]) agregarPrimeros(nodo *nodoABB[K,V]){

    if(desde != nil){
        for nodo != nil && iterador.cmp(nodo.clave, *iterador.desde)<0 {
            nodo = nodo.der
        }
    }
    

    if nodo != nil{
        iterador.agregarNodos(nodo)
    } 

}

func (iterador *iteradorABB[K,V]) agregarNodos(nodo *nodoABB[K,V]){
    diffHasta := -1

    if iterador.hasta != nil {
        diffHasta = iterador.cmp(nodo.clave, *iterador.hasta)
    }
    diffDesde := 1
    if iterador.desde != nil{
        diffDesde = iterador.cmp(nodo.clave, *iterador.desde)
    }
     

    for nodo!=nil && diffDesde>=0{
        if diffHasta<=0 { 
            iterador.aVisitar.Apilar(nodo)
        }
        nodo = nodo.izq
        diffHasta = iterador.cmp(nodo.clave, *iterador.hasta)
        diffDesde = iterador.cmp(nodo.clave, *iterador.desde)
    }
}

func crearIteradorABB[K comparable, V any](nodo *nodoABB[K,V],desde *K, hasta *K,
                                            cmp func (K,K) int, )
                                            hash.IterDiccionario[K,V]{

 
    iterador := new(iteradorABB[K,V])
    iterador.aVisitar = pila.CrearPilaDinamica[*nodoABB[K,V]]()
    iterador.cmp=cmp
    iterador.desde= desde
    iterador.hasta=hasta
    
    iterador.agregarNodos(nodo)

    return iterador
}


func (iterador *iteradorABB[K,V]) HaySiguiente() bool{
    return !iterador.aVisitar.EstaVacia()
}

func (iterador *iteradorABB[K,V]) panicTermino(){
    if(!iterador.HaySiguiente()){
        panic(ERROR_ITERADOR_TERMINO)
    }

}


func (iterador *iteradorABB[K,V]) VerActual() (K,V){
    iterador.panicTermino()

    actual:= iterador.aVisitar.VerTope()

    return actual.clave,actual.valor    
}

func (iterador *iteradorABB[K,V]) Siguiente() K{
    iterador.panicTermino()

    actual:= iterador.aVisitar.Desapilar()
    
    iterador.agregarNodos(actual.der)    

    return actual.clave
}

