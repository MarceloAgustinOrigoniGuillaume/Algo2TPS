package utils

func ordenar[T any](arr []T, inicio int, final int, cmp func(T, T) bool) ([]T, int) {
	pivot := arr[final]
	i := inicio
	for j := inicio; j < final; j++ {
		if cmp(arr[j], pivot) { // arr[j] < pivot returns true
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[final] = arr[final], arr[i]
	return arr, i
}

func QuickSort[T any](arr []T, low, high int, cmp func(T, T) bool) []T {
	if low < high {
		var p int
		arr, p = ordenar(arr, low, high, cmp)
		arr = QuickSort(arr, low, p-1, cmp)
		arr = QuickSort(arr, p+1, high, cmp)
	}
	return arr
}
