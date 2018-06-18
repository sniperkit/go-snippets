package main

import "fmt"
import "math"

type geometry interface {
	area() float64
	perim() float64
	show()
}

type rect struct {
	width, height float64
}

func (r rect) show() {
	fmt.Println(r.width, r.height)
}

func (r rect) area() float64 {
	return r.width * r.height
}

func (r rect) perim() float64 {
	return (r.width + r.height) * 2
}

type circle struct {
	radius float64
}

func (c circle) show() {
	fmt.Println(c.radius)
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

func measure(g geometry) {
	fmt.Println(g.area())
	fmt.Println(g.perim())
	g.show()
}

func main() {
	rect := rect{3, 2}
	circle := circle{3}

	measure(rect)
	measure(circle)

}
