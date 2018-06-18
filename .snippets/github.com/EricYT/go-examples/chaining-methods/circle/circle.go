package circle

import "math"

type Circle struct {
	x      float64
	y      float64
	radius float64
}

func (c Circle) Area() (area float64) {
	return math.Pi * c.radius * c.radius
}

func (c Circle) X() (x float64) {
	return c.x
}

func (c Circle) Y() (y float64) {
	return c.y
}

type CircleBuilder struct {
	x      float64
	y      float64
	radius float64
}

func NewCircleBuilder() *CircleBuilder {
	return &CircleBuilder{
		x:      0.0,
		y:      0.0,
		radius: 1.0,
	}
}

func (cb *CircleBuilder) X(coordinate float64) *CircleBuilder {
	cb.x = coordinate
	return cb
}

func (cb *CircleBuilder) Y(coordinate float64) *CircleBuilder {
	cb.y = coordinate
	return cb
}

func (cb *CircleBuilder) Radius(radius float64) *CircleBuilder {
	cb.radius = radius
	return cb
}

func (cb *CircleBuilder) Finalize() *Circle {
	return &Circle{
		x:      cb.x,
		y:      cb.y,
		radius: cb.radius,
	}
}
