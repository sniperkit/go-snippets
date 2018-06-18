package main

import "fmt"

type Car struct {
	Type  string
	Doors int
}

func (c Car) Info() {
	fmt.Println("Type: ", c.Type)
	fmt.Println("Doors: ", c.Doors)
}

func main() {
	bmw := Car{"M5", 5}
	bmw.Info()
}
