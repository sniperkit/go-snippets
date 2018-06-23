// +build sixth

package main

import (
	"fmt"
	"os"
	"time"
)

// Base :
type Base struct{}

// HelloTo :
func (f *Base) HelloTo(name string) {
	fmt.Printf("hello %v\n", name)
}

// World :
type World struct {
	Base
}

// Hello :
func (w *World) Hello() {
	w.HelloTo("world")
}

// Thing :
type Thing struct {
	Base
	Name string
}

// Hello :
func (t *Thing) Hello() {
	t.HelloTo(t.Name)
}

// Now :
type Now struct {
	Base
}

// Now :
func (t *Now) Hello() {
	t.HelloTo(time.Now().String())
}

// Greeter :
type Greeter interface {
	Hello()
}

// New :
func New() Greeter {
	name := os.Getenv("name")
	switch name {
	case "":
		return &World{}
	case "now":
		return &Now{}
	default:
		return &Thing{
			Name: name,
		}
	}
}
