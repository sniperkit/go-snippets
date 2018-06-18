package main

import (
	"fmt"
	"runtime"
)

var a int = 1
var c chan bool = make(chan bool)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// c := make(chan bool)
	for i := 0; i < 10; i++ {
		go GO(c, i)
	}
	// <-c

	for v := range c {
		fmt.Println(v)
	}
	// close(c)
}

func GO(c chan bool, index int) {
	// a := 1
	for i := 0; i < 100000000; i++ {
		a += i
	}
	fmt.Println(index, a)
	c <- true

	// if index == 9 {
	// 	c <- true
	// }
}
