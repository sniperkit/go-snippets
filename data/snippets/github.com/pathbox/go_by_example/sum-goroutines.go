package main

import (
	"fmt"
	"runtime"
	"time"
)

var c = make(chan int)

func main() {
	start := time.Now().UnixNano()
	fmt.Println(runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	mysum := 0
	max, num := 10000, 10
	for i := 0; i < max; i++ {
		go sum((max/num)*i+1, (max/num)*(i+1), i)
	}
	for i := 0; i < num; i++ {
		mysum = mysum + <-c
	}
	fmt.Println("sum:", mysum)
	fmt.Println("take times is : ", (time.Now().UnixNano() - start))
}

func sum(min, max, number int) {
	s := 0
	for i := min; i <= max; i++ {
		s = s + 1
	}
	c <- s
}
