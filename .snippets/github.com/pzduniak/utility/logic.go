package utility

func Is(condition bool, yes string, no string) string {
	if condition {
		return yes
	}

	return no
}

func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
