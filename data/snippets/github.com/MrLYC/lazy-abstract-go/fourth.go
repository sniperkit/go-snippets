// +build fourth

package main

import (
	"fmt"
	"os"
)

// World :
type World struct{}

// Hello :
func (w *World) Hello() {
	fmt.Println("hello world")
}

// Thing :
type Thing struct {
	Name string
}

// Hello :
func (t *Thing) Hello() {
	fmt.Printf("hello %v\n", t.Name)
}

// Greeter :
type Greeter interface {
	Hello()
}

// New :
func New() Greeter {
	name := os.Getenv("name")
	if name == "" {
		return &World{}
	}
	return &Thing{
		Name: name,
	}
}
