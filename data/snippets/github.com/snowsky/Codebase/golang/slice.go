package main

import (
	"fmt"
	"time"
)

func main() {
	t0 := time.Now()

	// Version 1: use accurate capacity for slice.
	for i := 0; i < 10000; i++ {
		values := make([]int, 1000)
		for x := 0; x < 1000; x++ {
			values = append(values, x)
			// fmt.Println(values)
		}
		if values[0] != 0 {
			fmt.Println(0)
		}
	}

	t1 := time.Now()

	// Version 2: use empty slice.
	for i := 0; i < 10000; i++ {
		values := []int{}
		for x := 0; x < 1000; x++ {
			values = append(values, x)
		}
		if values[0] != 0 {
			fmt.Println(0)
		}
	}

	t2 := time.Now()
	// Benchmark results.
	fmt.Println(t1.Sub(t0))
	fmt.Println(t2.Sub(t1))
}
