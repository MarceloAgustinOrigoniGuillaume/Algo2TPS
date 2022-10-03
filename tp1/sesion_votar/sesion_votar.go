package sesion_votar

import (
	TDACola "cola"
	"strconv"
)


type sesionVotar struct{
	identificadores_tipos []string
	listaCandidatos [][]candidatoStruct
	aVotar []Votante
	esperandoAVotar TDACola.Cola[Votante]
	registroDeVotos Registro
	votosImpugnados int // al parecer votos impugnados significaba votos en blanco?
}

//Funcion crear

func CrearSesion(identificadores_tipos []string,candidatos_file string,padrones_file string) (SesionVotar,error){
	sesion := new(sesionVotar)


	sesion.esperandoAVotar = TDACola.CrearColaEnlazada[Votante]()
	sesion.registroDeVotos = CrearRegistroDeVotos()
	// Cargar candidatos y padrones/dnis

	sesion.identificadores_tipos = identificadores_tipos
	var err error = nil
	if(candidatos_file == BASIC_SAMPLE){
		sesion.listaCandidatos = popularCandidatosBasico(len(sesion.identificadores_tipos)+1)
	} else{
		sesion.listaCandidatos,err = popularCandidatos(candidatos_file,len(sesion.identificadores_tipos)+1)
	}

	if(padrones_file == BASIC_SAMPLE){
		sesion.aVotar = popularVotantesBasico(len(identificadores_tipos))
	} else{
		sesion.aVotar,err = popularVotantes(padrones_file,len(identificadores_tipos))
	}

	if(err != nil || len(sesion.listaCandidatos) < len(sesion.identificadores_tipos)+1){
		err = new(ErrorLecturaArchivos)
	}
	


	return sesion,err
}

// funciones auxiliares

// deberia hacerse busqueda binaria   
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


func (sesion *sesionVotar) indiceTipo(tipo string) int{
	for i,valor := range sesion.identificadores_tipos{
		if(valor == tipo){
			return i
		}
	}
	return -1
}


// interfaz

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

	// Verificaciones
	if(sesion.esperandoAVotar.EstaVacia()){
		return new(ErrorFilaVacia)
	}
	// Se prefiere a hacer los chequeos en sesion votar
	// ya que no se quiere guardar en votante la cantidad de candidatos ni tampoco hacer un switch con errores
	tipo:= sesion.indiceTipo(tipoStr)

	if tipo == -1{
		return new(ErrorTipoInvalido)
	}


	candidato,err := strconv.Atoi(candidatoStr)
	if(err != nil || candidato<0 || candidato >= len(sesion.listaCandidatos[tipo])){
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
				sesion.listaCandidatos[tipo][candidato].votantes++
			})	
		}

		

	}

	return err
}

func (sesion *sesionVotar) Finalizar() error{
	var err error = nil
	
	if !sesion.esperandoAVotar.EstaVacia() {
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

		candidatos:= sesion.listaCandidatos[tipo]

		if (!visitar("Votos en blanco",candidatos[0].votantes)){
			return
		}

		i:= 1

		for i<len(candidatos) && visitar(candidatos[i].Credencial(),candidatos[i].votantes){
			i++
		}
}