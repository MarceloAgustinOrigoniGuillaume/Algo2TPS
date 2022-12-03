package generators

import "tp3/grafos"
import grafosLib "tp3/grafos/lib"
import "tp3/utils/PJ"
import "tp3/utils"
import "strconv"
import "strings"
import "errors"
import "fmt"
import "time"


const EULER_TEST = "euler"
const TOPOLOGICO_TEST = "Topologico"
const DIJKSTRA_TEST = "Dijkstra"
const KRUSKAL_TEST = "Kruskal"
const READ_PJ = "PJ"
const READ_RECOMENDACIONES = "CSV_R"
const NIL_MARK = "*nil*"
const MS = 1000 //000
//const ITINERARIO_TEST = "Itinerario"

// Writing
func WritePj(file string,grafo grafos.Grafo[string,int]) error{

	builder,err := pj.CrearPJ(file)
	if err != nil{
		return err
	}

	defer builder.ClosePJ()

	builder.StartPJ(grafo.CantidadVertices(),grafo.CantidadAristas())
	coord_x := 0
	coord_y := 0

	grafo.IterarVertices(func (nombre string) bool{

		coord_y++
		if coord_y > 100{
			coord_x++
			coord_y= 0
		}

		builder.AddCity(nombre,strconv.Itoa(coord_x),strconv.Itoa(coord_y))
		

		return true
	})

	grafo.IterarAristas(func (desde string,hasta string, peso int) bool{
		builder.AddArista(desde,hasta,peso)
		return true
	})

	return nil
}





func WriteEuler(outFile string,pj_file string,titulo string,dirigido bool,
	desde *string,hasta *string,pesoTotal int,tipo int) error{
	archivo, err := utils.AbrirOCrearArchivo(outFile)

	if err != nil{
		return err
	}
	args:= EULER_TEST+", "+titulo+", "+pj_file+", "

	if dirigido{
		args+= "true"
	} else{
		args+= "false"
	}


	args += "\n"+strconv.Itoa(tipo)+", "

	if tipo == TIPO_CAMINO_EULERIANO{
		args+= *desde
		args+= ", "+(*hasta) // ambos deberian ser != nil
	//} else if tipo == TIPO_CICLO_EULERIANO{
		//args = NIL_MARK
		//args += ", "+args // copia desde
	} else {
		args += NIL_MARK+", "+NIL_MARK
	}

	args += ", "+strconv.Itoa(pesoTotal)
	archivo.WriteString(args)

	return nil
}


func BasicRead(line string) (string,error){
	splitted := strings.SplitN(line,",",3)

	if len(splitted) != 3{
		return "",errors.New(fmt.Sprintf("Tenia cantidad incorrecta de datos, deberian ser 3... %v",splitted))
	}


	return splitted[0],nil
}


type structTestEuler struct{
	titulo,pj_file string
	dirigido bool

	desde,hasta string
	pesoTotal int
	tipo int
}

func readEuler(test_file string) (structTestEuler,error){
	linea := 0

	test_struct := structTestEuler{}
	var errorL error
	err := utils.LeerArchivo(test_file, func (line string) bool{
		if linea == 0{
			splitted := strings.SplitN(line,", ",4)

			if len(splitted) <4{
				errorL = errors.New(fmt.Sprintf("Linea de credenciales tenia mal argumentos deberia ser <test>, <titulo>, <pj_file>, <es dirigido> fue : '%s'",line))
			}

			test_struct.titulo = splitted[1]
			test_struct.pj_file = splitted[2]

			if splitted[3] == "true"{
				test_struct.dirigido = true
			}


		} else if linea == 1{
			splitted := strings.SplitN(line,", ",4)

			if len(splitted) <4{
				errorL = errors.New(fmt.Sprintf("Linea de argumentos tenia mal argumentos, deberia ser <tipo>, <desde>, <hasta>, <peso total> fue : '%s'",line))
			}

			test_struct.tipo,errorL = strconv.Atoi(splitted[0])
			if errorL != nil{
				return false
			}
			
			test_struct.pesoTotal,errorL = strconv.Atoi(splitted[3])

			if errorL != nil{
				return false
			}

			test_struct.desde = splitted[1]
			test_struct.hasta = splitted[2]

		}
		linea ++
		return errorL != nil || linea <2
	})
	
	if errorL == nil{
		errorL = err
	} 

	return test_struct,errorL
}


func runCicloEulerTest(grafo grafos.Grafo[string,int], origin string, test_info structTestEuler) (int64,error){
	fmt.Printf("-->Ciclo euleriano:\n")
	timeTest := time.Now().UnixNano()
	timeEnd := timeTest
	camino,errC := grafosLib.CicloEuleriano(grafo,origin)

	if errC != nil{
		if test_info.tipo == TIPO_CICLO_EULERIANO{
			return (time.Now().UnixNano()-timeTest),errC 
		}
		timeEnd = time.Now().UnixNano()
		fmt.Printf("No encontro ciclo euleriano y no debia\n")
	} else{


		if test_info.hasta != NIL_MARK{
			if test_info.desde != camino[0] {
				return (time.Now().UnixNano()-timeTest),errors.New(fmt.Sprintf("Desde no era correcto, should %s, have %s",test_info.desde,camino[0]))
			}

			if test_info.hasta != camino[len(camino)-1] {
				return (time.Now().UnixNano()-timeTest),errors.New(fmt.Sprintf("Hasta no era correcto, should %s, have %s",test_info.hasta,camino[len(camino)-1]))
			}

		}

		tiempoTotal := 0
		anterior := camino[0]
		for _,elem := range camino[1:]{
			tiempoTotal += grafo.ObtenerPeso(anterior,elem)
			anterior = elem
		}

		timeEnd = time.Now().UnixNano()


		if tiempoTotal != test_info.pesoTotal {
			return timeEnd- timeTest,errors.New(fmt.Sprintf("La suma total de pesos era distinta, should %d, have %d",test_info.pesoTotal,tiempoTotal))
		}

		//fmt.Printf("Test ciclo euleriano: PASS took %d ns\n",timeEnd - timeTest)
	}

	return (timeEnd-timeTest),nil 
}




func runCaminoEulerTest(grafo grafos.Grafo[string,int], origin string, test_info structTestEuler) (int64,error){
	fmt.Printf("-->Camino euleriano:\n")

	timeTest := time.Now().UnixNano()
	timeEnd := timeTest
	camino,errC := grafosLib.CaminoEuleriano(grafo,origin)

	if errC != nil{
		if test_info.tipo != TIPO_NADA{
			return (time.Now().UnixNano()-timeTest),errC
		}
		timeEnd = time.Now().UnixNano()
		
		fmt.Printf("No encontro camino euleriano y no debia\n")
	} else{

		if test_info.hasta != NIL_MARK{
			if test_info.desde != camino[0] {
				return (time.Now().UnixNano()-timeTest),errors.New(fmt.Sprintf("Desde no era correcto, should %s, have %s",test_info.desde,camino[0]))
			}

			if test_info.hasta != camino[len(camino)-1] {
				return (time.Now().UnixNano()-timeTest),errors.New(fmt.Sprintf("Hasta no era correcto, should %s, have %s",test_info.hasta,camino[len(camino)-1]))
			}
		}

		tiempoTotal := 0
		anterior := camino[0]
		for _,elem := range camino[1:]{
			tiempoTotal += grafo.ObtenerPeso(anterior,elem)
			anterior = elem
		}
		timeEnd = time.Now().UnixNano()

		if tiempoTotal != test_info.pesoTotal {
			return (time.Now().UnixNano()-timeTest),errors.New(fmt.Sprintf("La suma total de pesos era distinta, should %d, have %d",test_info.pesoTotal,tiempoTotal))
		}

	}

	return (timeEnd-timeTest),nil 
}

func ExecuteEulerTest(test_file string) (grafos.Grafo[string,int],error){
	timeInit := time.Now().UnixNano()

	test_info,errorL := readEuler(test_file)

	if errorL != nil{
		return nil,errorL
	} 

	grafo := grafos.GrafoNumericoPesado[string,int](test_info.dirigido)
	errorL = pj.LeerPJ(test_info.pj_file,grafo,BasicRead,strconv.Atoi)

	if errorL != nil{
		return nil,errorL
	} 

	fmt.Printf("test: %s pj: %s build took: %.2f pico\n",test_info.titulo,test_info.pj_file,(float64)(time.Now().UnixNano()-timeInit)/MS)

	desdeValor := test_info.desde

	if desdeValor == NIL_MARK{
		grafo.IterarVertices(func (vert string) bool{ // para que no salte que no existe
			desdeValor = vert
			return false
		})
	}

	timeTaken, err := runCicloEulerTest(grafo,desdeValor,test_info)
	if err != nil{
		fmt.Printf("Error at ciclo euler: %s, took %.2f pico\n",err.Error(),(float64)(timeTaken)/MS)
	} else{		
		fmt.Printf("Ciclo euler: PASS, took %.2f pico \n",(float64)(timeTaken)/MS)
	}
	timeTaken, errorL = runCaminoEulerTest(grafo,desdeValor,test_info)

	if errorL != nil{
		fmt.Printf("Error at camino euler: %s, took %.2f pico \n",errorL.Error(),(float64)(timeTaken)/MS)
	} else{		
		fmt.Printf("Camino euler: PASS, took %.2f pico \n",(float64)(timeTaken)/MS)
	}


	fmt.Printf("===================================\n")
	errFinal := err
	if errorL != nil{
		if errFinal != nil{
			errFinal = errors.New(fmt.Sprintf("Ciclo: %s\n Camino: %s \n",err.Error(),errorL.Error()))
		}


		showGrados(grafo)
		cant := 0
		first := -1
		summed := 0
		tipo_act:= TIPO_CAMINO_EULERIANO
		grafosLib.DFS_ALL(grafo,func (_ string){cant++}, func(){
			fmt.Printf("Comp with %d\n",cant)
			summed += cant
			if first == -1{
				first = cant
				cant = 0
				fmt.Printf("---------->Una comp no tenia todos los vertices\n")
				return
			}

			if cant != 1{
				tipo_act = TIPO_NADA
				fmt.Printf("---------->No era conexo ... found other component %d\n",cant)
			}
			cant = 0

		})


		fmt.Printf("dfs all suma %d , cantidad vertices %d \n",summed, grafo.CantidadVertices())
		if tipo_act == TIPO_NADA{
			fmt.Printf("-------------->Cambio a tipo nada.. dfs all suma %d , cantidad vertices %d \n",summed, grafo.CantidadVertices())
		}


		errFinal = errorL
	}

	return grafo,errFinal
}








