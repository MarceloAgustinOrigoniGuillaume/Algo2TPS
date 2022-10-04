package main

import (
	TDASesion "sesion_votar"
	"fmt"
	"os"
	"bufio"
	"time"
)

const ERROR_ARGUMENTOS = "Insuficientes argumentos"

func ejecutarComandosSucesion( sesion TDASesion.SesionVotar,comandos []string){
	for _,comando := range comandos{
		fmt.Fprintf(os.Stdout,"\n%s",comando)
		fmt.Fprintf(os.Stdout,"\n%s",TDASesion.AccionComandoAString(sesion,comando))
	}
}

func test(partidos, padron, in, out string){
	fmt.Fprintf(os.Stdout,"test %s , %s \n",partidos,padron)
	start:= time.Now()
	defer fmt.Fprintf(os.Stdout,"\ntook %s",time.Since(start))

	sesion,err := TDASesion.CrearSesion([]string{"Presidente","Gobernador","Intendente"},partidos,padron)

	if(err != nil){
		fmt.Fprintf(os.Stdout,err.Error())
		return
	}

	fmt.Fprintf(os.Stdout,"in: %s , expected : %s \n",in,out)
	err = TDASesion.TestearComandosArchivos(sesion,in,out)
	res:= "TODO OK"
	if(err != nil){
		res = err.Error()
	}

	fmt.Fprintf(os.Stdout,res)
}


func main(){

	if(len(os.Args) < 3){
		fmt.Fprintf(os.Stdout,ERROR_ARGUMENTOS)
		return
	}

	TIPOS_VOTOS := []string{"Presidente","Gobernador","Intendente"}
	
	// dios me perdone pero quise hacerlo para hacer mejor uso desde consola
	// funcionalidad extra para testear desde consola si se quiere
	if(os.Args[1] == "-test:" && len(os.Args) == 3){ 
		test(os.Args[2]+"_partidos",os.Args[2]+"_padron",os.Args[2]+"_in",os.Args[2]+"_out")
		return
	} else if(len(os.Args)== 5){
		test(os.Args[1],os.Args[2],os.Args[3],os.Args[4])
		return
	}

	sesion,err := TDASesion.CrearSesion(TIPOS_VOTOS,os.Args[1],os.Args[2])

	if(err != nil){
		fmt.Fprintf(os.Stdout,err.Error())
		return
	}

	inputUsuario := bufio.NewScanner(os.Stdin)

	for inputUsuario.Scan(){
		comando := inputUsuario.Text()
		if(comando == ""){ // esto es para cuando se hace por consola
			break
		}

		fmt.Fprintf(os.Stdout,"%s\n",TDASesion.AccionComandoAString(sesion,comando))
	}

	sesion.Finalizar()

	res:= sesion.Finalizar()

	if(res != nil){
		fmt.Fprintf(os.Stdout,"\n%s",res.Error())
	}

	TDASesion.MostrarEstado(sesion,TIPOS_VOTOS)
}