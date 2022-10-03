package main

import (
	TDASesion "sesion_votar"
	"fmt"
	"os"
	"bufio"
)

const ERROR_ARGUMENTOS = "Insuficientes argumentos"

func ejecutarComandosSucesion( sesion TDASesion.SesionVotar,comandos []string){
	for _,comando := range comandos{
		fmt.Fprintf(os.Stdout,"\n%s",comando)
		fmt.Fprintf(os.Stdout,"\n%s",TDASesion.AccionComandoAString(sesion,comando))
	}
}


func main(){

	if(len(os.Args) < 3){
		fmt.Fprintf(os.Stdout,ERROR_ARGUMENTOS)
		return
	}

	TIPOS_VOTOS := []string{"Presidente","Gobernador","Intendente"}

	//fmt.Fprintf(os.Stdout,"SESION VOTAR INITED")

	sesion,err := TDASesion.CrearSesion(TIPOS_VOTOS,os.Args[1],os.Args[2])

	if(err != nil){
		fmt.Fprintf(os.Stdout,err.Error())
		return
	}

	if(len(os.Args)== 5){ // funcionalidad extra para testear desde consola si se quiere
		fmt.Fprintf(os.Stdout,"Test->in: %s , expected : %s \n",os.Args[3],os.Args[4])
		err:= TDASesion.TestearComandosArchivos(sesion,os.Args[3],os.Args[4])
		
		res:= "TODO OK"
		if(err != nil){
			res = err.Error()
		}

		fmt.Fprintf(os.Stdout,res)
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

	/*

	// Hard coded commands
	ejecutarComandosSucesion(sesion, []string{"ingresar 1","votar Presidente 2","votar Presidente 3",
												"deshacer","votar Gobernador 3","ingresar 35","ingresar -3","fin-votar"})


	ejecutarComandosSucesion(sesion, []string{"ingresar 1","ingresar 2","votar Presidente 2","votar Presidente 3",
												"deshacer","votar Gobernador 3","ingresar 35","fin-votar"})


	ejecutarComandosSucesion(sesion,[]string{"ingresar 1","fin-votar"})
	ejecutarComandosSucesion(sesion,[]string{"ingresar 1","deshacer"})

	

	

	*/
}