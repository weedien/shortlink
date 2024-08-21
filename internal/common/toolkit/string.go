package toolkit

func Split(value string, delimiter rune) []string {
	var result []string
	for _, v := range value {
		if v == delimiter {
			continue
		}
		result = append(result, string(v))
	}
	return result
}
