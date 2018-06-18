package main

import "github.com/e154/go-candyjs"

type MyStruct struct {
	Number int
}

func (m *MyStruct) Multiply(x int) int {
	return m.Number * x
}

func main() {
	ctx := candyjs.NewContext()
	ctx.PushGlobalStruct("golangStruct", &MyStruct{10})

	ctx.EvalString(`print(golangStruct.number);`) //10
	ctx.EvalString(`print(golangStruct.multiply(5));`) //50
}
