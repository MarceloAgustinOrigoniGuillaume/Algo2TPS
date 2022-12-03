package main
import "fmt"
import (
	"tp3/grafos"
	"strconv"
	"math/rand"
	"time"
	"os"
	pj "tp3/utils/PJ"
	"tp3/testing/generators"
	grafosLib "tp3/grafos/lib"
	pila "tp3/pila"
)


func showGrados[T any](grafo grafos.Grafo[string,T]){
	entradas := grafosLib.GradosDeEntrada(grafo)

	grafosLib.GradosDeSalida(grafo).Iterar(func (persona string,grados int) bool{
		fmt.Printf("%s ama a %d personas y es amadx por %d personas \n",persona,grados,entradas.Obtener(persona))
		return true
	})
}

func main(){

	if len(os.Args) < 3{
		fmt.Printf("Insuficientes argumentos se requieren 3, la ruta del pj y la cantidad minima de vertices y aristas ")
		return
	}

	
	
	var cantVert int 
	var cantEdges int
	var err error

	cantVert,err = strconv.Atoi(os.Args[2])

	if err != nil{
		fmt.Printf("Error : %s",err.Error())
		return
	}

	cantEdges,err = strconv.Atoi(os.Args[3])

	if err != nil{
		fmt.Printf("Error : %s",err.Error())
		return
	}
	fmt.Printf("GRAPH WITH CANT VERT %d, edges :%d\n",cantVert,cantEdges)

	labels := []string{"Makise","Inaba","Misaka","Miu","Asuka","Etc"}
	dirigido := false


	//grafo := generators.BuildGrafoTopologico[int](dirigido,labels,cantVert)
	
	grafo := generators.BuildVertices[int](dirigido,labels,cantVert, haceNada,haceNada)
	desde,hasta,pesoTotal := generators.BuildGrafoEuler(time.Now().UnixNano(),grafo,cantEdges,func (edge int) int {return rand.Intn(edge+1)},generators.TIPO_CICLO_EULERIANO)
	
	showGrados(grafo)

	fmt.Printf("Hay camino desde Desde %v, hasta %v, pesoTotal %v, cantidad aristas: %d\n",desde,hasta,pesoTotal,grafo.CantidadAristas())

	desdeValor := "Makise"

	if desde != nil{
		fmt.Printf("Hay camino desde Desde %v, hasta %v\n",(*desde),(*hasta))
		desdeValor = *desde
	}

	camino,errC := grafosLib.CaminoEuleriano(grafo,desdeValor)

	if errC != nil{
		fmt.Printf("Err : %s\n",errC.Error())
		err = writePj(os.Args[1],grafo)
		return 
	}



	tiempoTotal := 0
	anterior := camino[0]
	for _,elem := range camino[1:]{
		tiempoTotal += grafo.ObtenerPeso(anterior,elem)
		anterior = elem
	}

	fmt.Printf("Camino euleriano start: %s end: %s pesoTotal: %v\n",camino[0],camino[len(camino)-1],tiempoTotal)

	//buildAristasAzar(grafo,cantEdges,func (edge int) int {return edge})


	err = writePj(os.Args[1],grafo)

	if err != nil{
		fmt.Printf("\nError : %s",err.Error())
	} else{
		fmt.Printf("\nsaved to %s",os.Args[1])		
	}
}