package sesion_votar

import (
	TDACola "cola"
	TDALista "lista"
	"strconv"
	"strings"
)

type sesionVotar struct{
	identificadores_tipos []string
	listaCandidatos [][]candidatoStruct
	aVotar []Votante
	esperandoAVotar TDACola.Cola[Votante]
	registroDeVotos Registro
	votosImpugnados int
}

//Funcion crear

func CrearSesion(identificadores_tipos []string,candidatos_file string,padrones_file string) (SesionVotar,error){
	sesion := new(sesionVotar)
	sesion.identificadores_tipos = identificadores_tipos

	sesion.esperandoAVotar = TDACola.CrearColaEnlazada[Votante]()
	sesion.registroDeVotos = CrearRegistroDeVotos()

	// valores iniciales por si saltan errores
	sesion.listaCandidatos = make([][]candidatoStruct,0) 
	sesion.aVotar = make([]Votante,0)

	err := sesion.popularCandidatos(candidatos_file)
	

	if(err == nil){
		err = sesion.popularVotantes(padrones_file,len(identificadores_tipos))
	}

	

	if(err != nil){
		err = new(ErrorLecturaArchivos)
	}
	


	return sesion,err
}

// funciones auxiliares


// devuelve un arreglo de Votantes dado archivo de los dnis
// Estaran ordenados por dni de menor a mayor
func (sesion *sesionVotar) popularVotantes(archivo string,opciones int) error{
	if(padrones_file == BASIC_SAMPLE){
		sesion.aVotar = popularVotantesBasico(len(identificadores_tipos))
		return nil
	} 

	sesion.aVotar = make([]Votante,1024)
	i:= 0
	errArchivo := LeerArchivo(archivo,func (datos []byte) bool{
			dni,err := strconv.Atoi(string(datos))
			if(err != nil){
				return true
			}

			if(i == len(sesion.aVotar)){
				sesion.aVotar= RedimensionarSlice(sesion.aVotar,2*len(sesion.aVotar))
			}
			sesion.aVotar[i] = CrearVotante(dni,opciones)
			i++
		return err == nil
	})

	if(errArchivo != nil){
		return errArchivo
	}

	if(i != len(sesion.aVotar)){
		sesion.aVotar= RedimensionarSlice(sesion.aVotar,i)
	}

	sesion.aVotar = QuickSort(sesion.aVotar,0,i-1)

	return nil
}

// crea una lista dada la cantidad de tipos de votos y un archivo de los candidatos, devuelve si hay error
func (sesion *sesionVotar) popularCandidatos(archivo string) error{
	if(candidatos_file == BASIC_SAMPLE){
		sesion.listaCandidatos = popularCandidatosBasico(len(sesion.identificadores_tipos)+1)
		return nil
	}

	tipos:= len(sesion.identificadores_tipos)

	candidatosArchivo,errArchivo := CrearArregloDeArchivo(archivo,func (lista TDALista.Lista[[]candidatoStruct],bytes []byte) error{
				splitted := strings.Split(string(bytes),",")
				if(len(splitted) < tipos+1){ // + 1 por el partido
					return new(ErrorLecturaArchivos)
				}

				candidatosPartido := make([]candidatoStruct,tipos)
				i:= 1
				for ind_candidato:= range candidatosPartido{
					candidatosPartido[ind_candidato] = CrearCandidato(splitted[0],splitted[i])
					i++
				}				

				lista.InsertarUltimo(candidatosPartido)
				return nil
			})

	if(errArchivo != nil){
		return errArchivo
	}

	sesion.listaCandidatos = make([][]candidatoStruct,len(candidatosArchivo)+1)

	sesion.listaCandidatos[0] = make([]candidatoStruct,tipos)

	for i:= range sesion.listaCandidatos[0]{
		sesion.listaCandidatos[0][i] = candidatoStruct{}
	}

	copy(sesion.listaCandidatos[1:],candidatosArchivo)

	return nil

}


//Busca binariamente, aprovechando popularVotantes.
func buscarVotante(lista_dni []Votante,dni int) Votante { 
		
	inicio := 0
	fin := len(lista_dni)
	medio := (fin+inicio)/2


	for medio != inicio && lista_dni[inicio].DNI()<dni && lista_dni[fin-1].DNI()>dni{
		if(lista_dni[medio].DNI() == dni){
			return lista_dni[medio]
		}

		if(lista_dni[medio].DNI() > dni){
			fin = medio
		} else{
			inicio = medio+1
		}

		medio = (fin+inicio)/2
	} 
	
	if(lista_dni[inicio].DNI() == dni){
		return lista_dni[inicio]
	} 

	if(lista_dni[fin-1].DNI() == dni){
		return lista_dni[fin-1]
	} 

	
	return nil
}


// funcion auxiliar para indexar los tipos de voto
func (sesion *sesionVotar) indiceTipo(tipo string) int{
	for i,valor := range sesion.identificadores_tipos{
		if(valor == tipo){
			return i
		}
	}
	return -1
}


// funciones de la interfaz

func (sesion *sesionVotar) HayVotante() bool{
	return !sesion.esperandoAVotar.EstaVacia()
}

func (sesion *sesionVotar) IngresarVotante(dniStr string) error{
	dni,err := strconv.Atoi(dniStr)
		
	if(err != nil || dni <0){
		return new(ErrorDNIInvalido)
	}

	votante := buscarVotante(sesion.aVotar,dni)

	if(votante == nil){
		return new(ErrorDNINoEsta)
	} 
	sesion.esperandoAVotar.Encolar(votante)
	return nil
}

func (sesion *sesionVotar) Votar(tipoStr string, candidatoStr string) error{

	if(sesion.esperandoAVotar.EstaVacia()){
		return new(ErrorFilaVacia)
	}
	// Se prefiere a hacer los chequeos en sesion votar
	// para ahorrar memoria
	tipo:= sesion.indiceTipo(tipoStr)

	if tipo == -1{
		return new(ErrorTipoInvalido)
	}

	candidato,err := strconv.Atoi(candidatoStr)
	
	if(err != nil || candidato<0 || candidato >= len(sesion.listaCandidatos)){
		return new(ErrorAlternativaInvalida)
	}

	if sesion.esperandoAVotar.VerPrimero().YaVoto() {
		return CrearErrorFraude(sesion.esperandoAVotar.Desencolar().DNI())
	}

	// Cambio de voto
	sesion.registroDeVotos.Agregar(sesion.esperandoAVotar.VerPrimero().CambiameElVoto(tipo,candidato))

	return nil
}

func (sesion *sesionVotar) Deshacer() error{
	if(sesion.esperandoAVotar.EstaVacia()){
		return new(ErrorFilaVacia)
	}

	if sesion.esperandoAVotar.VerPrimero().YaVoto(){
		return CrearErrorFraude(sesion.esperandoAVotar.Desencolar().DNI())
	}

	return sesion.registroDeVotos.Borrar()
}

func (sesion *sesionVotar) SiguienteVotante() error {

	if(sesion.esperandoAVotar.EstaVacia()){
		return new(ErrorFilaVacia)
	}

	sesion.registroDeVotos.Vaciar()
	votante:= sesion.esperandoAVotar.Desencolar()
	err := votante.FinalizarVoto()

	if(err == nil){
		if(votante.Impugnado()){
			sesion.votosImpugnados++
		} else{
			votante.MirarVotos(func(tipo int, candidato int) {
				sesion.listaCandidatos[candidato][tipo].votantes++
			})	
		}

		

	}

	return err
}

func (sesion *sesionVotar) Finalizar() error{
	var err error = nil
	
	if !sesion.esperandoAVotar.EstaVacia() {

		// para evitar se use a futuro, no es necesario si se asume se van solos los votantes.
		for !sesion.esperandoAVotar.EstaVacia(){ 
			sesion.esperandoAVotar.Desapilar()
		}

		err = new(ErrorSinTerminar)
	}

	// para evitar su uso a futuro una vez se finalizo
	sesion.aVotar = make([]Votante,0) 

	return err
}


// Funciones para tests/informacion

func (sesion *sesionVotar) VotosImpugnados() int{
	return sesion.votosImpugnados
}

func (sesion *sesionVotar) IterarVotos(identificador string,visitar func(string,int) bool){
		tipo := sesion.indiceTipo(identificador)
		if(tipo == -1){
			panic(ERROR_TIPO_INVALIDO)
		}

		if (!visitar("Votos en blanco",sesion.listaCandidatos[0][tipo].votantes)){
			return
		}
		i:= 1
		for i<len(sesion.listaCandidatos) && visitar(sesion.listaCandidatos[i][tipo].Credencial(),
										sesion.listaCandidatos[i][tipo].votantes){
			i++
		}
}