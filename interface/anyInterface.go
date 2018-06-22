package main

import (
	"fmt"
	"strconv"
)

type IDisposable interface {
	dispose() string
}

type Car struct {
	name string
	year int
}

type Bike struct {
	name string
}

func (b Bike) dispose() string {
	return "Bike has been manufactured by " + b.name
}

func (c Car) dispose() string {
	return c.name + " was bought in " + strconv.Itoa(c.year)
}

func PrintAnyInterface(any interface{}) {

	// REF : https://research.swtch.com/interfaces
	// if v, ok := any.(Stringer); ok {
	//     return v.String()
	// }

	switch v := any.(type) {
	case IDisposable:
		fmt.Println(v.dispose())
		return
	default:
		fmt.Println(v)
	}

}

func main() {

	var i = &Car{"Toyota", 2010}
	PrintAnyInterface(i)

	var x = new(Car)
	x.name = "Honda"
	x.year = 2017
	PrintAnyInterface(x)

	var b = &Bike{name: "Suzuki"}
	PrintAnyInterface(b)
}
