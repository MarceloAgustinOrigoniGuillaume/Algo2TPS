package test_ground_test

import (
	"fmt"
	Testeando "test_ground"
	"testing"
)

func Test1(t *testing.T) {
	inst := Testeando.CrearTestStruct[Testeando.Interfaz2]()

	inst.PoneleDato(Testeando.CrearPersona("Miu", 18))

	t.Log(fmt.Sprintf("obtuvo %v", inst.DameDato()))

	inst.DameDato().Edad()
}
