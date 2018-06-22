package lex

// ListSubstrs lists all the substrings of a word of given length.
func ListSubstrs(word string, length int) []string {
	// Cast word as a rune slice.
	runes := []rune(word)
	n := len(runes)
	subcnt := n - length + 1
	if subcnt < 0 {
		subcnt = 0
	}
	substrs := make([]string, subcnt)
	if n >= length {
		max := n - length
		for i := 0; i <= max; i++ {
			// Cast portion of rune slice back to string.
			substr := string(runes[i : i+length])
			substrs = append(substrs, substr)
		}
	}
	return substrs
}
