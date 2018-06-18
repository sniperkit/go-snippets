package main

import (
	"github.com/e154/go-candyjs"
)

func main() {

	ctx := candyjs.NewContext()
	ctx.PushGlobalGoFunction("golangMultiply", func(a, b int) int {
		return a * b
	})

	ctx.EvalString(`print(golangMultiply(5, 10));`) //50

}


