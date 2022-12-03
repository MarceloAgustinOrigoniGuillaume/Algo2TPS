package main

import (
	"fmt"
	"os"
	"tp3/testing/generators"
	"tp3/grafos"
	"strconv"
	//"strings"
	//"errors"
	"time"
	"math/rand"
	grafosLib "tp3/grafos/lib"
)


func showGrados[T any](grafo grafos.Grafo[string,T]){
	entradas := grafosLib.GradosDeEntrada(grafo)

	grafosLib.GradosDeSalida(grafo).Iterar(func (persona string,grados int) bool{
		fmt.Printf("%s ama a %d personas y es amadx por %d personas \n",persona,grados,entradas.Obtener(persona))
		return true
	})
}





func generarTestEuler(base_file string,grafo grafos.Grafo[string,int],cantEdges int,desdeValor string){

	var copy grafos.Grafo[string,int]
	
	values := []int{generators.TIPO_CICLO_EULERIANO,generators.TIPO_CAMINO_EULERIANO,generators.TIPO_NADA}
	labels_types := []string{"ciclo_euler","camino_euler","nada"}

	var pj_url string
	var url_test string
	var errWriting error
	for ind_value,curr_value := range values{
		//fmt.Printf("Generating Euler->%s\n",labels_types[ind_value])
		copy = generators.CopyVerts(grafo,grafo.EsDirigido())
		desde,hasta,pesoTotal := generators.BuildGrafoEuler(time.Now().UnixNano(),copy,cantEdges,func (edge int) int {return rand.Intn(edge+1)},curr_value)

		tipo_test := generators.CorroborarTipoTestEuler(copy,curr_value) // Por problemas al generar el test

		if tipo_test == generators.TIPO_NADA{
			desde = nil
			hasta = nil
		}
			
		pj_url = base_file+labels_types[ind_value]
		url_test = pj_url+".test"
		pj_url+= ".pj"
		fmt.Printf("------------>Writing pj: %s",pj_url)
			
		errWriting = generators.WritePj(pj_url,copy)

		if errWriting != nil{
			fmt.Printf("\nErr Writing test PJ: %s\n",errWriting.Error())
			continue
		}
		fmt.Printf(" test: %s\n",url_test)
		errWriting = generators.WriteEuler(url_test,pj_url, "test_"+labels_types[ind_value],grafo.EsDirigido(),desde,hasta,pesoTotal,tipo_test)

		if errWriting!= nil{
			fmt.Printf("Err Writing test : %s\n",errWriting.Error())
			continue
		}
		
	}
}


func generateTests(directory string, cantidad_iter int, generator func(base_file string)){
	var timeInit int64

	for i:= 0;i<cantidad_iter;i++{
		timeInit = time.Now().UnixNano()

		generator(directory+strconv.Itoa(i)+".")

		fmt.Printf("Finished test: %d...Took %d ns\n",i+1,(time.Now().UnixNano() - timeInit))

	}
}


func generatorDeVertices(labels []string,dirigido bool,cantidad_vertices int, generator func(string,grafos.Grafo[string,int])) func(string){
	var grafo grafos.Grafo[string,int]

	return func(base_file string) {
		grafo = generators.BuildVertices[int](dirigido,labels,cantidad_vertices, generators.HaceNada,generators.HaceNada)
		generator(base_file,grafo)
	}

}



func generarTestVolumen(base_file string,grafo grafos.Grafo[string,int],cantEdges int){

	generators.BuildAristasAzar(time.Now().UnixNano(),grafo,cantEdges, func (edge int) int {return rand.Intn(edge+1)})

	base_file = base_file+generators.TEST_VOLUMEN
	pj_url := base_file+".pj"
	url_test := base_file+".test"
	fmt.Printf("------------>Writing pj: %s",pj_url)
			
	errWriting := generators.WritePj(pj_url,grafo)

	if errWriting != nil{
		fmt.Printf("\nErr Writing test PJ: %s\n",errWriting.Error())
		return
	}
	fmt.Printf(" test: %s\n",url_test)
	errWriting = generators.WriteVolumen(url_test,pj_url, "test_"+generators.TEST_VOLUMEN,grafo.EsDirigido())

	if errWriting!= nil{
		fmt.Printf("Err Writing test : %s\n",errWriting.Error())
		return
	}
}


func main(){
	if len(os.Args) < 4{
		fmt.Printf("Insuficientes argumentos se requieren minimo 3,<Tipo de test a generar>,<directory tests>,<cantidad de tests>")
		return
	}


	if os.Args[1] != generators.EULER_TEST && os.Args[1] != generators.TEST_VOLUMEN{
		fmt.Printf("'%s' Tipo de test no implementado\n",os.Args[1])
		return
	}

	if os.Args[1] == generators.EULER_TEST && len(os.Args) < 5{
		fmt.Printf("Insuficientes argumentos se requieren minimo 4 para euler, <Tipo de test a generar>,<directory tests>,<cantidad de tests>, <desde>")		
	}


	var genFromVertices func(string,grafos.Grafo[string,int])

	cantidad_iter := 1
	cantidad_vertices := 200
	cantEdges:= 500 // a futuro capaz cambiar, cantidad de aristas de referencia.
	dirigido := false

	// check cantidad tests
	var err error
	cantidad_iter,err = strconv.Atoi(os.Args[3])
	if err != nil{
		fmt.Printf("Error : %s",err.Error())
		return
	}

	labels := []string{"Makise","Inaba","Misaka","Miu","Asuka","Etc"} // a reemplazar, capaz desde archivo?


	if os.Args[1] == generators.TEST_VOLUMEN{
		// optional args


		if len(os.Args) >= 5{
			cantidad_vertices,err = strconv.Atoi(os.Args[4])

			if err != nil{
				fmt.Printf("Error : %s",err.Error())
				return
			}
		}

		if len(os.Args) >= 6{
			cantEdges,err = strconv.Atoi(os.Args[5])

			if err != nil{
				fmt.Printf("Error : %s",err.Error())
				return
			}
		}

		if len(os.Args) >= 7 && os.Args[6] == "true"{
			dirigido = true
		}

		genFromVertices = func (base string, g grafos.Grafo[string,int]) {
			generarTestVolumen(base,g,cantEdges)
		} 
	} else{
		// optional args

		if len(os.Args) >= 6{
			cantidad_vertices,err = strconv.Atoi(os.Args[7])

			if err != nil{
				fmt.Printf("Error : %s",err.Error())
				return
			}
		}

		if len(os.Args) >= 7{
			cantEdges,err = strconv.Atoi(os.Args[6])

			if err != nil{
				fmt.Printf("Error : %s",err.Error())
				return
			}
		}



		if len(os.Args) >= 8 && os.Args[7] == "true"{
			dirigido = true
		}


		genFromVertices = func (base string, g grafos.Grafo[string,int]) {
			generarTestEuler(base,g,cantEdges,os.Args[4])
		} 

	}

	generateTests(os.Args[2],cantidad_iter, generatorDeVertices(labels,dirigido,cantidad_vertices,genFromVertices))
}