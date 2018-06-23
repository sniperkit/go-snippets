package test

import (
	"fmt"
	"testing"
)

func TestStructs(t *testing.T) {

	// 1. How to declare?
	type Point struct {
		X, Y int
	}

	// 2. How to initialize?
	p1 := Point{1, 2}
	fmt.Printf("\n2. p1 = %v\n", p1)

	p2 := &Point{3, 4}
	fmt.Printf("\n2. p2 = %v\n", p2)
	// p2 initialization is equivalent with p3
	p3 := new(Point)
	*p3 = Point{5, 6}
	fmt.Printf("\n2. p3 = %v\n", p3)

	// 3. How to embedd struct types?
	// Circle
	type Circle struct {
		Point
		Radius int
	}
	// Wheel
	type Wheel struct {
		Circle
		Spokes int
	}
	wheel := Wheel{
		Circle: Circle{
			Point: Point{
				X: 1,
				Y: 2,
			},
			Radius: 3,
		},
		Spokes: 4,
	}
	fmt.Printf("3. wheel = %v\n", wheel)
	// or equivalent to
	wheel.Point = Point{
		X: 5,
		Y: 6,
	}
	wheel.Radius = 7
	wheel.Spokes = 8
	fmt.Printf("3. wheel = %v\n", wheel)

}
