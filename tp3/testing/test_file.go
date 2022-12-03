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

	_,_ = generators.ExecuteTest(os.Args[1])

}
