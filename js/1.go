package main

import (
	"github.com/e154/go-candyjs"
)

const script1 =  `

  function factorial(n) {
    if (n === 0) return 1;
    return n * factorial(n - 1);
  }

  print(factorial(10));

  (function(a){
  	print(a);
  })(12)

`

func main() {

	ctx := candyjs.NewContext()
	ctx.EvalString(script1)

}


