package main

import (
	"github.com/e154/go-candyjs"
	"fmt"
)

const script2 =  `
var a = "a";
`

const script3 =  `
var b = 23;		`

const script4 =  `

var main;

main = function() {
	print('main', a, b);
	return -1;
};

main();
`

func main() {
	ctx := candyjs.NewContext()

	//ctx.EvalString(script3)
	//ctx.EvalString(script4)

	if err := ctx.PevalString(script2); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := ctx.PevalString(script3); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := ctx.PevalString(script4); err != nil {
		fmt.Println(err.Error())
		return
	}


	result := ctx.SafeToString(-1)
	fmt.Println("result:",result)
}
