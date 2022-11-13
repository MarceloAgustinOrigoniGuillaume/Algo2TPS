package usuario

// El usuario algogram extiende el DatosUsuario agregando su indice, de archivo, o insercion si se quiere.
type UsuarioAlgogram interface{
	DatosUsuario

	Indice() int
}


type usuarioAlgogram struct {
	nombre string
	indice int
}

func CrearUsuarioAlgogram(nombre string, indice int) UsuarioAlgogram {
	return &usuarioAlgogram{nombre,indice}
}

func (user *usuarioAlgogram) Nombre() string {
	return user.nombre
}

func (user *usuarioAlgogram) Indice() int {
	return indice
}
