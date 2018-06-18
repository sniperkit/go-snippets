package main

import (
	"fmt"
	"regexp"
)

func main() {
	v1 := []byte("0.1.a.4.4.")

	match, err := regexp.Match("^[0-9]+(\\.[0-9]+)*$", v1)
	if err != nil {
		panic(err)
	}
	fmt.Println(match)
}
