package sesion

//import "tp2/interfaces"

type usuarioAlgogram struct {
	nombre string
	indice int
}

func crearUsuarioAlgogram(nombre string, indice int) *usuarioAlgogram {
	return &usuarioAlgogram{nombre, indice}
}

func (user *usuarioAlgogram) Nombre() string {
	return user.nombre
}

func (user *usuarioAlgogram) Indice() int {
	return user.indice
}
