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
	runtime.GOMAXPROCS(1)
	c = make(chan int, 2)
	go ready("Tea", 2)
	go ready("Coffee", 1)
	fmt.Println("I'm waiting")
	// time.Sleep(time.Duration(5) * time.Second)
	<-c
	<-c
	defer close(c)
}
