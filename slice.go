package main

import (
	"fmt"
)

func main() {
	s := make([]string, 3)
	fmt.Println("emp: ", s)
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println(cap(s))
	fmt.Println(s)

	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println(cap(s))
	fmt.Println(s)

	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("cpy: ", c)

	l := s[1:]
	fmt.Println(l)
}
