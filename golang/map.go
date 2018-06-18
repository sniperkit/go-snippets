package main

import "fmt"

func main() {
	var a map[int]string
	a = map[int]string{}

	var b map[int]string = map[int]string{}

	var c map[int]string
	c = make(map[int]string)

	var d map[int]string = make(map[int]string)

	m := make(map[int]string)

	fmt.Println(a, b, c, d, m)
}
