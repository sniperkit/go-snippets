/**
bool
string
int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64 uintptr
byte // alias for uint8
rune // alias for int32
     // represents a Unicode code point
float32 float64
complex64 complex128
 */

package main

import (
	"math/cmplx"
	"fmt"
	"math"
)

var (
	ToBe	bool		= false
	MaxInt	uint64		= 1<<64 - 1
	z       complex128	= cmplx.Sqrt(-5 + 12i)
)

func main() {
	const f = "%T(%v)\n"
	fmt.Printf(f, ToBe, ToBe)
	fmt.Printf(f, MaxInt, MaxInt)
	fmt.Printf(f, z, z)

	//Type conversions
	var x, y int = 3, 4
	var fl64 float64 = math.Sqrt(float64(x*x + y*y))
	var z int = int(fl64)
	fmt.Println(x, y, z)


	// interface
	fmt.Println("interface")
	var v interface {}
	v = 7

	switch v.(type) {
	case int, int8, int16, int32, int64:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	}

	n, ok := v.(int)
	if ok {
		fmt.Println(n)
	}

}
