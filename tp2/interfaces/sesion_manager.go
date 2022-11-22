package interfaces

type SesionManager interface {
	Login(string) error
	Logout() error
	Publicar(string) error
	VerSiguientePost() (string, error)
	Likear(string) error
	MostrarLikes(string) (string, error)
}
