package hash

func ToBytes(objeto interface{}) []byte {

	str, esString := objeto.(string)

	if esString {
		return []byte(str)
	}

	return []byte(fmt.Sprintf("%v", objeto))
}
