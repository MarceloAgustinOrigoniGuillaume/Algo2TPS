package sesion_votar

import TDAPila "pila"



type Registro interface{
	Borrar() error
	Agregar(func())
	Vaciar()
}

type registroVotos struct{
	votosRegistrados TDAPila.Pila[func()]
}


func CrearRegistroDeVotos() Registro{
	registro := new(registroVotos)
	registro.votosRegistrados = TDAPila.CrearPilaDinamica[func()]()
	return registro
}


func (registro *registroVotos) Borrar() error {

	if(registro.votosRegistrados.EstaVacia()){
		return new(ErrorSinRegistro)
	}

	registro.votosRegistrados.Desapilar()()
	return nil
}


func (registro *registroVotos) Agregar(nuevo func()){
	if(nuevo != nil){
		registro.votosRegistrados.Apilar(nuevo)
	}
}

func (registro *registroVotos) Vaciar(){
	for !registro.votosRegistrados.EstaVacia(){
		registro.votosRegistrados.Desapilar()
	}
}