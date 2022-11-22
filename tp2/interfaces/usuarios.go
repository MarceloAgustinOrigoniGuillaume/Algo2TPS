package interfaces

type DatosUsuario interface {
	Nombre() string
}

// El usuario algogram extiende el DatosUsuario agregando su indice, de archivo, o insercion si se quiere.
type UsuarioAlgogram interface{
	DatosUsuario
	Indice() int
}