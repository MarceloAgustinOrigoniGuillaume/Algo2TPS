package main

import "os"
import algogram "tp2/sesion"
import "strings"
import "bufio"
import "fmt"

const ErrorFaltanParametros = "Faltan parametros"
const ErrorComandoInvalido = "Comando invalido"

func devolverErrorODefault(def string, err error) string {
	if err != nil {
		return err.Error()
	}

	return def
}

func AccionDesdeComando(sesion algogram.Sesion, comando string) string {
	if comando == "logout" {
		return devolverErrorODefault("Adios", sesion.Logout())
	} else if comando == "ver_siguiente_feed" {
		return devolverErrorODefault(sesion.VerSiguientePost())
	}

	args := strings.SplitN(comando, " ", 2)

	if len(args) < 2 {
		return ErrorFaltanParametros
	}

	if args[0] == "login" {
		return devolverErrorODefault("Hola "+args[1], sesion.Login(args[1]))
	}

	if args[0] == "publicar" {
		return devolverErrorODefault("Post publicado", sesion.Publicar(args[1]))
	}

	if args[0] == "likear_post" {
		return devolverErrorODefault("Post likeado", sesion.Likear(args[1]))
	}

	if args[0] == "mostrar_likes" {
		return devolverErrorODefault(sesion.MostrarLikes(args[1]))
	}

	return ErrorComandoInvalido
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "INSUFICIENTES ARGS\n")
		return
	}

	sesion, err := algogram.CrearSesion(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err.Error())
		return
	}
	inputUsuario := bufio.NewScanner(os.Stdin)

	for inputUsuario.Scan() {
		fmt.Fprintf(os.Stdout, "%s\n", AccionDesdeComando(sesion, inputUsuario.Text()))
	}

}
