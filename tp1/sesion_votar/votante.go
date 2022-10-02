package sesion_votar

type votanteStruct struct{
	dni int
	votos []int
}

func crearVotante(dni int,cant_tipo_voto int) *votanteStruct{
	persona := new(votanteStruct)
	persona.dni = dni
	persona.votos = make([]int,cant_tipo_voto,cant_tipo_voto)

	return persona
}