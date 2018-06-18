package main

import (
	"fmt"
)

func do(s string) string {
	if length := len(s); length > 10 {
		return fmt.Sprintf("%c%d%c", s[0], length-2, s[length-1])
	}
	return s
}

// func Test(t *testing.T) {
// 	a := assert.New(t)

// 	a.Equal("word", do("word"))
// 	a.Equal("l10n", do("localization"))
// 	a.Equal("i18n", do("internationalization"))
// 	a.Equal("p43s", do("pneumonoultramicroscopicsilicovolcanoconiosis"))
// 	a.Equal("abcdefghij", do("abcdefghij"))
// }

func main() {
	var n int
	var s string
	fmt.Scan(&n)
	for n > 0 {
		fmt.Scan(&s)
		fmt.Println(do(s))
		n--
	}
}
