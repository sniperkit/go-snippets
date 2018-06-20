// +build fifth

package main

import (
	"fmt"
	"os"
)

// BaseGreeter :
type BaseGreeter struct{}

// HelloTo :
func (f *BaseGreeter) HelloTo(name string) {
	fmt.Printf("hello %v\n", name)
}

// World :
type World struct {
	BaseGreeter
}

// Hello :
func (w *World) Hello() {
	w.HelloTo("world")
}

// Thing :
type Thing struct {
	BaseGreeter
	Name string
}

// Hello :
func (t *Thing) Hello() {
	t.HelloTo(t.Name)
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
