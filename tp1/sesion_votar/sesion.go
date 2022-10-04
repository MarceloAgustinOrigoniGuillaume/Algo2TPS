package sesion_votar

//Constante para testeo 
const BASIC_SAMPLE = ";Sample;" // una url no permite normalmente el ;, constante para tests

// Interfaz 

type SesionVotar interface{

	// Ingresa un votante que deberia votar a futuro, en el orden de ingreso. Devuelve el error
	// correspondiente si el dni es invalido o no esta en el padron
	IngresarVotante(dniStr string) error

	// Avanza al siguiente votante, finalizando el voto del actual.
	// Puede devolver un error de fraude o no hay votante segun corresponda
	SiguienteVotante() error	

	//Devuelve si hay un votante
	HayVotante() bool

	// Deshace la ultima accion del votante actual, puede devolver error por fraude
	// error no hay votante o error sin que deshacer.
	Deshacer() error
	
	// Se registra un cambio de voto en el votante actual
	// Se utiliza el 0 para impugnar el voto, lo cual genera que se ignore futuras llamadas
	// A menos se deshaga. 
	Votar(tipoStr string, candidatoStr string) error
	
	// Si finaliza la sesion de votacion, devuelve un error si no todos los que vinieron a votar
	// lo hicieron, y cierra el sistema
	Finalizar() error
	
	// Funciones para ver resultados
	VotosImpugnados() int
	IterarVotos(identificador string,visitar func(string,int) bool)
}