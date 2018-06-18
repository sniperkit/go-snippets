package main

import "fmt"

type Equaler interface {
	Equal(interface{}) bool
}

type Point struct {
	X, Y int
}

func (p *Point) Equal(v interface{}) bool {

	q, ok := v.(*Point)
	if !ok {
		return false
	}

	if p.X != q.X {
		return false
	}
	if p.Y != q.Y {
		return false
	}

	return true
}

type Vertex struct {
	X, Y int
}

func (p *Vertex) Equal(v interface{}) bool {

	q, ok := v.(*Vertex)
	if !ok {
		return false
	}

	if p.X != q.X {
		return false
	}
	if p.Y != q.Y {
		return false
	}

	return true
}

func equalSamples(a, b Equaler) {

	var s string

	if a.Equal(b) {
		s = "=="
	} else {
		s = "!="
	}

	fmt.Printf("%v %s %v\n", a, s, b)
}

func main() {

	var es = []Equaler{
		&Point{12, 0},
		&Point{-2, 14},
		&Point{-2, 14},
		&Vertex{12, 0},
		&Vertex{12, 0},
	}

	n := len(es)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			equalSamples(es[i], es[j])
		}
	}
}
