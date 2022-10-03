package sesion_votar_test

import(
	TDASesion "sesion_votar"
	"github.com/stretchr/testify/require"
	"testing"
	"fmt"
)

func crearSesionBasica() TDASesion.SesionVotar{
	return TDASesion.CrearSesion([]string{"Presidente","Gobernador","Intendente"},
						 TDASesion.BASIC_SAMPLE, TDASesion.BASIC_SAMPLE)
}

type TestPair struct{
	comando string
	expected string
}

func testearPairRequire(t *testing.T,sesion TDASesion.SesionVotar,pair TestPair){
	require.EqualValues(t,pair.expected,TDASesion.AccionComandoAString(sesion,pair.comando))
}

func testearPairRequireLog(t *testing.T,sesion TDASesion.SesionVotar,pair TestPair){
	t.Log(fmt.Sprintf("Ingresando y validando resultado comando '%s' res = '%s'",pair.comando,pair.expected))
	testearPairRequire(t,sesion,pair)
}



func testearComandosSucesionRequire(t *testing.T,sesion TDASesion.SesionVotar,pairs []TestPair){
	for _,pair := range pairs{
		//fmt.Printf("\n%s",TestPair)
		testearPairRequire(t,sesion,pair)
	}
}

func testearComandosSucesionRequireLog(t *testing.T,sesion TDASesion.SesionVotar,pairs []TestPair){
	for _,pair := range pairs{
		//fmt.Printf("\n%s",TestPair)
		testearPairRequireLog(t,sesion,pair)
	}
}


// Estas funciones de ahora son vestigios del pasado ja antes de darme cuenta podia usar el sistema de tests...
func testearPair(sesion TDASesion.SesionVotar,pair TestPair) string{
	res:= TDASesion.AccionComandoAString(sesion,pair.comando)
	if res != pair.expected{
		return fmt.Sprintf("expected '%s' got '%s'",pair.expected,res)
	}

	return ""
}

func testearComandosSucesion(sesion TDASesion.SesionVotar,pairs []TestPair){
	
	errores := ""
	
	for _,pair := range pairs{
		//fmt.Printf("\n%s",TestPair)
		error:= testearPair(sesion,pair)
		if(error != ""){
			errores+= "\n"+error
		}
	}

	if(errores == ""){
		errores = "\nTODO OK"
	}
	errores += "\n"

	fmt.Printf(errores)
}

// end vestigios



// funciones aux Desde archivos a implementar.....




func pairsDesdeArchivos(archivo_input string,archivo_output string) ([]TestPair,error){
	pares,err := TDASesion.CrearArregloDeArchivo[TestPair](archivo_input,func (bytes []byte) (TestPair,error){
			return TestPair{string(bytes),""},nil
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

	/*
	// defecto...
	return []TestPair{
			TestPair{"ingresar 1",TDASesion.OK},
			TestPair{"ingresar 2",TDASesion.OK},
			TestPair{"ingresar 200",TDASesion.ERROR_DNI_NO_ESTA},
			TestPair{"ingresar 50",TDASesion.ERROR_DNI_NO_ESTA}}

	*/
}

func testDesdeArchivosRequire(t *testing.T,candidatos_url string,padrones_url string,input_file string,out_put_file string){
	sesion := TDASesion.CrearSesion([]string{"Presidente","Gobernador","Intendente"},candidatos_url,padrones_url)
	
	testPairs, error := pairsDesdeArchivos(input_file,out_put_file)

	if(error != nil){
		t.Log(error)
		return
	}

	testearComandosSucesionRequire(t,sesion, testPairs)

	//sesion.Finalizar()
}

func testDesdeArchivosRequireLog(t *testing.T,candidatos_url string,padrones_url string,input_file string,out_put_file string){
	sesion := TDASesion.CrearSesion([]string{"Presidente","Gobernador","Intendente"},candidatos_url,padrones_url)
	
	testPairs, error := pairsDesdeArchivos(input_file,out_put_file)

	if(error != nil){
		t.Log(error)
		return
	}

	testearComandosSucesionRequireLog(t,sesion, testPairs)

	//sesion.Finalizar()
}


func testDesdeArchivosStreamRequire(t *testing.T,candidatos_url string,padrones_url string,input_file string,out_put_file string){
	sesion := TDASesion.CrearSesion([]string{"Presidente","Gobernador","Intendente"},candidatos_url,padrones_url)
	
	err := TDASesion.LeerArchivos(input_file,out_put_file, 
		func (linea_in []byte,linea_out []byte) bool {
			testearPairRequire(t,sesion,TestPair{string(linea_in),string(linea_out)})
			return true
		})

	if(err != nil){
		t.Log(err)
		return
	}

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


	testearComandosSucesionRequireLog(t,sesion, 
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

	sesion := crearSesionBasica()

	testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{"ingresar 1",TDASesion.OK},
			TestPair{"ingresar 2",TDASesion.OK},
			TestPair{"ingresar 200",TDASesion.ERROR_DNI_NO_ESTA},
			TestPair{"ingresar 50",TDASesion.ERROR_DNI_NO_ESTA}})	

}

func TestVotoEnBlanco(t *testing.T){
	t.Log("Se va a votar solo a presidente e intendente, eso es voto en blanco para todos, ya que no se voto en su totalidad")

	sesion := crearSesionBasica()


	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})

	testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{"votar Presidente 0",TDASesion.OK}, // capaz este deberia dar error tmbn
			TestPair{"votar Intendente 1",TDASesion.OK}, 
			TestPair{"votar Presidente 1",TDASesion.OK},
			TestPair{"fin-votar",TDASesion.OK}})

	i:= 0
	expected := []int{1,0,0,0}
	sesion.IterarVotos("Presidente",func (credencial string,votos int) bool{
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
			TestPair{"votar Presidente 0",TDASesion.OK}, // capaz este deberia dar error tmbn
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
	

	i:= 0
	expected := []int{0,1,0,9}
	sesion.IterarVotos("Presidente",func (credencial string,votos int) bool{
		require.EqualValues(t,expected[i],votos)
		i++
		return true
	})

	require.EqualValues(t,nil,sesion.Finalizar())

}



func TestFinVotar(t *testing.T){
	sesion := crearSesionBasica()


	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})

	testearComandosSucesionRequire(t,sesion, 
			[]TestPair{
			TestPair{"votar Presidente 0",TDASesion.OK}, // capaz este deberia dar error tmbn
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
	

	i:= 0
	expected := []int{0,1,0,0}
	sesion.IterarVotos("Presidente",func (credencial string,votos int) bool{
		require.EqualValues(t,expected[i],votos)
		i++
		return true
	})

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
	

	i:= 0
	expected := []int{0,10,0,0}
	sesion.IterarVotos("Presidente",func (credencial string,votos int) bool{
		require.EqualValues(t,expected[i],votos)
		i++
		return true
	})

	expected[3] = 10
	expected[1] = 0
	i=0
	sesion.IterarVotos("Intendente",func (credencial string,votos int) bool{
		require.EqualValues(t,expected[i],votos)
		i++
		return true
	})
}



// Fraudes y otras verificaciones de respuestas
func TestFraudulentos(t *testing.T){
	sesion := crearSesionBasica()


	testearPairRequireLog(t,sesion,TestPair{"ingresar 1",TDASesion.OK})
	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})

	testearPairRequire(t,sesion,TestPair{"votar Presidente 0",TDASesion.OK})
	testearPairRequireLog(t,sesion,TestPair{"fin-votar",TDASesion.OK})

	testearPairRequireLog(t,sesion,TestPair{"votar Presidente 0",fmt.Sprintf(TDASesion.ERROR_VOTANTE_FRAUDULENTO,1)})
	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})
	testearPairRequireLog(t,sesion,TestPair{"deshacer",fmt.Sprintf(TDASesion.ERROR_VOTANTE_FRAUDULENTO,1)})
	testearPairRequire(t,sesion,TestPair{"ingresar 1",TDASesion.OK})
	testearPairRequireLog(t,sesion,TestPair{"fin-votar",fmt.Sprintf(TDASesion.ERROR_VOTANTE_FRAUDULENTO,1)})

	require.EqualValues(t,3,sesion.VotosImpugnados())

}


// Test de resultados finales, con logica algo mas compleja

func TestsFuncionales(t *testing.T){

}





func TestDesdeArchivos(t *testing.T){
	testDesdeArchivosRequireLog(t,"../archivos/set1/listas","../archivos/set1/padrones","../archivos/set1/in","../archivos/set1/out")
	testDesdeArchivosStreamRequire(t,"../archivos/set1/listas","../archivos/set1/padrones","../archivos/set1/in","../archivos/set1/out")
}