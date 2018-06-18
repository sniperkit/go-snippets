package main

import "fmt"

type person struct {
	name string
	age  int
}

func main() {
	fmt.Println(person{name: "Bob", age: 20})

	s := person{"sean", 22}

	fmt.Println(s)

	fmt.Println(person{name: "bob"})
}
