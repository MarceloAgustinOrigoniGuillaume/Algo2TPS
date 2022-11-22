package managers

import "tp2/interfaces"
import "fmt"

// Funciones de comparacion
// Capaz se podrian dejar como parametro? .... No quise complicar demasiado el tema y mantener la abstraccion puede ser bueno
func comparacionAfinidad(usuario1, opcion1, opcion2 interfaces.UsuarioAlgogram) int {
	indice_target := usuario1.Indice()

	diff1 := indice_target - opcion1.Indice()
	if diff1 < 0 {
		diff1 = -diff1
	}

	diff2 := indice_target - opcion2.Indice()

	if diff2 < 0 {
		diff2 = -diff2
	}

	return diff2 - diff1

}

func comparacionPosts[D interfaces.UsuarioAlgogram, P interfaces.Post[int, D]](post1, post2 P) int {
	return post2.Id() - post1.Id()
}

// Los tipos D y P estan mas que nada para que Barbara se ahorre castear los tipos
// Estos tienen que ser Tipos de struct en especifico, ya que si son interfaces habria problema seguramente con primitivas
type recomendadorAlgogram[D interfaces.UsuarioAlgogram, P interfaces.Post[int, D]] struct {
	feeds interfaces.IdManager[int, interfaces.MapConexiones[D, P]] //ColaPrioridad[*postYAutor[P]]
}

func (recomendador *recomendadorAlgogram[D, P]) crearMapConexiones(usuario D) interfaces.MapConexiones[D, P] {
	return crearConexionUserPost[D, P](usuario,
		func(op1 P, op2 P) int {
			res := comparacionAfinidad(usuario, op1.Autor(), op2.Autor())

			if res == 0 {
				res = comparacionPosts[D, P](op1, op2)
			}

			return res
		})
}

func CrearEmptyRecomendadorAlgogram[D interfaces.UsuarioAlgogram, P interfaces.Post[int, D]]() interfaces.Recomendador[int, D, P] {
	return &recomendadorAlgogram[D, P]{CrearNumericalIdManager[interfaces.MapConexiones[D, P]]()}
}

// Al no necesitar agregar nuevos usuarios esto seria mas rapido, pero para mayor simpleza de codigo se usara la primitiva
func CrearRecomendadorAlgogram[D interfaces.UsuarioAlgogram, P interfaces.Post[int, D]](cantidadUSuarios int, userProvider func(int) D) interfaces.Recomendador[int, D, P] {

	recomendador := new(recomendadorAlgogram[D, P])

	recomendador.feeds = initializedNumericalIdManager[interfaces.MapConexiones[D, P]](cantidadUSuarios,
		func(indice int) interfaces.MapConexiones[D, P] {
			return recomendador.crearMapConexiones(userProvider(indice))
		})

	return recomendador
}

func (recomendador *recomendadorAlgogram[D, P]) AgregarUsuario(usuario D) {
	recomendador.feeds.Agregar(recomendador.crearMapConexiones(usuario))
}

// Agrega un post a todos los feeds.
func (recomendador *recomendadorAlgogram[D, P]) AgregarPost(autor D, post P) {
	indiceAutor := autor.Indice()

	recomendador.feeds.Iterar(func(indiceUsuario int, mapa interfaces.MapConexiones[D, P]) bool {
		if mapa == nil {
			fmt.Printf("QUE CHOTAS?\n")
			return true
		}

		if indiceAutor != indiceUsuario { // No parece muy lindo chequear asi pero no veo otra forma
			mapa.AgregarConexion(post)
		}

		return true
	})

}

// Obtiene las recomendaciones para un usuario
func (recomendador *recomendadorAlgogram[D, P]) ObtenerRecomendaciones(usuario D) interfaces.MapConexiones[D, P] {
	return recomendador.feeds.Obtener(usuario.Indice())
}
