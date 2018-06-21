package router

// getHashValue is used as compute string's hash value
func getHashValue(str string) (sum int) {
	stringLength := len(str)
	for i := 0; i < stringLength; i++ {
		sum += int(str[i])
	}
	return sum
}
