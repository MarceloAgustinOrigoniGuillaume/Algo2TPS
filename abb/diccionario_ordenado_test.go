package diccionario_test

import (
	//"fmt"
	"github.com/stretchr/testify/require"
	TDAABB "abb"
	"testing"
)

func TestUno(t *testing.T){
	abb := TDAABB.CrearABB[string,string](func(str1 string,str2 string) int {
		return len(str1)-len(str2) // ponele
	}) 
	require.NotNil(t,abb)
	require.EqualValues(t,0,abb.Cantidad())
}