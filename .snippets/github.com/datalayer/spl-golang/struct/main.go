package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Go's _structs_ are typed collections of fields.
// They're useful for grouping data together to form
// records.

// A struct is a collection of fields.
type Vertex struct {
	X int
	Y int
}

// This `person` struct type has `name` and `age` fields.
type person struct {
	name string
	age  int
}

type T struct {
	A int16
	B int8
	C [256]byte
}

func main() {

	// Create a struct and write it.
	t := T{A: 99, B: 10}
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, t)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf)
	// Read into an empty struct.
	t = T{}
	err = binary.Read(buf, binary.BigEndian, &t)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d %d", t.A, t.B)

	// This syntax creates a new struct.
	fmt.Println(person{"Bob", 20})

	// You can name the fields when initializing a struct.
	fmt.Println(person{name: "Alice", age: 30})

	// Omitted fields will be zero-valued.
	fmt.Println(person{name: "Fred"})

	// An `&` prefix yields a pointer to the struct.
	fmt.Println(&person{name: "Ann", age: 40})

	p1 := person{name: "Ann", age: 40}
	fmt.Println(p1)
	hash := []byte(fmt.Sprintf("%v", p1))
	fmt.Println(hash)

	// Access struct fields with a dot.
	s := person{name: "Sean", age: 50}
	fmt.Println(s.name)

	// You can also use dots with struct pointers - the
	// pointers are automatically dereferenced.
	sp := &s
	fmt.Println(sp.age)

	// Structs are mutable.
	sp.age = 51
	fmt.Println(sp.age)

	fmt.Println(Vertex{1, 2})

}
