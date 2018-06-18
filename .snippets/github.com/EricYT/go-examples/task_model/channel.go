package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool)
	go func() {
		time.Sleep(time.Millisecond * 1)
		c <- true
	}()

	<-c
	fmt.Println("end")
}
