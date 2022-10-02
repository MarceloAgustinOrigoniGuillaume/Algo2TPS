package sesion_votar
import (
	"fmt"
	"os"
	"bufio"
	TDALista "lista"
)

type ErrorOmicion struct{
}

func (err *ErrorOmicion) Error() string{
	return "ERROR: elemento invalido, pero se puede ignorar"
}

type ErrorMissMatchSizeOut struct{
}

func (err *ErrorMissMatchSizeOut) Error() string{
	return "ERROR: Habia mas lineas en el archivo out que en el in. Se ignoraron las sobrantes"
}


type ErrorMissMatchSizeIn struct{
}

func (err *ErrorMissMatchSizeIn) Error() string{
	return "ERROR: Habia mas lineas en el archivo in que en el out. Se ignoraron las sobrantes"
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


func CrearArregloDeArchivo[T any](url string , convertidor func([]byte) (T,error)) ([]T,error){
	listaAux:= TDALista.CrearListaEnlazada[T]()

	var elem T
	var err error
	LeerArchivo(url,func (datos []byte) bool{
		elem,err = convertidor(datos)
		


		if err != nil{
			
			switch err.(type){
				case *ErrorOmicion:
					// pass
					err = nil
				default:
					return false


			}

		}
		listaAux.InsertarUltimo(elem)
		return true
	})

	if(err != nil){
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


// devuelve un arreglo de ints, para testeo
func popularDNISBasico() []int{
	return []int{1,2,3,4,5,6,7,8,9,10}
}


// devuelve un arreglo de candidatos, para testeo
func popularCandidatosBasico() [][]candidatoStruct{
	return [][]candidatoStruct{
		{candidatoStruct{},CrearCandidato("1","Pre 1"),CrearCandidato("2","Pre 2"),CrearCandidato("3","Pre 3")},
		{candidatoStruct{},CrearCandidato("1","Gob 1"),CrearCandidato("2","Gob 2"),CrearCandidato("3","Gob 3")},
		{candidatoStruct{},CrearCandidato("1","Dip 1"),CrearCandidato("2","Dip 2"),CrearCandidato("3","Dip 3")},
	}
}


// devuelve un arreglo int dado archivo de los dnis
func popularDNIS(archivo string) []int{
	return popularDNISBasico()
}




// devuelve una lista dada la cantidad de tipos de votos y un archivo de los candidatos
func popularCandidatos(archivo string) [][]candidatoStruct{
	return popularCandidatosBasico()
}

// Busqueda binaria deberia hacerse como se observa devuelve false si no esta y true si lo esta 
func Contiene(lista_dni []int,dni int) bool { 
	for _,dni_comp := range lista_dni{
		if(dni_comp == dni){
			return true
		}
	}

	return false
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