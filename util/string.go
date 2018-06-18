package util

// Collection of string functions
//
// Made by Vegard Løkken <vegard@loekken.org>
// License MIT

// Contains checks wether or not a rune is present in a string
// It returns a boolean value of true or false
func StringContains(search rune, haystack string) bool {
	for _, v := range haystack {
		if v == search {
			return true
		}
	}
	return false
}
