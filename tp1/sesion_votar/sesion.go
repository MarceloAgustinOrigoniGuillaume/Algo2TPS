package sesion_votar

// TO DO , CAMBIOS A HACER // Haces los que puedas/quieras

// MUST DO, Obligatorios

// En manejo_datos.go, hacer correctamente las funciones popular desde archivos. Hacen lo que hacen.
// idealmente usar las funciones ya creadas de CrearArregloDesdeArchivo y LeerArchivo.
// En manejo_datos.go mejorar la funcion Contiene usada para el dni (lineal) que se tiene, asumiendo o bien dicho habiendo ordenado los padrones en popular

// Hacer que main.go reciba input del usuario por consola.. capaz ver de reutilizar codigo al hacerlo

// Seguramente tmbn obligatorio, cambiar el sistema de errores a usar errores como tal, no es muy dificil la vd solo cambiar los return type y asi.
// DUDOSO TEMA

// A la hora de ver si ya voto alguien se usa una lista enlazada, si fuera dinamica vaya y pase pero por mas 
// se ordene por insercion que es sencillo en una lista enlazada no se puede aprovechar la busqueda binaria....
// Como mejorar eso??? Y hacerlo... funciones envueltas = yaVoto y contarVotante de sesion_votar.go
// Propuesta es usar algun tipo de lista dinamica. O un arreglo a redimnesionar.. ni en hago algo con una lista enlazada

// Menos urgentes

// Comentar lo que hacen las funciones en la interfaz.... es bastante directo ja y los posibles estados de salida

// Agregar mas tests como estan en archivos/set1 y agregarlos a la funcion de test de archivos( el ultimo test)

// Agregar los tests de la catedra, y hacer la funcion correspondiente en sesion_votar_tests

// Cambiar a usar errores, habria que hacer un cambio de tipos y poner los structs en la carpeta de errores, bla bla




//Constantes 

const BASIC_SAMPLE = "<BasicSample>" // una url no permite normalmente almenos el < ni el >, constante para tests
const OK = "OK"

const ERROR_COMANDO_INVALIDO = "Comando invalido"
const ERROR_FALTAN_PARAMETROS = "ERROR: Faltan parametros"
const ERROR_DNI_INVALIDO = "ERROR: DNI incorrecto"
const ERROR_DNI_NO_ESTA = "ERROR: DNI fuera del padron"
const ERROR_FILA_VACIA = "ERROR: Fila vacia"
const ERROR_TIPO_INVALIDO = "ERROR: Tipo de voto invalido"
const ERROR_ALTERNATIVA_INVALIDA = "ERROR: Alternativa invalida"
const ERROR_VOTANTE_FRAUDULENTO = "ERROR: Votante FRAUDULENTO: %d"
const ERROR_SIN_VOTO_DESHACER = "ERROR: Sin voto a deshacer"
const ERROR_SIN_TERMINAR = "ERROR: Ciudadanos sin terminar de votar"


// Interfaz 

type SesionVotar interface{

	// Esta podria bien estar sola con Finalizar pero se prefiere permitir usar las primitivas directamente
	AccionDesdeComando(comando string) string 

	IngresarVotante(dni int) string

	SiguienteVotante() string

	HayVotante() bool

	Deshacer() string
	
	Votar(tipo int, candidato int) string
	
	Finalizar() string
	

	// Funciones con proposito de testeo

	VotosImpugnados() int
	IterarVotos(identificador string,visitar func(string,int) bool)
}