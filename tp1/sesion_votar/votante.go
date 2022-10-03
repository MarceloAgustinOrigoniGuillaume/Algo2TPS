package sesion_votar

type Votante interface{
	YaVoto() bool
	Impugnado() bool
	CambiameElVoto(tipo int, candidato int) func()
	FinalizarVoto() error
	MirarVotos(observar func(int,int))
	DNI() int
}

type votanteStruct struct{
	dni int
	votos []int
	yaVoto bool
	impugnado bool
}

func (votante *votanteStruct) DNI() int{
	return votante.dni
}

func CrearVotante(dni int,cant_tipo_voto int) Votante{
	persona := new(votanteStruct)
	persona.dni = dni
	persona.votos = make([]int,cant_tipo_voto,cant_tipo_voto)

	return persona
}

func (votante *votanteStruct) YaVoto() bool{
	return votante.yaVoto
}

func (votante *votanteStruct) Impugnado() bool{
	return votante.impugnado
}


func (votante *votanteStruct) CambiameElVoto(tipo int, candidato int) func(){
	if(votante.impugnado){
		return nil
	}



	anterior:= votante.votos[tipo]
	votante.votos[tipo] = candidato
	
	if(candidato == 0){
		votante.impugnado = true
		return func(){ 
			votante.votos[tipo] = anterior
			votante.impugnado = false
		}
	}
	return func(){ votante.votos[tipo] = anterior}

}


func (votante *votanteStruct) FinalizarVoto() error{
	if(votante.yaVoto){
		return CrearErrorFraude(votante.dni)
	}
	votante.yaVoto = true
	for _,voto := range votante.votos{
		if(voto == 0){
			votante.votos = make([]int,len(votante.votos))
			return nil
		}
	}

	return nil
}

func (votante *votanteStruct) MirarVotos(observar func(int,int)){
	for tipo,candidato := range votante.votos{
		observar(tipo,candidato)
	}
}