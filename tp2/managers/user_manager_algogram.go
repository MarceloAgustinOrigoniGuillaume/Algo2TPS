package managers

import "tp2/interfaces"
import hash "hash/hashCerrado"

// Un id Manager con Nombre como id, y usuarios como dato
type userManagerAlgogram[D interfaces.UsuarioAlgogram] struct {
	usuarios hash.Diccionario[string, D]
}

func CrearUserManagerAlgogram[D interfaces.UsuarioAlgogram]() interfaces.IdManager[string, D] {
	return &userManagerAlgogram[D]{hash.CrearHash[string, D]()}
}

func (manager *userManagerAlgogram[D]) Agregar(usuario D) string {
	manager.usuarios.Guardar(usuario.Nombre(), usuario)
	return usuario.Nombre()
}

func (manager *userManagerAlgogram[D]) Existe(nombre string) bool {
	return manager.usuarios.Pertenece(nombre)
}

func (manager *userManagerAlgogram[D]) Obtener(nombre string) D {
	return manager.usuarios.Obtener(nombre)
}

// Undefined, el id usado por manager algogram es obtenido de los datos del usuario
func (manager *userManagerAlgogram[D]) NuevoId() string {
	return ""
}

func (manager *userManagerAlgogram[D]) Iterar(visitar func(string, D) bool) {
	manager.usuarios.Iterar(visitar)
}
