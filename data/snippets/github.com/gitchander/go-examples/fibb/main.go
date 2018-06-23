package main

import (
	"fmt"
	"math/big"
)

func Fibb() func() uint64 {
	var a, b uint64 = 0, 1
	return func() uint64 {
		a, b = b, a+b
		return a
	}
}

func BigFibb() func() *big.Int {
	var (
		a = big.NewInt(0)
		b = big.NewInt(1)
		c = new(big.Int)
	)
	return func() *big.Int {
		c.Add(a, b)
		a.Set(b)
		b.Set(c)
		return a
	}
}

func main() {
	next := BigFibb()
	for i := 0; i < 1000; i++ {
		fmt.Println(next())
	}
}
