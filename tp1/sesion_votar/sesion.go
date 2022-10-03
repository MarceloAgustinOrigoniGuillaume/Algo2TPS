package sesion_votar

// TO DO , CAMBIOS A HACER // Haces los que puedas/quieras

// MUST DO, Obligatorios

// En manejo_datos.go, hacer correctamente las funciones popular desde archivos. Hacen lo que hacen.
// idealmente usar las funciones ya creadas de CrearArregloDesdeArchivo y LeerArchivo.
// En sesion_votar.go mejorar la funcion buscar usada para el dni (lineal) que se tiene, asumiendo o bien dicho habiendo ordenado por padron al popular

// Hacer que main.go reciba input del usuario por consola.. capaz ver de reutilizar codigo al hacerlo

// Menos urgentes

// Comentar lo que hacen las funciones en la interfaz.... es bastante directo ja y los posibles estados de salida

// Agregar mas tests como estan en archivos/set1 y agregarlos a la funcion de test de archivos( el ultimo test)

// Agregar los tests de la catedra, y hacer la funcion correspondiente en sesion_votar_tests

// Problemas de refactorizacion.....
// Corregir errores en valores expected en sesion_votar_test.go, si no se voto para todo cargo es voto en blanco
// no hay voto en blanco parcial al parecer... igual mejor verificar con los tests de la catedra.


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