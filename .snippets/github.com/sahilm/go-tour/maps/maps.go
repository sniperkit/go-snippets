package maps

import "strings"

// WordCount returns a map of string->int which contains
// the number of times each words occurs in string s. A space is considered
// to be a word boundary.
func WordCount(s string) map[string]int {
	wordCount := make(map[string]int)
	for _, word := range strings.Fields(s) {
		wordCount[word]++
	}
	return wordCount
}
