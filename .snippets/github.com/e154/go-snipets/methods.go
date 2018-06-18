package main

import (
	"math"
	"fmt"
)

type Vertex struct { X, Y float64}
// данный код определяет метод - Abs(), ассоциированный с Vertex,
// в теле функции получаетель назван - v
func (v *Vertex) Abs() float64 {return math.Sqrt(v.X*v.X + v.Y*v.Y)}
func (v *Vertex) sum() float64 {return v.X + v.Y}
func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

type MyFloat float64
func (f MyFloat) Abs() float64 {
	if f<0 {return float64(-f)}
	return float64(f)
}

func main() {
	v := &Vertex{3,4}
	fmt.Println(v, v.Abs())

	//////////////////
	//Methods continued
	//////////////////
	fmt.Println("")
	fmt.Println("Methods continued")
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())

	//////////////////
	//Methods with pointer receivers
	//////////////////
	fmt.Println("")
	fmt.Println("Methods with pointer receivers")
	v.Scale(5)
	fmt.Println(v, v.Abs())

}
