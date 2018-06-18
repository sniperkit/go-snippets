package main

import (
	"fmt"
	"strings"
)

func do(s string) string {
	uppercased := strings.ToUpper(s)
	lowercased := strings.ToLower(s)
	if s == uppercased {
		return lowercased
	}
	if s[0] == lowercased[0] && s[1:] == uppercased[1:] {
		return fmt.Sprintf("%c%s", uppercased[0], lowercased[1:])
	}
	return s
}

// func Test(t *testing.T) {
// 	a := assert.New(t)
// 	a.Equal("Caps", do("cAPS"))
// 	a.Equal("Lock", do("Lock"))
// 	a.Equal("cAPSlOCK", do("cAPSlOCK"))
// 	a.Equal("oops", do("OOPS"))
// }

func main() {
	var s string
	fmt.Scan(&s)
	fmt.Println(do(s))
}
