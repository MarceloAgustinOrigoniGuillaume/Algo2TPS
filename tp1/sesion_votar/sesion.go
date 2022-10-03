package sesion_votar
//Constante para testeo 
const BASIC_SAMPLE = "<BasicSample>" // una url no permite normalmente almenos el < ni el >, constante para tests

// Interfaz 

type SesionVotar interface{

	// Esta podria bien estar sola con Finalizar pero se prefiere permitir usar las primitivas directamente
	IngresarVotante(dniStr string) error

	SiguienteVotante() error

	HayVotante() bool

	Deshacer() error
	
	Votar(tipoStr string, candidatoStr string) error
	
	Finalizar() error
	

	// Funciones con proposito de testeo
	VotosImpugnados() int
	IterarVotos(identificador string,visitar func(string,int) bool)
}