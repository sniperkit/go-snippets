package main

import (
	"github.com/e154/go-candyjs"
	"fmt"
)

const script7 =  `
var i = 0;

var main = function() {
	print('i =',i);
	i++;
}
`

func main() {

	ctx := candyjs.NewContext()
	defer func() {
		ctx.DestroyHeap()
		ctx.Destroy()
	}()

	ctx.PevalString(script7)
	ctx.PushGlobalObject()
	if err := ctx.PushTimers(); err != nil {
		fmt.Println("error:", err.Error())
	}

	for i:=0;i<10;i++ {

		//ctx.GetPropString(-1, "main")
		if b := ctx.GetGlobalString("main"); !b {
			fmt.Println("main function not found")
			return
		}

		if r := ctx.Pcall(0); r != 0 {
			result := ctx.SafeToString(-1)
			fmt.Println("error:", result)
			return
		}

		if result := ctx.SafeToString(-1); result != "undefined" {
			fmt.Println("result:", result)
		}
	}


}