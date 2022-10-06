package sesion_votar

// Se pone la interfaz en el mismo archivo por que idealmente solo existira una implementacion
// Pero no se quiere permitir usar la struct directamente
type Votante interface {
	// indica si ya voto
	YaVoto() bool

	// indica si esta impugnado
	Impugnado() bool

	// devuelve el dni
	DNI() int

	// cambia el voto y devuelve la accion para deshacer este cambio
	CambiameElVoto(tipo int, candidato int) func()

	// finaliza el voto, o devuelve un error si ya voto
	FinalizarVoto() error

	// iterador para fijarse los votos
	MirarVotos(observar func(int, int))
}

type votanteStruct struct {
	dni       int
	votos     []int
	yaVoto    bool
	impugnado bool
}

func (votante *votanteStruct) DNI() int {
	return votante.dni
}

func CrearVotante(dni int, cant_tipo_voto int) Votante {
	persona := new(votanteStruct)
	persona.dni = dni
	persona.votos = make([]int, cant_tipo_voto, cant_tipo_voto)

	return persona
}

func (votante *votanteStruct) YaVoto() bool {
	return votante.yaVoto
}

func (votante *votanteStruct) Impugnado() bool {
	return votante.impugnado
}

func (votante *votanteStruct) CambiameElVoto(tipo int, candidato int) func() {
	if votante.impugnado {
		// retornaria nil pero su sistema al parecer requiere que aunque no se tomasen en cuenta
		// se deshagan
		return func() {}
	}

	anterior := votante.votos[tipo]
	votante.votos[tipo] = candidato

	if candidato == 0 {
		votante.impugnado = true
		return func() {
			votante.votos[tipo] = anterior
			votante.impugnado = false
		}
	}
	return func() { votante.votos[tipo] = anterior }

}

func (votante *votanteStruct) FinalizarVoto() error {
	if votante.yaVoto {
		return CrearErrorFraude(votante.dni)
	}
	votante.yaVoto = true

	return nil
}

func (votante *votanteStruct) MirarVotos(observar func(int, int)) {
	for tipo, candidato := range votante.votos {
		observar(tipo, candidato)
	}
}
