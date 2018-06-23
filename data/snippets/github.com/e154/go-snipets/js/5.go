package main

import (
	"github.com/e154/go-candyjs"
	"fmt"
)

const script5 =  `
"use strict";

print('script5');

function f(x){
  var a = 12;
  b = a + x*35; // error!
}
f();
`

const script6 =  `
	print('script6');
`


func main() {

	ctx := candyjs.NewContext()
	defer func() {
		ctx.Destroy()
	}()

	if err := ctx.PevalString(script5); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := ctx.PevalString(script6); err != nil {
		fmt.Println(err.Error())
		return
	}

	result := ctx.SafeToString(-1)
	fmt.Println("result:",result)
}