package sesion_votar
import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	TDALista "lista"
)


func ReadAll(archivo *os.File) string{
	res := ""
	var err error = nil
	count:= 1
	buffer := make([]byte,256)
	for err == nil && count >0{
		res += string(buffer[:count])
		count,err = archivo.Read(buffer)
	}

	return res
}

func IterarScanners(scanner *bufio.Scanner,scanner2 *bufio.Scanner,haceAlgo func([]byte,[]byte) bool) error{
	

	for scanner.Scan() && scanner2.Scan() && haceAlgo(scanner.Bytes(),scanner2.Bytes()){
	}

	err:= scanner.Err()
	if(err == nil){
		err = scanner2.Err()
	}

	return err


}



func LeerArchivos(url string,url2 string,haceAlgo func([]byte,[]byte) bool) error{
	archivo,err := os.Open(url)
	if(err != nil){
		return err
	}
	defer archivo.Close()

	archivo2,error2 := os.Open(url2)
	if(error2 != nil){
		return error2
	}

	defer archivo2.Close()

	return IterarScanners(bufio.NewScanner(archivo),bufio.NewScanner(archivo2),haceAlgo)
}



func LeerArchivo(url string,haceAlgo func([]byte) bool) error{
	archivo,error := os.Open(url)
	if(error != nil){
		return error
	}

	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)

	for scanner.Scan() && haceAlgo(scanner.Bytes()){
	}

	return scanner.Err()
}

func CrearArregloDeArchivo[T any](url string , insert func(TDALista.Lista[T],[]byte) error) ([]T,error){
	listaAux:= TDALista.CrearListaEnlazada[T]()

	var err error
	errArchivo := LeerArchivo(url,func (datos []byte) bool{
		err = insert(listaAux,datos)
		return err == nil
	})

	if(errArchivo != nil){
		return make([]T,0),errArchivo
	} else if(err != nil){
		return make([]T,0),err
	}



	res:= make([]T,listaAux.Largo())
	i:= 0
	listaAux.Iterar(func (elemento T) bool{
		res[i] = elemento
		i++
		return true
	})

	return res,nil
}


func particion(arr []Votante, inicio int, final int) ([]Votante, int) {
	pivot := arr[final]
	i := inicio
	for j := inicio; j < final; j++ {
		if arr[j].DNI() < pivot.DNI() {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[final] = arr[final], arr[i]
	return arr, i
}

func quickSort(arr []Votante, low, high int) []Votante {
	if low < high {
		var p int
		arr, p = particion(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func ordenar(lista_dni []Votante,dni int) int { 
	inicio := 0
	fin := len(lista_dni)
	medio := (fin+inicio)/2


	for medio != inicio && lista_dni[inicio].DNI()<dni && lista_dni[fin-1].DNI()>dni{
		if(lista_dni[medio].DNI() == dni){
			fmt.Printf("\n-> dni invalido, repetido")
			return -1 // No se permite repetidos
		}

		if(lista_dni[medio].DNI() > dni){
			fin = medio
		} else{
			inicio = medio+1
		}

		medio = (fin+inicio)/2
	} 
	
	if(lista_dni[inicio].DNI()> dni){
		return inicio
	} 

	if(lista_dni[fin-1].DNI() < dni){
		return fin
	} 

	fmt.Printf("\n-> dni invalido?? min %d, max %d, elem %d",lista_dni[inicio].DNI() ,lista_dni[fin-1].DNI(),dni)
	return -1
}


func redimensionarSlice(viejo []Votante,nuevo_largo int )[]Votante{
	nuevo := make([]Votante,nuevo_largo)

	copy(nuevo,viejo)

	return nuevo
}

//go run main.go ./archivos/catedra/03_partidos ./archivos/catedra/03_padron ./archivos/catedra/03_in ./archivos/catedra/03_out



// devuelve un arreglo de ints, para testeo
func popularVotantesBasico(opciones int) []Votante{
	votantes := make([]Votante,10)
	for i:= 1 ;i<11;i++{
		votantes[i-1] = CrearVotante(i,opciones)
	}
	return votantes
}


// devuelve un arreglo de candidatos, para testeo
func popularCandidatosBasico(tipos int) [][]candidatoStruct{
	candidatos := make([][]candidatoStruct,tipos)

	for i:= range candidatos{
		candidatos[i] = []candidatoStruct{
			candidatoStruct{},
			CrearCandidato("1",fmt.Sprintf("tip %d: %d",i,i)),
			CrearCandidato("2",fmt.Sprintf("tip %d: %d",i,i+1)),
			CrearCandidato("3",fmt.Sprintf("tip %d: %d",i,i+2)),
		}
	}


	return candidatos
}


// devuelve un arreglo de Votantes dado archivo de los dnis
// Estaran ordenados por dni de menor a mayor
func popularVotantes(archivo string,opciones int) ([]Votante,error){
	auxSlice:= make([]Votante,128)
	i:= 0
	errArchivo := LeerArchivo(archivo,func (datos []byte) bool{
			dni,err := strconv.Atoi(string(datos))
			if(err != nil){
				return true
			}

			if(i == len(auxSlice)){
				auxSlice= redimensionarSlice(auxSlice,2*len(auxSlice))
			}
			auxSlice[i] = CrearVotante(dni,opciones)
			i++
			// mientras vaya despues, ie el actual sea menor, itera

			//fmt.Printf("\n-> nuevo dni .... len = %d, cap = %d",i,len(auxSlice))
		return err == nil
	})

	if(errArchivo != nil){
		return make([]Votante,0),errArchivo
	}

	if(i != len(auxSlice)){
		auxSlice= redimensionarSlice(auxSlice,i)
	}



	return quickSort(auxSlice,0,i-1),nil
}

func popularVotantes2(archivo string,opciones int) ([]Votante,error){

	return CrearArregloDeArchivo(archivo,func (lista TDALista.Lista[Votante],bytes []byte) error{
				iterador := lista.Iterador()
				dni,error := strconv.Atoi(string(bytes))
				if(error != nil){
					return nil
				}
				votante := CrearVotante(dni,opciones)
				// mientras vaya despues, ie el actual sea menor, itera
				for iterador.HaySiguiente() && iterador.VerActual().DNI() < votante.DNI(){ 
					iterador.Siguiente()
				}

				iterador.Insertar(votante)


				fmt.Printf("\n-> nuevo dni .... len = %d",lista.Largo())
				return nil
			})
	//return popularVotantesBasico(opciones)
}




// devuelve una lista dada la cantidad de tipos de votos y un archivo de los candidatos
func popularCandidatos(archivo string,tipos int) ([][]candidatoStruct,error){
	candidatos := make([][]candidatoStruct,tipos)
	cant:= 4
	errArchivo := LeerArchivo(archivo,func (datos []byte) bool{
		cant = len(strings.Split(string(datos),","))//-1
		return false
	})

	for i:= range candidatos{
		candidatos[i] = make([]candidatoStruct,cant)
		candidatos[i][0] = candidatoStruct{}
		for j:=1;j<len(candidatos[i]);j++{
			candidatos[i][j] = CrearCandidato(fmt.Sprintf("partido %d",j),fmt.Sprintf("tip %d: %d",i,j))
		}
	}


	return candidatos,errArchivo
	
	/*
	return CrearArregloDeArchivo(archivo,func (lista TDALista.Lista[[]candidatoStruct],bytes []byte) error{
				lista.InsertarUltimo(
					[]candidatoStruct {candidatoStruct{},
					CrearCandidato("1","Pre 1"),
					CrearCandidato("2","Pre 2"),
					CrearCandidato("3","Pre 3")})

				//fmt.Printf("\nSE AGREGO? %d",lista.Largo())
				if(len(lista.VerUltimo()) != tipos){
					return new(ErrorLecturaArchivos)
				}
				return nil
			})

	*/
}



func AccionDesdeComando(sesion SesionVotar,comando string) error{	
	if(comando == "fin-votar"){
		return sesion.SiguienteVotante()
	} else if(comando == "deshacer"){
		return sesion.Deshacer()
	}

	args := strings.Split(comando," ")

	if(args[0] == "ingresar"){

		if(len(args) < 2){ // yo pondria != pero el error es solo en falta
			return new(ErrorFaltanParametros)
		}

		return sesion.IngresarVotante(args[1])

	}

	if(args[0] == "votar"){

		if(len(args) < 3){
			return new(ErrorFaltanParametros)
		}

		return sesion.Votar(args[1],args[2])
	}
	
	return new(ErrorComandoInvalido)

}

func AccionComandoAString(sesion SesionVotar,comando string) string{
	err:= AccionDesdeComando(sesion,comando)
	res:= OK
	if(err != nil){
		res = err.Error()
	}
	return res
}


// Se podria poner como primitiva pero se decidio mas elegante el no hacerlo..
func MostrarEstado(sesion SesionVotar,identificadores []string){
	for _,identificador :=range identificadores{

		fmt.Printf("\n\n%s",identificador)

		sesion.IterarVotos(identificador, func (credencial string,votos int) bool{
			fmt.Printf("\n %s: %d\n",credencial,votos)	
			return true
		})

	}

	fmt.Printf("\n\n Votos Impugnados: %d",sesion.VotosImpugnados())

}






// Tests Utilities

func TestearComando(sesion SesionVotar,comando string, expected string) error{
	res:= AccionComandoAString(sesion,comando)
	if res != expected{
		return CrearErrorTest(comando,expected,res)
	}

	return nil
}
func TestearComandosScanners(sesion SesionVotar,input_scanner *bufio.Scanner,output_scanner *bufio.Scanner) error{
	var err error = nil
	IterarScanners(input_scanner,output_scanner,func(comando []byte,expected []byte ) bool{
		err = TestearComando(sesion,string(comando),string(expected))
		return err == nil
	})	

	return err
}

func TestearComandosArchivos(sesion SesionVotar,input_file string, out_put_file string) error{
	in,err := os.Open(input_file)
	if(err != nil){
		return err
	}
	defer in.Close()

	expectedOut,error2 := os.Open(out_put_file)
	if(error2 != nil){
		return error2
	}

	defer expectedOut.Close()


	return TestearComandosScanners(sesion,bufio.NewScanner(in),bufio.NewScanner(expectedOut))
}