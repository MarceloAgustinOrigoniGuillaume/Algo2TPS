package utilities

import "os"
import "bufio"
import "strings"

// ../../pruebaTp2/users_test.txt
func LeerArchivo(url string, haceAlgo func(string) bool) error {
	archivo, error := os.Open(url)
	if error != nil {
		return error
	}

	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)

	for scanner.Scan() && haceAlgo(scanner.Text()) {
	}

	return scanner.Err()
}

func CompareLexico(user, other string) int {
	nombre := strings.ToLower(user)
	nombreOtro := strings.ToLower(other)
	toReturn := -1

	min := len(nombre)
	if len(nombreOtro) < min {
		toReturn = 1
		min = len(nombreOtro)
	} else if len(nombreOtro) == min {
		toReturn = 0
	}

	for i := 0; i < min; i++ {
		if nombre[i] < nombreOtro[i] {
			return -1
		} else if nombre[i] > nombreOtro[i] {
			return 1
		}
	}

	return toReturn

}
