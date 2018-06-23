package main

import "fmt"

type Integer int

func (i *Integer) Print() {
	fmt.Println(*i)
}

func main() {

	var a, b, c Integer
	a = 1
	b = 2
	c = a + b

	c.Print()
}
