package test_ground

import "fmt"

type Interfaz1 interface {
	Nombre() string
}

type Interfaz2 interface {
	Interfaz1
	Edad() int
}

type persona struct {
	nombre string
	edad   int
}

func CrearPersona(nombre string, edad int) Interfaz2 {
	return &persona{nombre, edad}
}

func (per *persona) Nombre() string {
	return per.nombre
}
func (per *persona) Edad() int {
	return per.edad
}

type GenericInterface[T Interfaz1] interface {
	DameDato() T
	PoneleDato(T)
}

type testStruct[T Interfaz1] struct {
	dato T
}

func CrearTestStruct[T Interfaz1]() GenericInterface[T] {
	return &testStruct[T]{}
}

func (t *testStruct[T]) DameDato() T {
	return t.dato
}

func (t *testStruct[T]) PoneleDato(dato T) {
	t.dato = dato

	fmt.Printf("SET DATO %v \n", dato)
}
