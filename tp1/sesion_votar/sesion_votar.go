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
	votosImpugnados int
}

//Funcion crear

func CrearSesion(identificadores_tipos []string,candidatos_file string,padrones_file string) SesionVotar{
	sesion := new(sesionVotar)


	sesion.esperandoAVotar = TDACola.CrearColaEnlazada[Votante]()
	sesion.registroDeVotos = CrearRegistroDeVotos()
	// Cargar candidatos y padrones/dnis

	sesion.identificadores_tipos = identificadores_tipos
	
	if(candidatos_file == BASIC_SAMPLE){
		sesion.listaCandidatos = popularCandidatosBasico()
	} else{
		sesion.listaCandidatos = popularCandidatos(candidatos_file)
	}

	if(candidatos_file == BASIC_SAMPLE){
		sesion.aVotar = popularVotantesBasico(len(identificadores_tipos))
	} else{
		sesion.aVotar = popularVotantes(padrones_file,len(identificadores_tipos))
	}
	

	return sesion
}

// funciones auxiliares

// deberia hacerse busqueda binaria   
func buscarVotante(lista_dni []Votante,dni int) Votante { 
	for _,votante := range lista_dni{
		if(votante.DNI() == dni){
			return votante
		}
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
		sesion.votosImpugnados++
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
		sesion.votosImpugnados++
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

	if(err != nil){ // unico error posible es fraude
		sesion.votosImpugnados++
	} else{
		votante.MirarVotos(func(tipo int, candidato int) {
			sesion.listaCandidatos[tipo][candidato].votantes++
		})

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