package utils

import "os"
import "bufio"
import "strings"

const URL_BASE = "archivosDados/tp3/"
const URL_BASE_OUT = "../../tp3Out/"

func parseameUrl(original string) string {

	original = strings.Replace(original, ":o/", URL_BASE_OUT, 1)
	return strings.Replace(original, ":i/", URL_BASE, 1)
}

func LeerArchivo(url string, haceAlgo func(string) bool) error {
	archivo, error := os.Open(parseameUrl(url))
	if error != nil {
		return error
	}

	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)

	for scanner.Scan() && haceAlgo(scanner.Text()) {
	}

	return scanner.Err()
}

func AbrirOCrearArchivo(url string) (*os.File, error) {
	archivo, err := os.Create(parseameUrl(url))
	if err != nil {
		archivo.Close()
		return nil, err
	}

	return archivo, nil
}
