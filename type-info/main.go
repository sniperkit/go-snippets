package main

import "fmt"

func main() {
	a := 42
	b := 6.543210
	c := "Hello World"
	d := true
	e := &a

	fmt.Printf("%T \n", a)
	fmt.Printf("%T \n", b)
	fmt.Printf("%T \n", c)
	fmt.Printf("%T \n", d)
	fmt.Printf("%T \n", e)

}
