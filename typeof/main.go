package main

import (
	"fmt"
	"reflect"
)

func main() {
	i := 42
	s := "a"
	b := true
	p := &i

	info(i)
	info(s)
	info(b)
	info(p)
}

func info(t interface{}) {
	fmt.Println(reflect.TypeOf(t))
}
