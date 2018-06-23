package main

import (
	"fmt"
	"runtime"
)

func x() *string {
	var x *string = new(string)
	*x = "hello"

	runtime.SetFinalizer(x, func(i *string) {
		fmt.Println("set finalizer x:", i)
	})

	return x
}

func main() {
	y := x()

	fmt.Println("main y:", *y)
}
