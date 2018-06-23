package main

import (
	"fmt"
	"math"
)

type MyInteger int

func (a MyInteger) MyMethod(b int) int {
	return int(a) + b
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func main() {

	var x MyInteger = 1
	fmt.Println(x.MyMethod(2))

	/*
	   You can declare a method on non-struct types, too.
	   In this example we see a numeric type MyFloat with an Abs method.
	   You can only declare a method with a receiver whose type is defined in the same package as the method. You cannot declare a method with a receiver whose type is defined in another package (which includes the built-in types such as int).
	*/

	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())

	/*
	   Pointer receivers
	   You can declare methods with pointer receivers.
	   This means the receiver type has the literal syntax *T for some type T. (Also, T cannot itself be a pointer such as *int.)
	   For example, the Scale method here is defined on *Vertex.
	   Methods with pointer receivers can modify the value to which the receiver points (as Scale does here). Since methods often need to modify their receiver, pointer receivers are more common than value receivers.
	   Try removing the * from the declaration of the Scale function on line 16 and observe how the program's behavior changes.
	   With a value receiver, the Scale method operates on a copy of the original Vertex value. (This is the same behavior as for any other function argument.) The Scale method must have a pointer receiver to change the Vertex value declared in the main function.
	*/v := Vertex{3, 4}
	v.Scale(10)
	fmt.Println(v.Abs())

}
