// +build second

package main

import (
	"fmt"
)

// World :
type World struct{}

// Hello :
func (w *World) Hello() {
	fmt.Println("hello world")
}

// Greeter :
type Greeter = *World

// New :
func New() Greeter {
	return &World{}
}
