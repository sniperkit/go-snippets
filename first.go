// +build first

package main

import (
	"fmt"
)

// Greeter :
type Greeter struct{}

// Hello :
func (w *Greeter) Hello() {
	fmt.Println("hello world")
}

// New :
func New() Greeter {
	return Greeter{}
}
