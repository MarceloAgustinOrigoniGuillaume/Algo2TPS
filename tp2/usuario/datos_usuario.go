package usuario

type DatosUsuario interface {
	Nombre() string
}

type datosUsuario struct {
	nombre string
}

func CrearDatosUsuario(nombre string) DatosUsuario {
	return &datosUsuario{nombre}
}

func (datos *datosUsuario) Nombre() string {
	return datos.nombre
}
