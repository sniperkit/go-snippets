package main

import "fmt"

type Vertex struct {
	X, Y int
}

func (Vertex) Test1() {
	fmt.Println("test1")
}

func (c Vertex) Test2() {
	fmt.Println("test2", c)
}

func (c *Vertex) Test3() {
	fmt.Println("test3", c)
}

func main() {
	fmt.Println(Vertex{1,2})

	//Struct Fields
	v := Vertex{1,2}
	v.X = 4
	fmt.Printf("v.X = %d",v.X)

	//Pointers
	p := Vertex{1,2}
	q := &p
	q.X = 1e9
	fmt.Println("Pointers")
	fmt.Println(p.X)

	//Struct Literals
	var (
		a = Vertex{1,2}		//has type Vertex
		s = &Vertex{1,2}	//has type *Vertex
		d = Vertex{X: 1}	//Y:0 is implicit
		f = Vertex{}		//X:0 and Y:0
	)
	fmt.Println("")
	fmt.Println("Struct Literals")
	fmt.Println(a, s, d, f)

	//The new function
	//The expression new(T) allocates a zeroed T value and returns a pointer to it.
	// var t *T = new(T)
	// or
	// t := new(T)
	fmt.Println("\n--------")
	fmt.Println("The new function")
	g := new(Vertex)
	fmt.Println(g)

	fmt.Println("\n--------")
	g.X, g.Y = 11, 9
	fmt.Println(g)

	fmt.Println("\n--------")
	p.Test1()
	p.Test2()
	p.Test3()
}
