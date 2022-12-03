package lib

import hash "tp3/diccionario"

type UnionFind[T comparable] struct {
	conjuntos []int
	indices   hash.Diccionario[T, int]
}

func CrearUnionFind[T comparable](elementos []T) *UnionFind[T] {
	union := UnionFind[T]{make([]int, len(elementos)), hash.CrearHash[T, int]()}

	for ind, elem := range elementos {
		union.conjuntos[ind] = ind
		union.indices.Guardar(elem, ind)
	}

	return &union

}

func (union *UnionFind[T]) Find(elem1 T) int {
	return union.find(union.indices.Obtener(elem1))
}

func (union *UnionFind[T]) find(indice int) int {
	if union.conjuntos[indice] == indice {
		return indice
	}

	res := union.find(union.conjuntos[indice])
	union.conjuntos[indice] = res // aplanar
	return res
}

func (union *UnionFind[T]) Unite(elem1, elem2 T) bool {

	ind1 := union.Find(elem1)
	ind2 := union.Find(elem2)

	if ind1 != ind2 {
		union.conjuntos[ind1] = ind2
		return true
	}

	return false
}
