package agencia

import kml "tp3/utils/KML" //utilitis
import pj "tp3/utils/PJ"
import utils "tp3/utils" //utilities
import "tp3/grafos"
import libGrafos "tp3/grafos/lib"
import hash "tp3/diccionario"

//.Arista

import "fmt"
import "errors"
import "strconv"
import "strings"



func showGrados[T any](grafo grafos.Grafo[string,T]){
	entradas := libGrafos.GradosDeEntrada(grafo)

	libGrafos.GradosDeSalida(grafo).Iterar(func (ciudad string,grados int) bool{
		fmt.Printf("De %s se va a %d lugares y se llega por %d\n",ciudad,grados,entradas.Obtener(ciudad))
		return true
	})

}

const MSJ_REDUCIR = "Peso total: %d"
const NAME_KML_LINE = "from: %s to: %s"
const NAME_CAMINO = "Peso total: %f"



type AgenciaViajes interface{
	Ir(desde,hasta, outFile string) (string ,error)
	ViajeDesde(desde, outFile string) (string ,error)
	Itinerario(recomendacionesFile string) (string ,error)
	ReducirCaminos(outFile string) (string ,error)
}



type ciudadStruct struct{
	nombre string
	latitud string
	longitud string
}

func crearCiudad(nombre, latitud, longitud string) ciudadStruct{
	return ciudadStruct{nombre,latitud,longitud}
}

type agenciaStruct struct{
	indexadoCiudades hash.Diccionario[string,ciudadStruct]
	grafoLugar grafos.Grafo[string,int]
}

func (agencia *agenciaStruct) itinerarioDesdeRecomendaciones(archivoRecomendaciones string) ([]string,error){
	precedencias := grafos.GrafoNumericoNoPesado[string](true)

	agencia.grafoLugar.IterarVertices(func (ciudad string) bool{
		precedencias.AgregarVertice(ciudad)
		return true
	})

	var errInterno error
	err := utils.LeerArchivo(archivoRecomendaciones, func(line string) bool{
		splitted := strings.SplitN(line,",",2)

		if len(splitted) != 2{
			errInterno = errors.New(fmt.Sprintf("Conexion de ciudades tiene cantidad incorrecta de datos, deberian ser 2... %v",splitted))
			return false
		}

		precedencias.AgregarArista(splitted[0],splitted[1],1)	
		return true
	})

	if errInterno == nil{
		errInterno = err
	}


	if errInterno != nil{
		return nil,errInterno
	}

	return libGrafos.SecuenciaTopologica(precedencias)



}

func (agencia *agenciaStruct) readCiudad(line string) (string,error){
	splitted := strings.SplitN(line,",",3)

	if len(splitted) != 3{
		return "",errors.New(fmt.Sprintf("Ciudad con cantidad incorrecta de datos, deberian ser 3... %v",splitted))
	}


	agencia.indexadoCiudades.Guardar(splitted[0],crearCiudad(splitted[0],splitted[1],splitted[2]))

	return splitted[0],nil
}

func CrearAgenciaViajes(archivoLugar string) (AgenciaViajes,error){
	agencia:= new(agenciaStruct)

	agencia.grafoLugar = grafos.GrafoNumericoPesado[string,int](false)
	agencia.indexadoCiudades = hash.CrearHash[string,ciudadStruct]()

	return agencia , pj.LeerPJ(archivoLugar,agencia.grafoLugar,agencia.readCiudad,strconv.Atoi)
}


func (agencia *agenciaStruct) iterarCiudadesCamino(ciudades []string, visitar func(bfr,act ciudadStruct) bool){
	bfr_ciudad := agencia.indexadoCiudades.Obtener(ciudades[0])
		
	var act_ciudad ciudadStruct

	for _,nombre_ciudad:= range ciudades[1:]{
		act_ciudad = agencia.indexadoCiudades.Obtener(nombre_ciudad)

		if !visitar(bfr_ciudad,act_ciudad){
			return
		}

		bfr_ciudad = act_ciudad
	}
}


func (agencia *agenciaStruct) Ir(desde, hasta, outFile string) (string,error){


	if !agencia.indexadoCiudades.Pertenece(desde) || !agencia.indexadoCiudades.Pertenece(hasta) {
		return "",libGrafos.ErrorRecorrido()
	}

	caminoEsp,errC := libGrafos.CaminoMinimoDijkstraHasta(agencia.grafoLugar,desde,hasta)

	if errC != nil{
		return "",libGrafos.ErrorRecorrido()
	}

	anterior := caminoEsp[0]
	res:= fmt.Sprintf("Camino desde %s hasta %s",desde,hasta)
	for _,ciudad := range caminoEsp[1:]{
		res += fmt.Sprintf("\nIr de %s a %s",anterior,ciudad)
		anterior = ciudad
	}

	builder,err := kml.CrearKML(outFile)
	if err == nil{
		builder.StartKML(fmt.Sprintf("Camino desde %s hasta %s",desde,hasta))
		origin:= agencia.indexadoCiudades.Obtener(desde)

		builder.AddPoint(origin.nombre,origin.latitud,origin.longitud)
		agencia.iterarCiudadesCamino(caminoEsp, func (bfr_ciudad,act_ciudad ciudadStruct) bool{
			builder.AddPoint(act_ciudad.nombre,act_ciudad.latitud,act_ciudad.longitud)
			builder.AddLine( // add a line
				fmt.Sprintf(NAME_KML_LINE,bfr_ciudad.nombre,act_ciudad.nombre),// title
				bfr_ciudad.latitud,bfr_ciudad.longitud, // start coords
				act_ciudad.latitud,act_ciudad.longitud) // end coords
			return true
		})
		builder.CloseKML()
	}


	return res,err
}


func (agencia *agenciaStruct) Itinerario(recomendacionesFile string) (string,error){
	itinerario,err := agencia.itinerarioDesdeRecomendaciones(recomendacionesFile)

	if err != nil{
		return "",libGrafos.ErrorRecorrido()
	}


	res := itinerario[0]
	for _,elem := range itinerario[1:]{

		res+= " -> "+elem
	}

	return res,nil
}

func (agencia *agenciaStruct) writePj(file string,aristas []libGrafos.Arista[string,int]) (int,error){

	pesoTotal := 0
	builder,err := pj.CrearPJ(file)
	if err != nil{
		return 0,err
	}

	defer builder.ClosePJ()

	builder.StartPJ(agencia.indexadoCiudades.Cantidad(),len(aristas))
	agencia.indexadoCiudades.Iterar(func (nombre string,ciudad ciudadStruct) bool{
		builder.AddCity(nombre,ciudad.latitud,ciudad.longitud)
		return true
	})

	for _,arista := range aristas{
		builder.AddArista(arista.Desde(),arista.Hasta(),arista.Peso())
		pesoTotal+= arista.Peso()
	}


	return pesoTotal,nil
}


func (agencia *agenciaStruct) ReducirCaminos(destinoFile string) (string,error){

	aristasMST:= libGrafos.MSTKruskal(agencia.grafoLugar,libGrafos.QuickSortAristas[string,int])
	pesoTotal,err:= agencia.writePj(destinoFile,aristasMST)

	return fmt.Sprintf(MSJ_REDUCIR,pesoTotal),err
}




func (agencia *agenciaStruct) ViajeDesde(desde string, outFile string) (string,error){


	if !agencia.indexadoCiudades.Pertenece(desde){
		return "", libGrafos.ErrorRecorrido()
	}



	camino,errC := libGrafos.CicloEuleriano(agencia.grafoLugar,desde)

	if errC != nil{
		return "",libGrafos.ErrorRecorrido()
	}
	res := camino[0]
	tiempoTotal := 0
	anterior := camino[0]
	for _,elem := range camino[1:]{

		res+= " -> "+elem

		tiempoTotal += agencia.grafoLugar.ObtenerPeso(anterior,elem)
		anterior = elem
	}

	res += fmt.Sprintf("\nTiempo total:%d",tiempoTotal)


	builder,err := kml.CrearKML(outFile)

	if err == nil{
		builder.StartKML(fmt.Sprintf("Camino desde y hasta %s pasando por todas las rutas",desde))

		agencia.indexadoCiudades.Iterar(func (nombre string,ciudad ciudadStruct) bool{
			builder.AddPoint(nombre,ciudad.latitud,ciudad.longitud)
			return true
		})

		agencia.iterarCiudadesCamino(camino, func (bfr_ciudad,act_ciudad ciudadStruct) bool{
			builder.AddLine( // add a line
				fmt.Sprintf(NAME_KML_LINE,bfr_ciudad.nombre,act_ciudad.nombre),// title
				bfr_ciudad.latitud,bfr_ciudad.longitud, // start coords
				act_ciudad.latitud,act_ciudad.longitud) // end coords
			return true
		})
		builder.CloseKML()	
	}

	return res,err
}