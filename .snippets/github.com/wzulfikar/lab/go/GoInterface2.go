package main

import "fmt"

func main() {
	c := circle{}
	s := square{}
	t := triangle{}

	// passing exactly 1 shape
	fmt.Println(getAreaOfAnyShape(c))

	// passing more than 1 shape
	fmt.Println(getAreaOfAnyShape(c, s, t))

	// passing no shape
	fmt.Println(getAreaOfAnyShape())
}

// interface: set of methods
type shape interface {
	area() string
}

type circle struct{}
type square struct{}
type triangle struct{}

// using interface as function param
// the `...` means we can pass zero to many shape
func getAreaOfAnyShape(shapes ...shape) []string {
	// use slice for dynamic allocation
	areas := make([]string, len(shapes))
	for key, s := range shapes {
		areas[key] = s.area()
	}
	return areas
}

// `circle`, `square` and `triangle`
// (implicitly) implements `shape` interface
// because it has all the methods defined in `shape` interface

func (c circle) area() string {
	return "area of circle is: pi"
}

func (s square) area() string {
	return "area of square is: p x l"
}

func (t triangle) area() string {
	return "area of triangle is: p x l x t"
}
