package hash

type IterDiccionario[K comparable, V any]{
	Siguiente() (K,V)
}

type Hash[K comparable, V any] interface{

	Guardar(clave K, valor V)
	Pertence(clave K) bool
	Obtener(clave K) V
	Borrar(clave K) V
	Cantidad() int
	Iterar(func(clave K, dato V) bool)
	Iterador() IterDiccionario[K,V]


}