package sesion
import "os"
import "bufio"

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
