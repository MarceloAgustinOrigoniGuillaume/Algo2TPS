package main

import (
	"bufio"
	"fmt"
	"os"
	TDASesion "sesion_votar"
	errores "sesion_votar/errores"
	"time"
)

func test(partidos, padron, in, out string) {
	fmt.Fprintf(os.Stdout, "test %s , %s \n", partidos, padron)
	fmt.Fprintf(os.Stdout, "in: %s , expected : %s \n", in, out)

	start := time.Now()
	defer fmt.Fprintf(os.Stdout, "\ntook %s", time.Since(start))
	err := TDASesion.TestearArchivos(partidos, padron, in, out)

	res := "TODO OK"

	if err != nil {
		res = err.Error()
	}
	fmt.Fprintf(os.Stdout, res)
}

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stdout, errores.ERROR_FALTAN_PARAMETROS+"\n")
		return
	}

	TIPOS_VOTOS := []string{"Presidente", "Gobernador", "Intendente"}

	// dios este con nosotros pero se quiso hacerlo para hacer un uso mas facil desde consola
	// para testear desde consola si se quiere
	if os.Args[1] == "-test:" && len(os.Args) == 3 {
		// formato de archivos de la catedra
		test(os.Args[2]+"_partidos", os.Args[2]+"_padron", os.Args[2]+"_in", os.Args[2]+"_out")
		return
	} else if len(os.Args) == 5 {
		test(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
		return
	}

	sesion, err := TDASesion.CrearSesion(TIPOS_VOTOS, os.Args[1], os.Args[2])

	if err != nil {
		fmt.Fprintf(os.Stdout, err.Error()+"\n")
		return
	}

	inputUsuario := bufio.NewScanner(os.Stdin)

	for inputUsuario.Scan() {
		fmt.Fprintf(os.Stdout, "%s\n", TDASesion.AccionComandoAString(sesion, inputUsuario.Text()))
	}

	res := sesion.Finalizar()

	if res != nil {
		fmt.Fprintf(os.Stdout, "%s\n", res.Error())
	}

	TDASesion.MostrarEstado(sesion, TIPOS_VOTOS, func(mensaje string) bool {
		fmt.Fprintf(os.Stdout, mensaje+"\n")
		return true
	})

	inputUsuario.Scan() // para que pause si es por consola a mano
}
