package sesion_votar
type candidatoStruct struct{
	partido string
	nombre string
	votantes int
}

func CrearCandidato(partido string, nombre string) candidatoStruct{
	return candidatoStruct{partido,nombre,0}
}

func (candidato *candidatoStruct) Credencial() string{
	return candidato.partido+" - "+candidato.nombre
}