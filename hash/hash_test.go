package diccionario_test

import (
	Hash "diccionario"
	"github.com/stretchr/testify/require"
	"testing"
)



func TestAgregado(t *testing.T){
	hashCuckoo := Hash.CrearHash[string,int]()

	require.NotNil(t,hashCuckoo,"Crear devolvio nil")


	hashCuckoo.Guardar("Uno",1)
	hashCuckoo.Guardar("Dos",2)
	hashCuckoo.Guardar("Tres",3)

	require.EqualValues(t,3,hashCuckoo.Cantidad(),"Cantidad incorrecta")
}


func TestPertenece(t *testing.T){
	hashCuckoo := Hash.CrearHash[string,int]()

	require.NotNil(t,hashCuckoo)


	hashCuckoo.Guardar("Uno",1)
	hashCuckoo.Guardar("Dos",2)
	hashCuckoo.Guardar("Tres",3)

	require.EqualValues(t,3,hashCuckoo.Cantidad(),"Cantidad incorrecta")

	require.False(t,hashCuckoo.Pertenece("O"),"Se invento una clave 'O'")
	require.False(t,hashCuckoo.Pertenece(""),"Se invento una clave ''")
	require.False(t,hashCuckoo.Pertenece("Makise"),"Se invento una clave 'Makise'")

	require.True(t,hashCuckoo.Pertenece("Uno"),"Se olvido de la clave 'Uno'")
	require.True(t,hashCuckoo.Pertenece("Dos"),"Se olvido de la clave 'Dos'")
	require.True(t,hashCuckoo.Pertenece("Tres"),"Se olvido de la clave 'Tres'")
}



func TestObtener(t *testing.T){
	hashCuckoo := Hash.CrearHash[string,int]()

	require.NotNil(t,hashCuckoo)


	hashCuckoo.Guardar("Uno",1)
	hashCuckoo.Guardar("Dos",2)
	hashCuckoo.Guardar("Tres",3)

	require.EqualValues(t,3,hashCuckoo.Cantidad(),"Cantidad incorrecta")

	require.PanicsWithValue(t,Hash.ERROR_NO_ESTABA,func() { hashCuckoo.Obtener("Makise")})
	require.PanicsWithValue(t,Hash.ERROR_NO_ESTABA,func() { hashCuckoo.Obtener("Miu")})
	require.PanicsWithValue(t,Hash.ERROR_NO_ESTABA,func() { hashCuckoo.Obtener("Inaba")})
	require.PanicsWithValue(t,Hash.ERROR_NO_ESTABA,func() { hashCuckoo.Obtener("Zero Two")})

	require.EqualValues(t,1,hashCuckoo.Obtener("Uno"))
	require.EqualValues(t,2,hashCuckoo.Obtener("Dos"))
	require.EqualValues(t,3,hashCuckoo.Obtener("Tres"))
}



func TestBorrar(t *testing.T){
	hashCuckoo := Hash.CrearHash[string,int]()

	require.NotNil(t,hashCuckoo)

	hashCuckoo.Guardar("Uno",1)
	hashCuckoo.Guardar("Dos",2)
	hashCuckoo.Guardar("Tres",3)
	hashCuckoo.Guardar("Makise",10)
	hashCuckoo.Guardar("Inaba",9)

	require.EqualValues(t,5,hashCuckoo.Cantidad(),"Cantidad incorrecta")

	
	require.EqualValues(t,10,hashCuckoo.Borrar("Makise"))
	require.EqualValues(t,1,hashCuckoo.Borrar("Uno"))

	require.PanicsWithValue(t,Hash.ERROR_NO_ESTABA,func() { hashCuckoo.Borrar("Makise")})
	require.PanicsWithValue(t,Hash.ERROR_NO_ESTABA,func() { hashCuckoo.Borrar("Misaka")})
	require.PanicsWithValue(t,Hash.ERROR_NO_ESTABA,func() { hashCuckoo.Obtener("Makise")})
	require.PanicsWithValue(t,Hash.ERROR_NO_ESTABA,func() { hashCuckoo.Obtener("Uno")})

	require.EqualValues(t,9,hashCuckoo.Obtener("Inaba"))
	require.EqualValues(t,2,hashCuckoo.Obtener("Dos"))
	require.EqualValues(t,3,hashCuckoo.Obtener("Tres"))
	require.EqualValues(t,3,hashCuckoo.Cantidad(),"Cantidad incorrecta tras borrar")
}


