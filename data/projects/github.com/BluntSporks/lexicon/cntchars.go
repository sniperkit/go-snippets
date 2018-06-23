package lex

// CntChars counts all the characters in a language lexicon.
func CntChars(words map[string]bool) map[rune]int {
	chars := make(map[rune]int)
	for word := range words {
		for _, ch := range word {
			chars[ch]++
		}
	}
	return chars
}
