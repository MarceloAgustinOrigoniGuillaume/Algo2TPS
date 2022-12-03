package utilities

import "fmt"

func toBytes(objeto interface{}) []byte {

	str, esString := objeto.(string)

	if esString {
		return []byte(str)
	}

	return []byte(fmt.Sprintf("%v", objeto))
}

func _JenkinsHashFunction(bytes []byte) int {
	res := 0
	for i := 0; i < len(bytes); i++ {
		res += int(bytes[i])
		res += res << 10
		res ^= res >> 6
	}

	return res
}

func AplicarFuncionDeHash[K comparable](clave K, maximo int) int { // paso intermedio para hacer mas facil cambios
	return _JenkinsHashFunction(toBytes(clave)) % maximo
}
