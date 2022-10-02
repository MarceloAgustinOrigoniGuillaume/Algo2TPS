package sesion_votar

import (
	TDACola "cola"
	TDAPila "pila"
	TDALista "lista"
	"fmt"
	"strings"
	"strconv"
)
type sesionVotar struct{
	aVotar TDACola.Cola[*votanteStruct]
	yaVotaron TDALista.Lista[int]
	listaCandidatos [][]candidatoStruct
	padronesValidos []int
	identificadores_tipos []string
	accionesVotante TDAPila.Pila[func()] // se podria hacer desde el Votante.. pero no quiero desperdiciar memoria
	votosImpugnados int
}

//Funcion crear

func CrearSesion(identificadores_tipos []string,candidatos_file string,padrones_file string) SesionVotar{
	sesion := new(sesionVotar)


	sesion.aVotar = TDACola.CrearColaEnlazada[*votanteStruct]()
	sesion.accionesVotante = TDAPila.CrearPilaDinamica[func()]()
	sesion.yaVotaron = TDALista.CrearListaEnlazada[int]()


	// Cargar candidatos y padrones/dnis

	sesion.identificadores_tipos = identificadores_tipos
	if(candidatos_file == BASIC_SAMPLE){
		sesion.listaCandidatos = popularCandidatosBasico()
	} else{
		sesion.listaCandidatos = popularCandidatos(candidatos_file)
	}

	if(candidatos_file == BASIC_SAMPLE){
		sesion.padronesValidos = popularDNISBasico()
	} else{
		sesion.padronesValidos = popularDNIS(padrones_file)
	}
	

	return sesion
}


// funciones para hacer con comandos strings.. podrian estar aparte? capaz

func (sesion *sesionVotar) indiceTipo(tipo string) int{
	for i,valor := range sesion.identificadores_tipos{
		if(valor == tipo){
			return i
		}
	}
	return -1
}

func (sesion *sesionVotar) AccionDesdeComando(comando string) string{
	// Me encantaria usar un hash de strings... ie comandos a funciones, igual sirven dos arreglos y ya
	// Se que se podria usar simplemente ifs, o un switch, pero ni ganas de hacerlos y ademas es mas facil de escalar de esta forma
	
	if(comando == "fin-votar"){
		return sesion.SiguienteVotante()
	} else if(comando == "deshacer"){
		return sesion.Deshacer()
	}



	args := strings.Split(comando," ")

	if(args[0] == "ingresar"){

		if(len(args) < 2){ // yo pondria != pero el error es solo en falta
			return ERROR_FALTAN_PARAMETROS
		}
		// validame el dni

		dni,error := strconv.Atoi(args[1])
		if(error != nil || dni <0){
			return ERROR_DNI_INVALIDO
		}

		return  sesion.IngresarVotante(dni)

	}

	if(args[0] == "votar"){

		if(len(args) < 3){
			return ERROR_FALTAN_PARAMETROS
		}
		tipo := sesion.indiceTipo(args[1])

		if tipo == -1{
			return ERROR_TIPO_INVALIDO
		}

		// validame el dni

		candidato,error := strconv.Atoi(args[2])
		if(error != nil){
			return ERROR_ALTERNATIVA_INVALIDA
		}

		// validame el tipo voto y candidato


		return sesion.Votar(tipo,candidato)
	}
	
	return ERROR_COMANDO_INVALIDO

}



//funciones de utilidad interna

func (sesion *sesionVotar) yaVoto(dni int) bool{
	// ESTA BUSQUEDA ES LINEAL , VALE LA PENA ORDENAR? LA BUSQUEDA BINARIA NO CONVENZE CON UNA LISTA ENLAZADA
	// SE TIENE QUE REVISAR....
	iterador := sesion.yaVotaron.Iterador() 

	for iterador.HaySiguiente() && iterador.VerActual() <= dni {
		if(iterador.Siguiente() == dni){ // Siguiente retorna el actual, por eso usamos el for con <=
			return true
		}
	}

	return false
}

func (sesion *sesionVotar) contarVotante(dni int) bool{
	// CAPAZ HAY UNA MEJOR FORMA
	// Se agrega el votante a los que votaron si no se agrego ya. Retorna si se agrego. Si no lo hizo era fraude
	// Se ingresara mediante insercion, al irse agregando de a uno es lo mas eficiente.
	// Y al mismo tiempo se verifica que no este.

	iterador := sesion.yaVotaron.Iterador()

	for iterador.HaySiguiente() && iterador.VerActual() <= dni {
		if(iterador.Siguiente() == dni){ // Siguiente retorna el actual, por eso usamos el for con <=
			return false
		}
	}
	iterador.Insertar(dni)

	return true


}

func (sesion *sesionVotar) limpiarAcciones(){
	for !sesion.accionesVotante.EstaVacia(){
		sesion.accionesVotante.Desapilar()
	}
}

// implementacion


func (sesion *sesionVotar) HayVotante() bool{
	return !sesion.aVotar.EstaVacia()
}


func (sesion *sesionVotar) IngresarVotante(dni int) string{
	if(!Contiene(sesion.padronesValidos,dni)){
		return ERROR_DNI_NO_ESTA
	}
	// Encolar, requiere Cola, struct votante.
	sesion.aVotar.Encolar(crearVotante(dni,len(sesion.listaCandidatos)))
	return OK
}

func (sesion *sesionVotar) Votar(tipo int, candidato int) string{
	if !sesion.HayVotante(){
		return ERROR_FILA_VACIA
	}

	if(sesion.yaVoto(sesion.aVotar.VerPrimero().dni)){
		sesion.votosImpugnados++
		return fmt.Sprintf(ERROR_VOTANTE_FRAUDULENTO,sesion.aVotar.Desencolar().dni)
	}


	if (tipo<0 ||tipo >= len(sesion.listaCandidatos)){
		return ERROR_TIPO_INVALIDO
	}

	candidatos := sesion.listaCandidatos[tipo]
	if (candidato<0 || candidato >= len(candidatos)){
		return ERROR_ALTERNATIVA_INVALIDA
	}


	votante := sesion.aVotar.VerPrimero()
	valor_actual := votante.votos[tipo]
	sesion.accionesVotante.Apilar(func() {votante.votos[tipo] = valor_actual} )
	votante.votos[tipo] = candidato


	return OK
	//return nil
}


func (sesion *sesionVotar) Deshacer() string{
	if(sesion.aVotar.EstaVacia()){
		return ERROR_FILA_VACIA
	}

	if(sesion.yaVoto(sesion.aVotar.VerPrimero().dni)){
		sesion.votosImpugnados++
		return fmt.Sprintf(ERROR_VOTANTE_FRAUDULENTO,sesion.aVotar.Desencolar().dni)
	}

	if(sesion.accionesVotante.EstaVacia()){
		return ERROR_SIN_VOTO_DESHACER
	}

	
	sesion.accionesVotante.Desapilar()()
	return OK
}

func (sesion *sesionVotar) SiguienteVotante() string {
	// Desencola , requiere Cola, struct votante.
	// Si no hay votante bla bla.

	if(sesion.aVotar.EstaVacia()){
		return ERROR_FILA_VACIA
	}

	sesion.limpiarAcciones()

	votante := sesion.aVotar.Desencolar()
	if(!sesion.contarVotante(votante.dni)){
		sesion.votosImpugnados++
		return fmt.Sprintf(ERROR_VOTANTE_FRAUDULENTO,votante.dni)
	}


	for tipo,candidato := range votante.votos{
		sesion.listaCandidatos[tipo][candidato].votantes++
	}

	//sesion.agregarVoto(votante) // se quiso usar una funcion auxiliar pero para un for, que solo se usa aca no lo vale
	//return nil
	return OK
}

func (sesion *sesionVotar) Finalizar() string{
	
	// Aca se muestra por defecto el resultado esto capaz se decide cambiarlo
	//MostrarEstado(sesion.identificadores_tipos,sesion)


	// Tambien se deberia mandar a dormir a la sesion, se supone se finalizo
	// Es decir pasar a modo solo lectura. Pero eso no se va a modificar porque si bien puede ser logico
	// tampoco es completamente necesario
	
	if !sesion.aVotar.EstaVacia() {
		return ERROR_SIN_TERMINAR
	}

	return OK
}


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