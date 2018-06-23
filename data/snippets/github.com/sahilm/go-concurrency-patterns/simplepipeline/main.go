/*
A pipeline where sq wraps the channel emitted by gen. Since the inbound and
outbound channel for sq is the same, we can compose many sqs together. The apply fn
demonstrates this composition.
*/
package main

import (
	"fmt"
)

func gen(nums ...int) chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in chan int) chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

type squarer func(chan int) chan int

func apply(times int, in chan int, fun squarer) chan int {
	out := in
	for i := 0; i < times; i++ {
		out = fun(out)
	}
	return out
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	out := apply(3, gen(nums...), sq)

	for n := range out {
		fmt.Println(n)
	}
}
