package sesion_votar
import (
	"fmt"
	"os"
	"bufio"
	"strings"
	TDALista "lista"
)


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
func popularVotantesBasico(opciones int) []Votante{
	votantes := make([]Votante,9)
	for i:= range votantes{
		votantes[i] = crearVotante(i,opciones)
	}
	return votantes
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
func popularVotantes(archivo string,opciones int) []Votante{
	return popularVotantesBasico(opciones)
}




// devuelve una lista dada la cantidad de tipos de votos y un archivo de los candidatos
func popularCandidatos(archivo string) [][]candidatoStruct{
	return popularCandidatosBasico()
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





