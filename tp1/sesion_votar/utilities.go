package sesion_votar

import (
	"bufio"
	"fmt"
	TDALista "lista"
	"os"
	"strings"
)

// devuelve un arreglo de ints, para testeo
func popularVotantesBasico(opciones int) []Votante {
	votantes := make([]Votante, 10)
	for i := 1; i < 11; i++ {
		votantes[i-1] = CrearVotante(i, opciones)
	}
	return votantes
}

// devuelve un arreglo de candidatos, para testeo
func popularCandidatosBasico(tipos int) [][]candidatoStruct {
	candidatos := make([][]candidatoStruct, 4)
	personas := []string{"Blanco", "Persona B", "Amigo A", "Otro"}
	for i := range candidatos {

		candidatos[i] = make([]candidatoStruct, tipos)
		partido := fmt.Sprintf("Partido %d", i)
		for j := range candidatos[i] {

			candidatos[i][j] = CrearCandidato(partido, fmt.Sprintf("%s%d", personas[i], j+1))
		}
	}

	return candidatos
}

func ReadAll(archivo *os.File) string {
	res := ""
	var err error = nil
	count := 1
	buffer := make([]byte, 256)
	for err == nil && count > 0 {
		res += string(buffer[:count])
		count, err = archivo.Read(buffer)
	}

	return res
}

func IterarScanners(scanner *bufio.Scanner, scanner2 *bufio.Scanner, haceAlgo func([]byte, []byte) bool) error {

	for scanner.Scan() && scanner2.Scan() && haceAlgo(scanner.Bytes(), scanner2.Bytes()) {
	}

	err := scanner.Err()
	if err == nil {
		err = scanner2.Err()
	}

	return err
}

// Por si se mueven los tests. Y para hacer mas facil escribir
func ParseameUrl(url string) string {
	url = strings.Replace(url, ":a/", "archivos/", 1)

	return strings.Replace(url, ":c/", "archivos/catedra/", 1)

}

func LeerArchivos(url string, url2 string, haceAlgo func([]byte, []byte) bool) (*bufio.Scanner, error) {
	archivo, err := os.Open(ParseameUrl(url))
	if err != nil {
		return nil, err
	}
	defer archivo.Close()

	archivo2, error2 := os.Open(ParseameUrl(url2))
	if error2 != nil {
		return nil, error2
	}

	defer archivo2.Close()
	outPutScan := bufio.NewScanner(archivo2)
	err = IterarScanners(bufio.NewScanner(archivo), outPutScan, haceAlgo)

	if err != nil {
		return nil, err
	}

	return outPutScan, nil // devuelve el restante del output ya sea para permitir mas tests, especialmente del resultado final

}

func LeerArchivo(url string, haceAlgo func([]byte) bool) error {
	archivo, error := os.Open(ParseameUrl(url))
	if error != nil {
		return error
	}

	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)

	for scanner.Scan() && haceAlgo(scanner.Bytes()) {
	}

	return scanner.Err()
}

func CrearArregloDeArchivo[T any](url string, insert func(TDALista.Lista[T], []byte) error) ([]T, error) {
	listaAux := TDALista.CrearListaEnlazada[T]()

	var err error
	errArchivo := LeerArchivo(url, func(datos []byte) bool {
		err = insert(listaAux, datos)
		return err == nil
	})

	if errArchivo != nil {
		return make([]T, 0), errArchivo
	} else if err != nil {
		return make([]T, 0), err
	}

	res := make([]T, listaAux.Largo())
	i := 0
	listaAux.Iterar(func(elemento T) bool {
		res[i] = elemento
		i++
		return true
	})

	return res, nil
}

func ordenar(arr []Votante, inicio int, final int) ([]Votante, int) {
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

func QuickSort(arr []Votante, low, high int) []Votante {
	if low < high {
		var p int
		arr, p = ordenar(arr, low, high)
		arr = QuickSort(arr, low, p-1)
		arr = QuickSort(arr, p+1, high)
	}
	return arr
}

func RedimensionarSlice[T any](viejo []T, nuevo_largo int) []T {
	nuevo := make([]T, nuevo_largo)

	copy(nuevo, viejo)

	return nuevo
}

func AccionDesdeComando(sesion SesionVotar, comando string) error {
	if comando == "fin-votar" {
		return sesion.SiguienteVotante()
	} else if comando == "deshacer" {
		return sesion.Deshacer()
	}

	args := strings.Split(comando, " ")

	if args[0] == "ingresar" {

		if len(args) < 2 { // yo pondria != pero el error es solo en falta de
			return new(ErrorFaltanParametros)
		}

		return sesion.IngresarVotante(args[1])

	}

	if args[0] == "votar" {

		if len(args) < 3 {
			return new(ErrorFaltanParametros)
		}

		return sesion.Votar(args[1], args[2])
	}

	return new(ErrorComandoInvalido)

}

func AccionComandoAString(sesion SesionVotar, comando string) string {
	err := AccionDesdeComando(sesion, comando)
	res := OK
	if err != nil {
		res = err.Error()
	}
	return res
}

// Se podria poner como primitiva pero se decidio mas elegante el no hacerlo..
// Y esta funcion serviria independientemente de la implementacion, ponele
func MostrarEstado(sesion SesionVotar, identificadores []string, mostrarLinea func(string) bool) {

	seguir := true
	i := 0

	for i < len(identificadores) && seguir {

		identificador := identificadores[i]
		if !mostrarLinea(identificador + ":") {
			return
		}

		i++

		sesion.IterarVotos(identificador, func(credencial string, votos int) bool {
			msg := fmt.Sprintf("%s: %d votos", credencial, votos) // hay que sacarle la s si es uno
			if votos == 1 {
				msg = msg[:len(msg)-1]
			}
			seguir = mostrarLinea(msg)
			return seguir
		})

		if !mostrarLinea("") {
			return
		}

	}

	if seguir {
		msg := fmt.Sprintf("Votos Impugnados: %d votos", sesion.VotosImpugnados()) // hay que sacarle la s si es uno
		if sesion.VotosImpugnados() == 1 {
			msg = msg[:len(msg)-1]
		}
		mostrarLinea(msg)
	}

}

// Tests Utilities

func TestearComando(sesion SesionVotar, comando string, expected string) error {
	res := AccionComandoAString(sesion, comando)
	if res != expected {
		return CrearErrorTest(comando, expected, res)
	}

	return nil
}
func TestearComandosScanners(sesion SesionVotar, input_scanner *bufio.Scanner, output_scanner *bufio.Scanner) error {
	var err error = nil
	IterarScanners(input_scanner, output_scanner, func(comando []byte, expected []byte) bool {
		err = TestearComando(sesion, string(comando), string(expected))
		return err == nil
	})

	return err
}

func TestearComandosArchivos(sesion SesionVotar, input_file string, out_put_file string) (*bufio.Scanner, error) {
	in, err := os.Open(ParseameUrl(input_file))
	if err != nil {
		return nil, err
	}
	defer in.Close()

	expectedOut, error2 := os.Open(ParseameUrl(out_put_file))
	if error2 != nil {
		return nil, error2
	}

	defer expectedOut.Close()

	outScanner := bufio.NewScanner(expectedOut)
	err = TestearComandosScanners(sesion, bufio.NewScanner(in), outScanner)

	if err != nil {
		return nil, err
	}

	return outScanner, err
}

func TestFinalResult(sesion SesionVotar, identificadores []string, outScanner *bufio.Scanner) error {
	errTest := sesion.Finalizar()

	if errTest != nil {
		if !outScanner.Scan() {
			return nil
		}

		if errTest.Error() != outScanner.Text() {
			return CrearErrorTest("Resultado", outScanner.Text(), errTest.Error())
		}
		errTest = nil
	}

	MostrarEstado(sesion, identificadores, func(linea string) bool {
		if !outScanner.Scan() {
			errTest = CrearErrorTest("Resultado", "EOF", linea)
		} else if outScanner.Text() != linea {
			errTest = CrearErrorTest("Resultado", outScanner.Text(), linea)
		}

		return errTest == nil
	})

	return errTest
}

func TestearArchivos(partidos_file, padron_file, input_file, out_put_file string) error {
	identificadores := []string{"Presidente", "Gobernador", "Intendente"}
	sesion, err := CrearSesion(identificadores, partidos_file, padron_file)

	if err != nil {
		return err
	}

	outScanner, errTest := TestearComandosArchivos(sesion, input_file, out_put_file)

	if errTest != nil {
		return errTest
	}
	return TestFinalResult(sesion, identificadores, outScanner)

}
