package main

import (
	"fmt"
	"unicode"
)

func nextInt(b []byte, i int) (int, int) {
	for ; i < len(b) && !unicode.IsDigit(b[i]); i++ {
	}
	x := 0
	for ; i < len(b) && unicode.IsDigit(b[i]); i++ {
		x = x*10 + int(b[i]) - '0'
	}
	return x, i
}

func main() {
	a := "a100b234"

	for i := 0; i < len(a); i++ {
		x, i := nextInt(a, i)
		fmt.Println(x)
	}
}
