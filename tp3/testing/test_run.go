package main

import (
	"fmt"
	"os"
	"tp3/testing/generators"
	"tp3/grafos"
	"strconv"
	"strings"
	"errors"
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


func basicRead(line string) (string,error){
	splitted := strings.SplitN(line,",",3)

	if len(splitted) != 3{
		return "",errors.New(fmt.Sprintf("Tenia cantidad incorrecta de datos, deberian ser 3... %v",splitted))
	}


	return splitted[0],nil
}

func main(){
	if len(os.Args) < 3{
		fmt.Printf("Insuficientes argumentos se requieren 3, la ruta del pj y el tipo de test y si es dirigido/pesado ")
		return
	}

	dirigido := false
	labels := []string{"Makise","Inaba","Misaka","Miu","Asuka","Etc"}

	if len(os.Args) >= 4 && os.Args[3] == "true"{
		dirigido = true
	}
	cantidad_iter := 10
	var err error
	if len(os.Args) >= 5{
		cantidad_iter,err = strconv.Atoi(os.Args[4])

		if err != nil{
			fmt.Printf("Error : %s",err.Error())
			return
		}
	}


	cantidad_vertices := 200
	if len(os.Args) >= 6{
		cantidad_vertices,err = strconv.Atoi(os.Args[4])

		if err != nil{
			fmt.Printf("Error : %s",err.Error())
			return
		}
	}

	cantEdges:= 500

	desdeValor := os.Args[2]

	grafo := generators.BuildVertices[int](dirigido,labels,cantidad_vertices, generators.HaceNada,generators.HaceNada)

	var copy grafos.Grafo[string,int]
	
	base_url := "tests/grafos/"
	ind := 0
	values := []int{generators.TIPO_CICLO_EULERIANO,generators.TIPO_CAMINO_EULERIANO,generators.TIPO_NADA}
	labels_types := []string{"TIPO_CICLO_EULERIANO","TIPO_CAMINO_EULERIANO","TIPO_NADA"}

	var timeInit int64
	var timeTest int64
	var timeEnd int64

	for ind_value,curr_value := range values{
		fmt.Printf("\n.................Curr value %s \n",labels_types[ind_value])

		for i:= 0;i<cantidad_iter;i++{
			timeInit = time.Now().UnixNano()
			copy = generators.CopyVerts(grafo,dirigido)
			desde,hasta,pesoTotal := generators.BuildGrafoEuler(time.Now().UnixNano(),copy,cantEdges,func (edge int) int {return rand.Intn(edge+1)},curr_value)



			if i == 0{
				fmt.Printf("------------>Saving BASIC TEST! ")
				ind++
				tipo_act := curr_value

				if tipo_act != generators.TIPO_NADA{
					cant := 0
					first := -1
					summed := 0
					grafosLib.DFS_ALL(copy,func (_ string){cant++}, func(){
						summed += cant
						if first == -1{
							first = cant
							cant = 0
							fmt.Printf("---------->Una comp no tenia todos los vertices\n")
							return
						}

						if cant != 1{
							tipo_act = generators.TIPO_NADA
							fmt.Printf("---------->No era conexo ... found other component %d\n",cant)
						}
						cant = 0

					})


					if tipo_act == generators.TIPO_NADA{
						fmt.Printf("-------------->Cambio a tipo nada.. dfs all suma %d , cantidad vertices %d \n",summed, copy.CantidadVertices())
					}
				}

				pj_url := base_url+"tests/"+strconv.Itoa(ind)+"_pj_"+strconv.Itoa(tipo_act)+".pj"
				errC2:= generators.WritePj(pj_url,copy)

				if errC2!= nil{
					fmt.Printf("Err Writing test PJ: %s\n",errC2.Error())
					return
				}
				url_test := base_url+"tests/"+strconv.Itoa(ind)+"_test_"+strconv.Itoa(tipo_act)+".test"
				errC2 = generators.WriteEuler(url_test,pj_url, "test_"+labels_types[ind_value],dirigido,desde,hasta,pesoTotal,tipo_act)

				if errC2!= nil{
					fmt.Printf("Err Writing test : %s\n",errC2.Error())
					return
				}

			}

			fmt.Printf(">Searching camino eulariano from %s, should have %d\n",desdeValor,pesoTotal)

			if desde != nil{
				fmt.Printf("//Builder said %s, %s, %d\n",*desde,*hasta,pesoTotal)
			}
			timeTest = time.Now().UnixNano()
			camino,errC := grafosLib.CicloEuleriano(copy,desdeValor)

			if errC != nil{
				if curr_value == generators.TIPO_CICLO_EULERIANO{
					fmt.Printf("Err : %s\n",errC.Error())
					ind++
					err = generators.WritePj(base_url+strconv.Itoa(ind)+"_t_"+strconv.Itoa(curr_value)+".pj",copy)
					continue 
				}
				fmt.Printf("No encontro ciclo y no debia \n")
			} else{

				tiempoTotal := 0
				anterior := camino[0]
				for _,elem := range camino[1:]{
					tiempoTotal += copy.ObtenerPeso(anterior,elem)
					anterior = elem
				}

				timeEnd = time.Now().UnixNano()

				fmt.Printf("Ciclo euleriano start: %s end: %s pesoTotal: %v, len %d\n",camino[0],camino[len(camino)-1],tiempoTotal,len(camino))
				fmt.Printf("Time build : %d Time test: %v\n\n",timeTest- timeInit, timeEnd - timeTest)				
			}


			if desde != nil{
				timeTest = time.Now().UnixNano()
				camino,errC := grafosLib.CaminoEuleriano(copy,*desde)

				if errC != nil{
					if curr_value == generators.TIPO_NADA{
						fmt.Printf("No encontro camino y no debia \n")
						continue
					}

					fmt.Printf("Err : %s\n",errC.Error())
					ind++
					err = generators.WritePj(base_url+strconv.Itoa(ind)+".pj",copy)
					continue 
				}



				tiempoTotal := 0
				anterior := camino[0]
				for _,elem := range camino[1:]{
					tiempoTotal += copy.ObtenerPeso(anterior,elem)
					anterior = elem
				}

				timeEnd = time.Now().UnixNano()

				fmt.Printf("Camino euleriano start: %s end: %s pesoTotal: %v\n",camino[0],camino[len(camino)-1],tiempoTotal)
				fmt.Printf("Time build : %d Time test: %v\n",timeTest- timeInit, timeEnd - timeTest)

			}
			
			fmt.Printf("===================================\n")

		}

	}
}