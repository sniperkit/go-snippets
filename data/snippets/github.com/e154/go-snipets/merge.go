package main

import (
	"fmt"
	"github.com/imdario/mergo"
)

type Foo struct {
	A string
	B int64
}

func main() {
	src := Foo{
		A: "one",
	}

	dest := Foo{
		A: "two",
		B: 2,
	}

	mergo.Merge(&dest, src)

	fmt.Println(dest)
	// Will print
	// {two 2}
}