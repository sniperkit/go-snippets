package main

import (
	"fmt"
	"runtime"
	"time"
)

var c chan int

func ready(w string, sec int) {
	time.Sleep(time.Duration(sec) * time.Second)
	fmt.Println(w, "is ready!")
	c <- sec
}

func main() {
	i := 0
	runtime.GOMAXPROCS()
	c := make(chan int, 2)
	go ready("Tea", 2)
	go ready("Coffee", 1)
	fmt.Println("I'm waiting")
L:
	for {
		select {
		case <-c:
			i++
			if i > 1 {
				break L
			}
		}
	}
	defer close(c)
}
