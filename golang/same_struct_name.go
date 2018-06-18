package main

import (
	"fmt"
)

type A struct {
	B
	// Name string
}

type B struct {
	C
	Name string
}

type C struct {
	Name string
}

func main() {
	a := A{B: B{Name: "B"}}
	// a := A{B: B{Name: "B"}, C: C{Name: "C"}}
	// a := A{Name: "A", B: B{Name: "B"}}
	fmt.Println(a.Name, a.B.Name)
}
