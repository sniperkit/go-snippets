package main

import (
	"fmt"

	"github.com/EricYT/go-examples/chaining-methods/circle"
)

func main() {
	c := circle.NewCircleBuilder().X(1.0).Y(2.0).Radius(2.0).Finalize()
	fmt.Println("area:", c.Area())
	fmt.Println("x:", c.X())
	fmt.Println("y:", c.Y())
}
