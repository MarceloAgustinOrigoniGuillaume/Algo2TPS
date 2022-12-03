package generators

import "tp3/grafos"
import "math/rand"
//import "fmt"

const TIPO_CICLO_EULERIANO = 1
const TIPO_CAMINO_EULERIANO = 2
const TIPO_NADA = 0


func buildGrafoEulerDirigido[V comparable, T grafos.Numero](grafo grafos.Grafo[V,T], edge_quantity int,pesoProvider func(int) T,tipo int) (*V,*V,T){
	//rand.Seed(time.Now().UnixNano())
	var pesoTotal T = 0
	vertices := grafo.ObtenerVertices()

	length := len(vertices)
	
	gruposCorrectos := (int)(edge_quantity/3)
	//aristasUsadas := set.CrearSetWith[int](gruposCorrectos)

	//var choosen int
	selectedDesde := -1
	count_aristas_desde := 0
	shouldKeepOne := tipo != TIPO_NADA
	var selectedHasta int
	var selectedIntermedio int
	var pesoCurr T
	//square := length*length
	//maxInd := square*length

	countGrupos := 0
	for countGrupos < gruposCorrectos{

		if shouldKeepOne{
			count_aristas_desde++
			if !grafo.EsDirigido(){
				count_aristas_desde++
			}

			if count_aristas_desde >= grafo.CantidadVertices()-1{
				shouldKeepOne = false
			}
		}

		if selectedDesde == -1 || !shouldKeepOne{
			selectedDesde = rand.Intn(length)
		}
		selectedHasta = rand.Intn(length)
		selectedIntermedio = rand.Intn(length)

		// max ind = length*length*length, primer indice = ind % length, max => 0 
		// segundo indice = (maxIndice/length) % length max => 0
		// tercer indice = maxIndice / square max => length


		// Verificamos no sean iguales, sin lazos
		if selectedDesde == selectedHasta || selectedDesde == selectedIntermedio || selectedHasta == selectedIntermedio{
			continue
		}


		// Verificamos no sea repetida
		if (grafo.ExisteArista(vertices[selectedDesde],vertices[selectedHasta]) || 
		grafo.ExisteArista(vertices[selectedHasta],vertices[selectedIntermedio])||
		grafo.ExisteArista(vertices[selectedIntermedio],vertices[selectedDesde])){
			continue
		}

		pesoCurr = pesoProvider(grafo.CantidadAristas())
		grafo.AgregarArista(vertices[selectedDesde],vertices[selectedHasta],pesoCurr)
		pesoTotal+=pesoCurr

		pesoCurr = pesoProvider(grafo.CantidadAristas())
		grafo.AgregarArista(vertices[selectedHasta],vertices[selectedIntermedio],pesoCurr)
		pesoTotal+=pesoCurr

		pesoCurr = pesoProvider(grafo.CantidadAristas())
		grafo.AgregarArista(vertices[selectedIntermedio],vertices[selectedDesde],pesoCurr)
		pesoTotal+=pesoCurr

		countGrupos++
	}



	if tipo == TIPO_CAMINO_EULERIANO{
		var hasta V 
		var desde *V = nil

		for desde == nil{
			ind:= rand.Intn(length)
			hasta = vertices[ind] 
			//fmt.Printf("EVAl hasta .. %v, ind := %d\n",hasta,ind)

			grafo.IterarAdyacentes(hasta,func (ady V,peso T) bool{
				//fmt.Printf("Found desde .. %v\n",ady)
				desde = &ady
				grafo.BorrarArista(hasta,ady)
				pesoTotal-= peso
				return false
			})			
		}


		return desde,&hasta,pesoTotal
	} 

	if tipo == TIPO_NADA && grafo.CantidadAristas() >= length-1{
		var hasta V 
		removidas := 0

		for removidas < 2{
			hasta = vertices[rand.Intn(length)] 

			grafo.IterarAdyacentes(hasta,func (ady V,peso T) bool{
				grafo.BorrarArista(hasta,ady)
				pesoTotal-= peso
				removidas++
				return removidas <2
			})			
		}



	}

	return nil,nil,pesoTotal

}

func BuildGrafoEuler[V comparable, T grafos.Numero](seed int64, grafo grafos.Grafo[V,T], edge_quantity int,pesoProvider func(int) T,tipo int) (*V,*V,T){
	rand.Seed(seed)
	desde,hasta,total := buildGrafoEulerDirigido(grafo,edge_quantity,pesoProvider,tipo)

	// Deberia chequearse que de haber quererse un camino euleriano, debe ser conexo, o almenos no disconjunto
	return desde,hasta,total
}
