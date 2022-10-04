package sesion_votar_test

import(
	TDASesion "sesion_votar"
	TDALista "lista"
	"github.com/stretchr/testify/require"
	"testing"
	"fmt"
	"os"
)
os.
func crearSesionBasica() TDASesion.SesionVotar{
	sesion,_:= TDASesion.CrearSesion([]string{"Presidente","Gobernador","Intendente"},
						 TDASesion.BASIC_SAMPLE, TDASesion.BASIC_SAMPLE)
	return sesion
}

type TestPair struct{
	comando string
	expected string
}

func testearPairRequire(t *testing.T,sesion TDASesion.SesionVotar,pair TestPair){
	require.EqualValues(t,pair.expected,TDASesion.AccionComandoAString(sesion,pair.comando))
}

func testearPairRequireLog(t *testing.T,sesion TDASesion.SesionVotar,pair TestPair){
	t.Log(fmt.Sprintf("->%s res = '%s'",pair.comando,pair.expected))
	testearPairRequire(t,sesion,pair)
}



func testearComandosSucesionRequire(t *testing.T,sesion TDASesion.SesionVotar,pairs []TestPair){
	for _,pair := range pairs{
		//fmt.Printf("%s",TestPair)
		testearPairRequire(t,sesion,pair)
	}
}

func testearComandosSucesionRequireLog(t *testing.T,sesion TDASesion.SesionVotar,pairs []TestPair){
	for _,pair := range pairs{
		//fmt.Printf("%s",TestPair)
		testearPairRequireLog(t,sesion,pair)
	}
}




// funciones aux Desde archivos a implementar.....




func pairsDesdeArchivos(archivo_input string,archivo_output string) ([]TestPair,error){
	

	pares,err := TDASesion.CrearArregloDeArchivo[TestPair](archivo_input,
		func (lista TDALista.Lista[TestPair],bytes []byte) error{
			lista.InsertarUltimo(TestPair{string(bytes),""})
			return nil
		})

	if(err == nil){
		i:= 0
		err = TDASesion.LeerArchivo(archivo_output,func (bytes []byte) bool{
			if(i >= len(pares)){
				i++ // para despues poder saber si hubo mas que comandos
				return false
			}

			pares[i].expected = string(bytes)
			i++
			return true
		})

		if(err == nil){
			if(i>len(pares)){
				err = new(TDASesion.ErrorMissMatchSizeOut)
			} else if (i < len(pares)){
				err = new(TDASesion.ErrorMissMatchSizeIn)
				copy_pares := make([]TestPair,i)
				copy(copy_pares,pares)
				pares = copy_pares

			}
		}
	}

	

	return pares,err

}

func testDesdeArchivosRequire(t *testing.T,candidatos_url string,padrones_url string,input_file string,out_put_file string){
	sesion,err := TDASesion.CrearSesion([]string{"Presidente","Gobernador","Intendente"},candidatos_url,padrones_url)
	
	if(err != nil){
		t.Log(err)
		return
	}

	testPairs, err2 := pairsDesdeArchivos(input_file,out_put_file)

	if(err2 != nil){
		t.Log(err2)
		return
	}

	testearComandosSucesionRequire(t,sesion, testPairs)

	//sesion.Finalizar()
}



func testDesdeArchivosRequireLog(t *testing.T,candidatos_url string,padrones_url string,input_file string,out_put_file string){
	sesion,err := TDASesion.CrearSesion([]string{"Presidente","Gobernador","Intendente"},candidatos_url,padrones_url)

	
	if(err != nil){
		t.Log(err)
		return
	}

	
	testPairs, error := pairsDesdeArchivos(input_file,out_put_file)

	if(error != nil){
		t.Log(error)
		return
	}

	testearComandosSucesionRequireLog(t,sesion, testPairs)

	//sesion.Finalizar()
}


func testDesdeArchivosStreamRequire(t *testing.T,candidatos_url string,padrones_url string,input_file string,out_put_file string) error{
	sesion,err := TDASesion.CrearSesion([]string{"Presidente","Gobernador","Intendente"},candidatos_url,padrones_url)
	
	if(err == nil){
		err = TDASesion.LeerArchivos(input_file,out_put_file, 
		func (linea_in []byte,linea_out []byte) bool {
			testearPairRequire(t,sesion,TestPair{string(linea_in),string(linea_out)})
			return true
		})	
	}
	

	return err

	//sesion.Finalizar()
}





// TESTS UNITARIOS O CERCANOS A ESO JA


// De funcionalidad general
func TestVacio(t *testing.T){
	sesion := crearSesionBasica()


	require.False(t,sesion.HayVotante())
	testearPairRequireLog(t,sesion,TestPair{"votar Presidente 1",TDASesion.ERROR_FILA_VACIA})
	testearPairRequireLog(t,sesion,TestPair{"fin-votar",TDASesion.ERROR_FILA_VACIA})
	testearPairRequireLog(t,sesion,TestPair{"deshacer",TDASesion.ERROR_FILA_VACIA})	

}

func TestIngresadoComandos(t *testing.T){
	sesion := crearSesionBasica()
	t.Log("Se va a probar que los comandos se puedan ejecutar, si no se es invalido y no faltan parametros")


	testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{"1",TDASesion.ERROR_COMANDO_INVALIDO},
			TestPair{"ingresar 2",TDASesion.OK},
			TestPair{"ingresar",TDASesion.ERROR_FALTAN_PARAMETROS},
			TestPair{"votar",TDASesion.ERROR_FALTAN_PARAMETROS},
			TestPair{"votar Presidente",TDASesion.ERROR_FALTAN_PARAMETROS},
			TestPair{"votar Presidente 1",TDASesion.OK},
			TestPair{"deshacer",TDASesion.OK},
			TestPair{"fin-votar",TDASesion.OK},
				})	

}


// De funcionalidad de comandos
func TestIngresarVotante(t *testing.T){
	t.Log("Se va a probar probando ingresar dnis invalidos y validos")

	sesion := crearSesionBasica()

	testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{"ingresar 1",TDASesion.OK},
			TestPair{"ingresar 2",TDASesion.OK},
			TestPair{"ingresar 200",TDASesion.ERROR_DNI_NO_ESTA},
			TestPair{"ingresar -2",TDASesion.ERROR_DNI_INVALIDO},
			TestPair{"ingresar 50",TDASesion.ERROR_DNI_NO_ESTA}})	

}

func TestVotoEnBlanco(t *testing.T){
	t.Log("Se va a votar solo a presidente e intendente, eso es voto en blanco para todos, ya que no se voto en su totalidad")

	sesion := crearSesionBasica()

	testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{"ingresar 1",TDASesion.OK},
			TestPair{"votar Intendente 1",TDASesion.OK}, 
			TestPair{"votar Presidente 1",TDASesion.OK},
			TestPair{"fin-votar",TDASesion.OK}})

	verificarVotos(t,sesion,"Presidente",[]int{1,0,0,0})

}

func verificarVotos(t *testing.T, sesion TDASesion.SesionVotar,tipo string,expected []int){
	i:= 0
	sesion.IterarVotos(tipo,func (credencial string,votos int) bool{
		if(i>= len(expected)){
			return false
		}
		require.EqualValues(t,expected[i],votos)
		i++
		return true
	})
}

func TestVotar(t *testing.T){
	sesion := crearSesionBasica()


	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})

	testearComandosSucesionRequireLog(t,sesion, 
			[]TestPair{
			TestPair{"votar Gobernador 1",TDASesion.OK}, 
			TestPair{"votar Intendente 1",TDASesion.OK}, 
			TestPair{"votar Presidente 4",TDASesion.ERROR_ALTERNATIVA_INVALIDA},
			TestPair{"votar Diputado 2",TDASesion.ERROR_TIPO_INVALIDO},
			TestPair{"votar Presidente 1",TDASesion.OK}})

	testearPairRequire(t,sesion,TestPair{"fin-votar",TDASesion.OK})


	t.Log("Se va a votar al presidente 1 y despues al 3 9 veces y verificar los votos al final")
	for i:=2;i<11;i++{
		testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{fmt.Sprintf("ingresar %d",i),TDASesion.OK},
			TestPair{"votar Gobernador 1",TDASesion.OK}, 
			TestPair{"votar Intendente 1",TDASesion.OK},
			TestPair{"votar Presidente 1",TDASesion.OK},
			TestPair{"votar Presidente 3",TDASesion.OK},
			TestPair{"fin-votar",TDASesion.OK}})
	}
	

	verificarVotos(t,sesion,"Presidente",[]int{0,1,0,9})	
	require.EqualValues(t,nil,sesion.Finalizar())

}

func TestVotosImpugnados(t *testing.T){
	sesion := crearSesionBasica()

	t.Log("Se va a votar al Gobernador 0 y despues varios otros votos se verificara este impugnado")
	t.Log("Tambien Se va a votar al Presidente 1 y Gobernador 0 y despues al 2 y 1 respectivamente.Despues deshacer y despues varios otros votos se verificara los votos")

	testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{"ingresar 1",TDASesion.OK},
			TestPair{"ingresar 2",TDASesion.OK},
			TestPair{"votar Gobernador 0",TDASesion.OK}, 
			TestPair{"votar Gobernador 1",TDASesion.OK}, 
			TestPair{"votar Intendente 1",TDASesion.OK}, 
			TestPair{"votar Intendente 2",TDASesion.OK}, 
			TestPair{"votar Presidente 3",TDASesion.OK},
			TestPair{"fin-votar",TDASesion.OK},
			TestPair{"votar Presidente 1",TDASesion.OK},
			TestPair{"votar Gobernador 0",TDASesion.OK},
			TestPair{"votar Presidente 2",TDASesion.OK}, 
			TestPair{"votar Gobernador 1",TDASesion.OK}, 
			TestPair{"deshacer",TDASesion.OK}, 
			TestPair{"votar Gobernador 2",TDASesion.OK}, 
			TestPair{"votar Intendente 1",TDASesion.OK}, 
			TestPair{"votar Intendente 2",TDASesion.OK},
			TestPair{"fin-votar",TDASesion.OK}})

	verificarVotos(t,sesion,"Presidente",[]int{0,1,0,0})
	verificarVotos(t,sesion,"Intendente",[]int{0,0,1,0})
	verificarVotos(t,sesion,"Gobernador",[]int{0,0,1,0})
	require.EqualValues(t,1,sesion.VotosImpugnados())
}

func TestFinVotar(t *testing.T){
	sesion := crearSesionBasica()


	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})

	testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{"votar Gobernador 1",TDASesion.OK}, 
			TestPair{"votar Intendente 1",TDASesion.OK}, 
			TestPair{"votar Presidente 4",TDASesion.ERROR_ALTERNATIVA_INVALIDA},
			TestPair{"votar Diputado 2",TDASesion.ERROR_TIPO_INVALIDO},
			TestPair{"votar Presidente 1",TDASesion.OK}})

	testearPairRequire(t,sesion,TestPair{"fin-votar",TDASesion.OK})


	t.Log("Se va a votar al presidente 1 y despues al 3 9 veces pero sin llamar fin-votar y verificar que no se haya votado y salte el mensaje que no se termino de votar")
	for i:=2;i<11;i++{
		testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{fmt.Sprintf("ingresar %d",i),TDASesion.OK},
			TestPair{"votar Gobernador 1",TDASesion.OK}, 
			TestPair{"votar Intendente 1",TDASesion.OK},
			TestPair{"votar Presidente 1",TDASesion.OK},
			TestPair{"votar Presidente 3",TDASesion.OK},})
	}
	

	verificarVotos(t,sesion,"Presidente",[]int{0,1,0,0})


	err:= sesion.Finalizar()
	require.NotNil(t,err)
	require.EqualValues(t,TDASesion.ERROR_SIN_TERMINAR,err.Error())
}


func TestDeshacer(t *testing.T){
	sesion := crearSesionBasica()

	t.Log("Se va a votar al presidente y gobernador 1 intendente 3 y despues al presidente 3 intendente 1 despues deshacer dos veces... 10 veces y verificar los votos al final")

	for i:=1;i<11;i++{
		testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{fmt.Sprintf("ingresar %d",i),TDASesion.OK},
			TestPair{"votar Presidente 1",TDASesion.OK},
			TestPair{"votar Gobernador 1",TDASesion.OK},
			TestPair{"votar Intendente 3",TDASesion.OK},
			TestPair{"votar Intendente 1",TDASesion.OK},
			TestPair{"votar Presidente 3",TDASesion.OK},
			TestPair{"deshacer",TDASesion.OK},
			TestPair{"deshacer",TDASesion.OK},
			TestPair{"fin-votar",TDASesion.OK}})
	}
	

	verificarVotos(t,sesion,"Presidente",[]int{0,10,0,0})
	verificarVotos(t,sesion,"Intendente",[]int{0,0,0,10})
}



func TestFraudulentos(t *testing.T){
	sesion := crearSesionBasica()

	testearPairRequireLog(t,sesion,TestPair{"ingresar 1",TDASesion.OK})
	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})
	testearPairRequire(t,sesion,TestPair{"votar Presidente 2",TDASesion.OK})
	testearPairRequireLog(t,sesion,TestPair{"fin-votar",TDASesion.OK})

	t.Log("Se va ingresar devuelta 1 y testear el fraude")


	testearPairRequireLog(t,sesion,TestPair{"votar Presidente 0",fmt.Sprintf(TDASesion.ERROR_VOTANTE_FRAUDULENTO,1)})
	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})
	testearPairRequireLog(t,sesion,TestPair{"deshacer",fmt.Sprintf(TDASesion.ERROR_VOTANTE_FRAUDULENTO,1)})
	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})
	testearPairRequireLog(t,sesion,TestPair{"fin-votar",fmt.Sprintf(TDASesion.ERROR_VOTANTE_FRAUDULENTO,1)})

	require.EqualValues(t,0,sesion.VotosImpugnados())

}

func TestDesdeArchivos(t *testing.T){
	t.Log("Se verificara que se pueda cargar el sistema desde archivos")
	testDesdeArchivosRequire(t,"../archivos/set1/partidos","../archivos/set1/padron","../archivos/set1/in","../archivos/set1/out")
	err:= testDesdeArchivosStreamRequire(t,"../archivos/set1/partidos","../archivos/set1/padron","../archivos/set1/in","../archivos/set1/out")

	if(err != nil){
		t.Log(err.Error())
	}
}


func getUrlBaseCatedra(num_test int) string{
	return fmt.Sprintf("../:c/%02d",num_test)
}

func TestCatedra(t *testing.T){
	t.Log("Se va a testear de forma iterativa los tests de la catedra")

	for i:= 1;i<11;i++{
		url := getUrlBaseCatedra(i)
		archivo,err := os.Open(TDASesion.ParseameUrl(url+".test"))
		if(err != nil){
			t.Log(err.Error())
			continue
		}
		t.Log(fmt.Sprintf("test:%d -----%s\n",i,TDASesion.ReadAll(archivo)))
		archivo.Close()


		err = testDesdeArchivosStreamRequire(t,url+"_partidos",url+"_padron",url+"_in",url+"_out")		

		if(err != nil){
			t.Log(err.Error())
		} else{
			t.Log("PASS")			
		}
	}



}
