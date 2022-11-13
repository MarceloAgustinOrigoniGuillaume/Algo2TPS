package cola_prioridad

import "fmt"

func CumplePropiedadDeHeap[T comparable](elementos []T, cmp func(T, T) int, ultimoIndice int) string {
	res := ""
	for i := 0; i <= (ultimoIndice+1)>>1; i++ {
		candidato := (i << 1) + 1 // izq?
		//fmt.Printf("\n ??? ind:%d = %v\n",i,elementos[i])
		if candidato <= ultimoIndice && cmp(elementos[i], elementos[candidato]) < 0 {
			res += fmt.Sprintf("ind:%d = %v y hijo izq ind:%d = %v, no cumple\n", i, elementos[i], candidato, elementos[candidato])
		} else if candidato <= ultimoIndice {
			//fmt.Printf("izq ind:%d = %v\n",candidato,elementos[candidato])

		}

		candidato++ // der

		if candidato <= ultimoIndice && cmp(elementos[i], elementos[candidato]) < 0 {
			res += fmt.Sprintf("ind:%d = %v y hijo der ind:%d = %v, no cumple\n", i, elementos[i], candidato, elementos[candidato])
		} else if candidato <= ultimoIndice {
			//fmt.Printf("der ind:%d = %v\n",candidato,elementos[candidato])
		}
	}
	return res
}
