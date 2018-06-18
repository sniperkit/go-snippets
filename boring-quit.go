package main

import (
	"fmt"
	"time"
)

func boring(msg string, quit chan bool) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s: %d", msg, i):
				time.Sleep(time.Duration(1000) * time.Millisecond)
			case <-quit:
				return
			}
		}
	}()
}

func main() {
	quit := make(chan bool)
	c := boring("Joe", quit)
	for i := 1000; i >= 0; i-- {
		fmt.Println(<-c)
	}
	quit <- true
}
