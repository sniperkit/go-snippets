package main

import (
	"fmt"
	"unicode/utf8"
)

func StringToRunes1(s string) []rune {
	return []rune(s)
}

func StringToRunes2(s string) (rs []rune) {
	for _, r := range s {
		rs = append(rs, r)
	}
	return
}

func StringToRunes3(s string) (rs []rune) {

	bs := []byte(s)
	n := utf8.RuneCount(bs)
	rs = make([]rune, n)

	for i := 0; i < n; i++ {
		r, size := utf8.DecodeRune(bs)
		bs = bs[size:]
		rs[i] = r
	}

	return
}

func StringToRunes4(s string) (rs []rune) {

	bs := []byte(s)

	for {
		r, size := utf8.DecodeRune(bs)
		if size == 0 {
			break
		}

		bs = bs[size:]
		rs = append(rs, r)
	}

	return
}

func main() {
	str := "123"
	fmt.Println(StringToRunes1(str))
}
