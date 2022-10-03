package main

import (
	TDASesion "sesion_votar"
	"fmt"
	"os"
	"strconv"
)

var TIPOS_VOTOS []string

func ejecutarComandosSucesion( sesion TDASesion.SesionVotar,comandos []string){
	for _,comando := range comandos{
		fmt.Printf("\n%s",comando)
		err:= TDASesion.AccionDesdeComando(sesion,comando)
		res:= TDASesion.OK
		if(err != nil){
			res = err.Error()
		}
		fmt.Printf("\n%s",res)
	}
}


func main(){

	if(len(os.Args) < 3){
		fmt.Printf("Insuficientes argumentos")
		return
	}

	TIPOS_VOTOS = []string{"Presidente","Gobernador","Intendente"}

	sesion := TDASesion.CrearSesion(TIPOS_VOTOS,os.Args[1],os.Args[2])


	if(len(os.Args)>3){
		//
		fmt.Printf("DEBERIA TESTEAR ARCHIVOS")

		for i:= 3;i<len(os.Args);i++{
			fmt.Printf("\n"+os.Args[i]+"\n")
			arreglo,err := TDASesion.CrearArregloDeArchivo[int](os.Args[i],func (inp []byte) (int,error){
				text:= string(inp)
				res,err := strconv.Atoi(text)
				if(err == nil){
					fmt.Printf("\nNew elem %d",res)
				} else{
					fmt.Printf("\n Omitio :"+text)
					err = new(TDASesion.ErrorOmicion)
				}
				return res,err
			})

			if(err != nil){
				fmt.Printf("\n"+err.Error()+"\n")
			} else{
				fmt.Printf("\n len final == %d , cap = %d\n",len(arreglo),cap(arreglo))
			}
			//TestFromArchivo(os.Args[i])
		}
		return
	}

	// Aca deberia encargarse de pedir input por consola no este test simple
	ejecutarComandosSucesion(sesion, []string{"ingresar 1","votar Presidente 2","votar Presidente 3",
												"deshacer","votar Gobernador 3","ingresar 35","ingresar -3","fin-votar"})


	ejecutarComandosSucesion(sesion, []string{"ingresar 1","ingresar 2","votar Presidente 2","votar Presidente 3",
												"deshacer","votar Gobernador 3","ingresar 35","fin-votar"})


	ejecutarComandosSucesion(sesion,[]string{"ingresar 1","fin-votar"})
	ejecutarComandosSucesion(sesion,[]string{"ingresar 1","deshacer"})

	res:= sesion.Finalizar()

	if(res != nil){
		fmt.Printf("\n"+res.Error())
	}

	TDASesion.MostrarEstado(sesion,TIPOS_VOTOS)


}