package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "wzulfikar@two.com"
	fmt.Println(maskString(s, "*"))

	s = "wzulfikarwzulfikar@two.com"
	fmt.Println(maskString(s, "*"))

	s = "One"
	fmt.Println(maskString(s, "*"))

	s = "Yunian Safaru"
	fmt.Println(maskString(s, "*"))

	s = "Yuni"
	fmt.Println(maskString(s, "*"))

	s = "u@as.co"
	fmt.Println(maskString(s, "*"))

	s = "Jay Ho"
	fmt.Println(maskString(s, "*"))

	s = "Satoshi Nakamoto"
	fmt.Println(maskString(s, "*"))

	s = "Sad Man Alone"
	fmt.Println(maskString(s, "*"))

	s = "One two three four"
	fmt.Println(maskString(s, "*"))
}

func maskString(s, mask string) string {
	s = strings.TrimSpace(s)
	l := len(s)

	maskFrom := l / 3
	maskEnd := l / 4

	// limit maskFrom & maskEnd
	if maskFrom > 4 {
		maskFrom = 4
	}
	if maskEnd > 3 {
		maskEnd = 3
	}

	masked := s[0:maskFrom]
	prev := ""
	for i := maskFrom; i < l; i++ {
		c := s[i : i+1]
		if c == "@" ||
			c == "." ||
			c == " " ||
			i >= l-maskEnd {
			masked += c
		} else if prev == " " {
			masked += c
			maskEnd--
		} else {
			masked += mask
		}
		prev = c
	}

	return masked
}
