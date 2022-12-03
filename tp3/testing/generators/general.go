package generators

import "tp3/grafos"
import grafosLib "tp3/grafos/lib"
import "fmt"
import "math/rand"
import "strconv"



func showGrados[T any](grafo grafos.Grafo[string,T]){
	entradas := grafosLib.GradosDeEntrada(grafo)
	differentes := 0
	salidas := 0
	ceroes := 0
	impares:= 0
	grafosLib.GradosDeSalida(grafo).Iterar(func (persona string,grados int) bool{
		//fmt.Printf("%s ama a %d personas y es amadx por %d personas \n",persona,grados,entradas.Obtener(persona))

		if (grados != entradas.Obtener(persona)){
			differentes++
		}

		if (grados > entradas.Obtener(persona)){
			salidas++
		}

		if (grados == entradas.Obtener(persona) && grados == 0){
			ceroes++
		}

		if (grados %2==1){
			impares++
		}



		return true
	})

	fmt.Printf("Habia %d diferentes grados, %d serian salidas, %d ceroes, impares %d\n",differentes,salidas,ceroes,impares)
}


func CopyVerts[V comparable, T grafos.Numero](grafo grafos.Grafo[V,T],dirigido bool) grafos.Grafo[V,T]{
	copy := grafos.GrafoNumericoPesado[V,T](dirigido)

	grafo.IterarVertices(func (vert V) bool{
		copy.AgregarVertice(vert)
		return true
	})

	return copy
}



func BuildVertices[T grafos.Numero](dirigido bool, labels []string,vert_quantity int,vertAgregado func(vert string),sufijoTerminado func(sufijo string)) grafos.Grafo[string,T] {
	grafo := grafos.GrafoNumericoPesado[string,T](dirigido)
	
	index := 0
	suffixCount := 0
	suffix := ""
	vert := ""
	for grafo.CantidadVertices() < vert_quantity{

		vert = labels[suffixCount] +suffix
		vertAgregado(vert)
		grafo.AgregarVertice(vert)

		suffixCount++ 
		if suffixCount == len(labels){
			sufijoTerminado(suffix)
			index += 1
			suffix = strconv.Itoa(index) 
			suffixCount = 0			
		}
	}
	return grafo
}



func HaceNada(_ string){

}



func BuildAristasAzar[T grafos.Numero](seed int64, grafo grafos.Grafo[string,T], edge_quantity int, pesoProvider func(int) T){

	rand.Seed(seed)
	//rand.Seed(time.Now().UnixNano())
	vertices := grafo.ObtenerVertices()
	
	var selectedDesde int
	var selectedHasta int
	
	length := len(vertices)

	for grafo.CantidadAristas() < edge_quantity{

		selectedDesde = rand.Intn(length)
		selectedHasta = rand.Intn(length)

		for selectedDesde == selectedHasta{
			selectedHasta = rand.Intn(length)
		}

		grafo.AgregarArista(vertices[selectedDesde],vertices[selectedHasta],pesoProvider(grafo.CantidadAristas()))
	}
}
