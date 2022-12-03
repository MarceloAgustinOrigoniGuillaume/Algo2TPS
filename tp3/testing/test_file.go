package main

import (
	"fmt"
	"os"
	"tp3/testing/generators"
	"tp3/grafos"
	//"tp3/utils/PJ"
	//"strconv"
	//"strings"
	//"errors"
	//"time"
	//"math/rand"
	grafosLib "tp3/grafos/lib"
)


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


func main(){
	if len(os.Args) < 2{
		fmt.Printf("Insuficientes argumentos se requiere la ruta del test")
		return
	}

	_,_ = generators.ExecuteEulerTest(os.Args[1])

}

/*
func other(){
	desdeValor:= os.Args[2]

	dirigido := false
	if len(os.Args) >= 4 && os.Args[3] == "true"{
		dirigido = true
	}


	grafo := grafos.GrafoNumericoPesado[string,int](dirigido)
	fmt.Printf("Grafo es dirigido? %v\n",grafo.EsDirigido())
	pj.LeerPJ(os.Args[1],grafo,generators.BasicRead,strconv.Atoi)

	fmt.Printf("Cantidad vertices vert ? %v\n",grafo.CantidadVertices())

	cant := 0


	grafosLib.DFS(grafo,desdeValor,func (_ string){cant++})

	fmt.Printf("From DFS? %v\n",cant)


	showGrados(grafo)


	camino,errC := grafosLib.CicloEuleriano(grafo,desdeValor)

	if errC != nil{
		fmt.Printf("Err : %s\n",errC.Error())
	} else{
			tiempoTotal := 0
			anterior := camino[0]
			for _,elem := range camino[1:]{
				tiempoTotal += grafo.ObtenerPeso(anterior,elem)
				anterior = elem
			}
			fmt.Printf("Ciclo euleriano start: %s end: %s pesoTotal: %v, len %d\n",camino[0],camino[len(camino)-1],tiempoTotal,len(camino))

	}


	camino,errC = grafosLib.CaminoEuleriano(grafo,desdeValor)

	if errC != nil{
		fmt.Printf("Err : %s\n",errC.Error())
	} else{
			tiempoTotal := 0 
			anterior := camino[0]
			for _,elem := range camino[1:]{
				tiempoTotal += grafo.ObtenerPeso(anterior,elem)
				anterior = elem
			}

			fmt.Printf("CaminoEuleriano euleriano start: %s end: %s pesoTotal: %v, len %d\n",camino[0],camino[len(camino)-1],tiempoTotal,len(camino))

	}
}



*/