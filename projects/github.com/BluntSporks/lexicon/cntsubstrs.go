package lex

// CntSubstrs counts all the substrings of given length in a language lexicon.
func CntSubstrs(words map[string]bool, length int) map[string]int {
	substrs := make(map[string]int)
	for word := range words {
		list := ListSubstrs(word, length)
		for _, substr := range list {
			substrs[substr]++
		}
	}
	return substrs
}
