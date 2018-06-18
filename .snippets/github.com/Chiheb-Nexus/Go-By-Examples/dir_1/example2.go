package main

import (
	"fmt"
	"math/cmplx"
)

var (
	goIsFun   bool       = true
	maxInt    uint64     = 1<<64 - 1
	v_complex complex128 = cmplx.Sqrt(-5 + 2i)
)

func main() {
	/*
		// We can also write variables like this:

		goIsFun  	:= true
		maxInt		:= uint64(1<<64 - 1)
		v_complex	:= cmplx.Sqrt(-5 + 2i)
	*/

	const f = "%T (%v)\n"
	fmt.Printf(f, goIsFun, goIsFun)
	fmt.Printf(f, maxInt, maxInt)
	fmt.Printf(f, v_complex, v_complex)
}
