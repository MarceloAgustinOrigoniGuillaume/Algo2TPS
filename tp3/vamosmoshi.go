package main

import "tp3/agencia_viajes"
import "os"
import "fmt"
import "bufio"
import "strings"

const ErrorParametros = "Cantidad parametros incorrecta"
const ErrorComandoInvalido = "Comando invalido"

func devolverErrorODefault(def string, err error) string {
	if err != nil {
		return err.Error()
	}

	return def
}

func AccionDesdeComando(empresa agencia.AgenciaViajes, comando string) string {
	args := strings.SplitN(comando, " ", 2)

	if len(args) < 2 {
		return ErrorParametros
	}

	comando = args[0]

	if comando == "itinerario" {
		return devolverErrorODefault(empresa.Itinerario(args[1]))
	}

	if comando == "reducir_caminos" {
		return devolverErrorODefault(empresa.ReducirCaminos(args[1]))
	}

	argsSplitted := strings.Split(args[1], ",")

	if comando == "viaje" {
		if len(argsSplitted) != 2 {
			return ErrorParametros
		}
		return devolverErrorODefault(empresa.ViajeDesde(strings.TrimSpace(argsSplitted[0]), strings.TrimSpace(argsSplitted[1])))
	}

	if comando == "ir" {
		if len(argsSplitted) != 3 {
			return ErrorParametros
		}
		return devolverErrorODefault(empresa.Ir(strings.TrimSpace(argsSplitted[0]), strings.TrimSpace(argsSplitted[1]), strings.TrimSpace(argsSplitted[2])))
	}

	return ErrorComandoInvalido
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "INSUFICIENTES ARGS\n")
		return
	}

	empresa, err := agencia.CrearAgenciaViajes(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error cargando el lugar: %s\n", err.Error())
		return
	}
	inputUsuario := bufio.NewScanner(os.Stdin)

	for inputUsuario.Scan() {
		fmt.Fprintf(os.Stdout, "%s\n", AccionDesdeComando(empresa, inputUsuario.Text()))
	}

}
