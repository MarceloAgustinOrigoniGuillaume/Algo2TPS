package sesion_votar

import "fmt"
// Constantes
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


// errores main

type ErrorComandoInvalido struct{}

func (err *ErrorComandoInvalido) Error() string{
	return ERROR_COMANDO_INVALIDO
}

type ErrorFaltanParametros struct{}

func (err *ErrorFaltanParametros) Error() string{
	return ERROR_FALTAN_PARAMETROS
}


// errores manejo archivos
type ErrorOmicion struct{}

func (err *ErrorOmicion) Error() string{
	return "Warning: se ignoro elemento invalido"
}

type ErrorMissMatchSizeOut struct{}

func (err *ErrorMissMatchSizeOut) Error() string{
	return "ERROR: Habia mas lineas en el archivo out que en el in. Se ignoraron las sobrantes"
}


type ErrorMissMatchSizeIn struct{}

func (err *ErrorMissMatchSizeIn) Error() string{
	return "ERROR: Habia mas lineas en el archivo in que en el out. Se ignoraron las sobrantes"
}


// errores sesion votar como tal

type ErrorDNIInvalido struct{}

func (err *ErrorDNIInvalido) Error() string{
	return ERROR_DNI_INVALIDO
}

type ErrorDNINoEsta struct{}

func (err *ErrorDNINoEsta) Error() string{
	return ERROR_DNI_NO_ESTA
}


type ErrorFraude struct{
	dniVotante int
}

func (err *ErrorFraude) Error() string{
	return fmt.Sprintf(ERROR_VOTANTE_FRAUDULENTO,err.dniVotante)
}

func CrearErrorFraude(dni int) error{
	err:= new(ErrorFraude)
	err.dniVotante = dni
	return err
}

type ErrorTipoInvalido struct{}

func (err *ErrorTipoInvalido) Error() string{
	return ERROR_TIPO_INVALIDO
}

type ErrorAlternativaInvalida struct{}

func (err *ErrorAlternativaInvalida) Error() string{
	return ERROR_ALTERNATIVA_INVALIDA
}



type ErrorSinTerminar struct{}

func (err *ErrorSinTerminar) Error() string{
	return ERROR_SIN_TERMINAR
}

type ErrorFilaVacia struct{}

func (err *ErrorFilaVacia) Error() string{
	return ERROR_FILA_VACIA
}


type ErrorSinRegistro struct{}

func (err *ErrorSinRegistro) Error() string{
	return ERROR_SIN_VOTO_DESHACER
}