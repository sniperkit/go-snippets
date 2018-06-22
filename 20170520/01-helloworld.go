package main

import (
	"fmt"

	helloworld "github.com/go-training/workshops/20170520/hello"
)

const (
	A = iota + 1
	B
	C
)

func main() {
	fmt.Println("Hello World!!")
	fmt.Println(helloworld.HelloWorld())

	var a int
	var b = 1
	c := 1
	a = 1
	hello := "Hello world"
	fmt.Println("a:", a)
	fmt.Println("b:", b)
	fmt.Println("c:", c)
	fmt.Println(hello)

	fmt.Println("A:", A)
	fmt.Println("B:", B)
	fmt.Println("C:", C)
}
