package main

import (
	"fmt"
	"runtime"
)

func printer(c chan int) {
	for i := 0; i < 10; i++ {
		c <- i
		// time.Sleep(time.Second)
	}
	close(c)
}

func reciever(c chan int) {
	for recievedMsg := range c {
		fmt.Println(recievedMsg)
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	newChanel := make(chan int)
	go printer(newChanel)
	reciever(newChanel)
}
