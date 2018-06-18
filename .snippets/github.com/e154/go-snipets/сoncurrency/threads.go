package main

import (
	"fmt"
	"time"
)
var ch1 chan int = make(chan int);

func foo(ch chan int) {
	for {fmt.Println(<-ch)}
}

func bar(ch chan int) {
	for i := 0; ; i++ {
		time.Sleep(1000000000);
		ch <- i
	}
}

func main(){
	go bar(ch1)
	go foo(ch1)
	time.Sleep(10000000000)
}
