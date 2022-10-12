package hash_test

import (
	Hash "hash/hashCuckoo2"
	"github.com/stretchr/testify/require"
	"testing"
	"fmt"
)

const ERROR_FUNCION_HASH = "Error: mala funcion de hash, redimension requeridad demasiadas veces seguidas"
const ERROR_NO_ESTABA = "La clave no pertenece al diccionario"
const ERROR_ITERADOR_TERMINO = "El iterador termino de iterar"



func TestAgregado(t *testing.T){
	hash := Hash.CrearHash[string,int]()

	require.NotNil(t,hash,"Crear devolvio nil")


	hash.Guardar("Uno",1)
	hash.Guardar("Dos",2)
	hash.Guardar("Tres",3)

	require.EqualValues(t,3,hash.Cantidad(),"Cantidad incorrecta")
}


func TestPertenece(t *testing.T){
	hash := Hash.CrearHash[string,int]()

	require.NotNil(t,hash)


	hash.Guardar("Uno",1)
	hash.Guardar("Dos",2)
	hash.Guardar("Tres",3)

	require.EqualValues(t,3,hash.Cantidad(),"Cantidad incorrecta")

	require.False(t,hash.Pertenece("O"),"Se invento una clave 'O'")
	require.False(t,hash.Pertenece(""),"Se invento una clave ''")
	require.False(t,hash.Pertenece("Makise"),"Se invento una clave 'Makise'")

	require.True(t,hash.Pertenece("Uno"),"Se olvido de la clave 'Uno'")
	require.True(t,hash.Pertenece("Dos"),"Se olvido de la clave 'Dos'")
	require.True(t,hash.Pertenece("Tres"),"Se olvido de la clave 'Tres'")
}



func TestObtener(t *testing.T){
	hash := Hash.CrearHash[string,int]()

	require.NotNil(t,hash)


	hash.Guardar("Uno",1)
	hash.Guardar("Dos",2)
	hash.Guardar("Tres",3)

	require.EqualValues(t,3,hash.Cantidad(),"Cantidad incorrecta")

	require.PanicsWithValue(t,ERROR_NO_ESTABA,func() { hash.Obtener("Makise")})
	require.PanicsWithValue(t,ERROR_NO_ESTABA,func() { hash.Obtener("Miu")})
	require.PanicsWithValue(t,ERROR_NO_ESTABA,func() { hash.Obtener("Inaba")})
	require.PanicsWithValue(t,ERROR_NO_ESTABA,func() { hash.Obtener("Zero Two")})

	require.EqualValues(t,1,hash.Obtener("Uno"))
	require.EqualValues(t,2,hash.Obtener("Dos"))
	require.EqualValues(t,3,hash.Obtener("Tres"))
}



func TestBorrar(t *testing.T){
	hash := Hash.CrearHash[string,int]()

	require.NotNil(t,hash)

	hash.Guardar("Uno",1)
	hash.Guardar("Dos",2)
	hash.Guardar("Tres",3)
	hash.Guardar("Makise",10)
	hash.Guardar("Inaba",9)

	require.EqualValues(t,5,hash.Cantidad(),"Cantidad incorrecta")

	
	require.EqualValues(t,10,hash.Borrar("Makise"))
	require.EqualValues(t,1,hash.Borrar("Uno"))

	require.PanicsWithValue(t,ERROR_NO_ESTABA,func() { hash.Borrar("Makise")})
	require.PanicsWithValue(t,ERROR_NO_ESTABA,func() { hash.Borrar("Misaka")})
	require.PanicsWithValue(t,ERROR_NO_ESTABA,func() { hash.Obtener("Makise")})
	require.PanicsWithValue(t,ERROR_NO_ESTABA,func() { hash.Obtener("Uno")})

	require.EqualValues(t,9,hash.Obtener("Inaba"))
	require.EqualValues(t,2,hash.Obtener("Dos"))
	require.EqualValues(t,3,hash.Obtener("Tres"))
	require.EqualValues(t,3,hash.Cantidad(),"Cantidad incorrecta tras borrar")
}



func BenchmarkIterador(b *testing.B) {
	b.Log("ESTA HACIENDO EL TEST?")
	n := 12500

	b.Run(fmt.Sprintf("Prueba %d elementos", n+8), func(b *testing.B) {
		hash := Hash.CrearHash[string,int]()
		hash.Guardar("Uno",1)
		hash.Guardar("Dos",2)
		hash.Guardar("Tres",3)
		hash.Guardar("Makise",10)
		hash.Guardar("Inaba",9)
		hash.Guardar("Misaka",9)
		hash.Guardar("Miu",9)
		hash.Guardar("Marcelo",9)	
		i := 0
		defer func(){
			if r:= recover(); r!= nil{
				b.Log(fmt.Sprintf("Err:%s ... key = %08d",r,i))
			}
		}()
		for i < n {
			hash.Guardar(fmt.Sprintf("%08d", i),i)
			i++
		}

	})		
}
