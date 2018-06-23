package main

import "fmt"

type Cloner interface {
	// Clone() Cloner // Необходимо иметь зависимость

	Clone() interface{} // Метод не имеет зависимости
}

type Alpha struct {
	Id      int
	Message string
}

func (a Alpha) Clone() interface{} {
	c := a
	return c
}

type Betha struct {
	Name     string
	Position int
}

func (b Betha) Clone() interface{} {
	c := b
	return c
}

func main() {

	a := Alpha{
		Id:      345,
		Message: "message for alpha",
	}

	b := Betha{
		Name:     "name betha",
		Position: 23534,
	}

	var cs = []Cloner{a, b}

	for _, c := range cs {
		clone := c.Clone()
		fmt.Printf("%+v\n", clone)
	}
}
