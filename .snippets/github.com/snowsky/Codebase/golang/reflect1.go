package main

import (
	"fmt"
	"reflect"
)

type s struct {
}

func (e *s) Error() string {
	return "a string"
}

func main() {
	x := (*s)(nil)
	if x == nil {
		fmt.Printf("x == nil, type of x = %s\n", reflect.TypeOf(x))
	}

	y := error(x)
	if y != nil {
		fmt.Printf("y != nil, type of y = %s\n", reflect.TypeOf(y))
	}
}
