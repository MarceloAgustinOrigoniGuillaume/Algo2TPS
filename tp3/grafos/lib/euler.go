package lib

import grafos "tp3/grafos"
import pila "tp3/pila"
import hash "tp3/diccionario"
import "fmt"


func clavesHash[K comparable, V any](dicc hash.Diccionario[K,V]) []K{
		res := make([]K, dicc.Cantidad())
		i:= 0
		dicc.Iterar(func (key K, value V ) bool{
			res[i] = key
			i++
			return true
		})	

		return res
}
// Recorrido per se

func _dfsAristasNoDirigido[V comparable, T any](aristasAVisitar hash.Diccionario[V,hash.Diccionario[V,T]],aristas hash.Diccionario[V,T],origen V, visitar func(visitado V)){

	var aristasNext hash.Diccionario[V,T]
	// se toma las claves del hash en vez de solo iterar para hacer una "buena practica" y evitar problemas
	// que se podrian generar por la forma en la que itera el hash
	
	for _,ady := range clavesHash(aristas){ 

		// Si no pertenece ya se visito
		// Y como al visitar se borra del hash del origen como del de llegada tampoco hace falta borrarla de currAvisitar
		aristasNext = aristasAVisitar.Obtener(ady)
		if aristasNext.Pertenece(origen) {
			// se asume solo hay una arista conectando dos vertices, esto es para no dirigidos	
			aristas.Borrar(ady)
			aristasNext.Borrar(origen)

			_dfsAristasNoDirigido(aristasAVisitar,aristasNext,ady,visitar) // Visitas a partir de este
		}
	}

	// Esta asegurado que no le queda a quien visitar
	// Visita este vertice, dado que es dfs primero se visitaria el ultimo en el camino
	visitar(origen) 
}


func _dfsAristasDirigido[V comparable, T any](aristasAVisitar hash.Diccionario[V,hash.Diccionario[V,T]],origen V, visitar func(visitado V)){

	aristas :=aristasAVisitar.Obtener(origen)

	// se toma las claves del hash en vez de solo iterar para hacer una "buena practica" y evitar problemas
	// que se podrian generar por la forma en la que itera el hash
	for _,ady := range clavesHash(aristas){ 

		if aristas.Pertenece(ady) { // se pudo haber visto en un llamado recursivo
			// se borra para no repetir	
			aristas.Borrar(ady)
			_dfsAristasDirigido(aristasAVisitar,ady,visitar) // Visitas a partir de este
		}
	}

	// Esta asegurado que no le queda a quien visitar
	// Visita este vertice, dado que es dfs primero se visitaria el ultimo en el camino
	visitar(origen) 
}



func recorridoEulerDirigido[V comparable, T any](grafo grafos.Grafo[V,T],origen V,aceptarSemi bool,visitar func(visitado V)) bool{

	aristasAVisitar,grados_entrada := AristasDeSalidaYGradosEntrada(grafo)
		
	// Verificacion tiene euleriano
	// y con start en origen
	var end *V
	var diff int
		
	aristasAVisitar.Iterar(func (vert V, salidas hash.Diccionario[V,T]) bool{
		diff = grados_entrada.Obtener(vert)-salidas.Cantidad()

		if diff != 0{ // cantidad distinta

			if (diff < -1 || diff >1){
				end = nil // para saber despues que fallo
				return false // no es ni semiEulariano
			}

			if diff < 0{

				if(vert != origen){
					end = nil // para saber despues que fallo
					return false // ya no podria ser el start el origen					
				}

				return true
			}

			if end != nil{
				end = nil
				return false // no es ni semi euleriano
			}
		
			end = &vert
		}

			return true
		})


	if ((diff != 0 && end == nil) || (end != nil && !aceptarSemi)){ // no se pudo o no se acepta
		//fmt.Printf("ni intentes el camino euleriano, dirigido\n")
		return false
	}
	//fmt.Printf("intenta el camino euleriano, dirigido\n")


	//recorrido
	_dfsAristasDirigido(aristasAVisitar,origen,visitar)

	return end == nil
}

func recorridoEulerNoDirigido[V comparable, T any](grafo grafos.Grafo[V,T],origen V,aceptarSemi bool,visitar func(visitado V)) bool{

	// No dirigidos tienen un manejo algo distinto
	cantidadImpares := 0
	aristasAVisitar := AristasDeSalida(grafo)

	// Verificacion existe y con start en origen
	var end *V
	
	aristasAVisitar.Iterar(func (vert V, salidas hash.Diccionario[V,T]) bool{
		if (salidas.Cantidad() % 2) != 0{ // cantidad impar de salidas
			cantidadImpares++
			if !aceptarSemi{
				return false
			}

			if vert == origen{
				return true
			}

			if end == nil{
				end = &vert
			} else if end != nil{
				// podria ser 2 o 3 este contador, sumas uno para que sea 3 o 4, siempre rompiendo
				cantidadImpares++ 
				// no puede haber mas de dos con %2 para ser semi euleriano
				// Y se toma al origen como el start para este caso en particular
				return false 
			}
		}

		return true
	})

	// si hay distinto de 2 impares no hay , o si no se acepta semi paras aca
	if (cantidadImpares != 0 && (!aceptarSemi || cantidadImpares != 2)){ 
		////fmt.Printf("ni intentes el camino euleriano, no dirigido\n")
		return false
	}
	////fmt.Printf("intenta , no dirigido\n")

	_dfsAristasNoDirigido(aristasAVisitar,aristasAVisitar.Obtener(origen),origen,visitar)
	
	return cantidadImpares == 0 // devuelve si fue euleriano
}

func RecorridoEuleriano[V comparable, T any](grafo grafos.Grafo[V,T],origen V,aceptarSemi bool,visitar func(visitado V)) bool{
	
	if grafo.EsDirigido() {
		////fmt.Printf("euler dirigido\n")

		return recorridoEulerDirigido(grafo,origen,aceptarSemi,visitar)
	} 
	////fmt.Printf("euler no dirigido\n")

	return recorridoEulerNoDirigido(grafo,origen,aceptarSemi,visitar)
}



// ciclos / caminos  de euler


func CicloEuleriano[V comparable, T any](grafo grafos.Grafo[V,T],origen V) ([]V,error){

	if !grafo.ExisteVertice(origen){
		return nil,CrearErrorGrafo("Vertice no existia al grafo")
	}
	camino := pila.CrearPilaDinamica[V]()
	longitud := 0
	RecorridoEuleriano(grafo,origen,false, func(visitado V) {
		camino.Apilar(visitado)
		////fmt.Printf("CicloEuleriano visitando -->%v\n",visitado)
		longitud++
	})

	longitudCamino := grafo.CantidadAristas()
	if !grafo.EsDirigido(){
		longitudCamino = longitudCamino/2
	}

	if longitud != longitudCamino+1{
		return nil,CrearErrorGrafo("El grafo NO tenia ciclo euleriano")
	}

	res := make([]V,longitud)

	for i:=0;i<longitud;i++{
		res[i] = camino.Desapilar()
	}


	return res,nil

}



func CaminoEuleriano[V comparable, T any](grafo grafos.Grafo[V,T],origen V) ([]V,error){
	if !grafo.ExisteVertice(origen){
		return nil,CrearErrorGrafo("Vertice no existia al grafo")
	}

	camino := pila.CrearPilaDinamica[V]()
	longitud := 0
	RecorridoEuleriano(grafo,origen,true, func(visitado V) {
		camino.Apilar(visitado)
		////fmt.Printf("CaminoEuleriano visitando -->%v\n",visitado)
		longitud++
	})

	longitudCamino := grafo.CantidadAristas()
	if !grafo.EsDirigido(){
		longitudCamino = longitudCamino/2
	}

	if longitud != longitudCamino+1{

		//fmt.Printf("LONG SHOULD BE %d is %d\n",longitudCamino+1,longitud)
		return nil,CrearErrorGrafo(fmt.Sprintf("El grafo NO tenia un camino euleriano que empieze en %v, %d",origen,longitud))
	}

	res := make([]V,longitud)

	for i:=0;i<longitud;i++{
		res[i] = camino.Desapilar()
	}


	return res,nil

}
