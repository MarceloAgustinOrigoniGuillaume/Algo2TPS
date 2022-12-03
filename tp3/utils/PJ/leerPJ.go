package pj

import "tp3/grafos"
import "tp3/utils"
import "errors"
import "strconv"
import "strings"
import "fmt"

func LeerPJ[T grafos.Numero](file string, grafo grafos.Grafo[string, T], convVert func(string) (string, error), convPeso func(string) (T, error)) error {
	numeroVertices := -1
	numeroAristas := -1
	var errInterno error
	var tiempo T
	var vert string
	linea := 0
	err := utils.LeerArchivo(file, func(line string) bool {
		res := ""
		grafo.IterarVertices(func(vertx string) bool {
			res += " " + vertx
			return true
		})
		linea++

		if grafo.CantidadVertices() == 0 && numeroVertices < 0 {
			numeroVertices, errInterno = strconv.Atoi(line)
		} else if numeroVertices > 0 {
			vert, errInterno = convVert(line)

			if errInterno == nil {
				grafo.AgregarVertice(vert)
				numeroVertices--
			}

		} else if grafo.CantidadAristas() == 0 && numeroAristas < 0 {
			numeroAristas, errInterno = strconv.Atoi(line)
		} else if numeroAristas > 0 {

			splitted := strings.SplitN(line, ",", 3)

			if len(splitted) != 3 {
				errInterno = errors.New(fmt.Sprintf("Conexion tiene cantidad incorrecta de datos, deberian ser 3... %v", splitted))
				return false
			}

			if !grafo.ExisteVertice(splitted[0]) {
				errInterno = errors.New(fmt.Sprintf("No pertenecia %v", splitted[0]))
			} else if !grafo.ExisteVertice(splitted[1]) {
				errInterno = errors.New(fmt.Sprintf("No pertenecia %v", splitted[1]))
			} else {
				tiempo, errInterno = convPeso(splitted[2])

				if errInterno == nil {

					grafo.AgregarArista(splitted[0], splitted[1], tiempo)
					numeroAristas--
				}

			}

		}

		return errInterno == nil
	})

	if errInterno == nil {
		return err
	}

	return errInterno
}
