package main

import (
	"github.com/e154/go-candyjs"
	"fmt"
)

const script8 =  `
print('>>>')
print(IC.test());
print('>>>')
print(IC.test2());
print('>>>')
print(IC.test3.print());
print('>>>')

`

type Test struct {}

func (t *Test) Print() string {
	return "ok3"
}

func addScript1(ctx *candyjs.Context) {
	//ctx := candyjs.NewContext()
	if b := ctx.GetGlobalString("IC"); !b {
		fmt.Println("Object not found")
		return
	}
	ctx.PushObject()
	ctx.PushGoFunction(func() string {
		return "ok"
	})
	ctx.PutPropString(-3, "test")
	ctx.Pop()
}

func addScript2(ctx *candyjs.Context) {
	//ctx = candyjs.NewContext()
	if b := ctx.GetGlobalString("IC"); !b {
		fmt.Println("Object not found")
		return
	}
	ctx.PushObject()
	ctx.PushGoFunction(func() string {
		return "ok2"
	})
	ctx.PutPropString(-3, "test2")
	ctx.Pop()
}

func addStruct(ctx *candyjs.Context) {
	//ctx = candyjs.NewContext()
	if b := ctx.GetGlobalString("IC"); !b {
		fmt.Println("Object not found")
		return
	}
	ctx.PushObject()
	ctx.PushStruct(&Test{})
	ctx.PutPropString(-3, "test3")
	ctx.Pop()
}

func main() {
	ctx := candyjs.NewContext()
	ctx.PevalString(`var IC = {};`)

	addScript1(ctx)

	addScript2(ctx)

	addStruct(ctx)

	ctx.PevalString(script8)
}
