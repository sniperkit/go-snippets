package main

import (
	"time"
)

func main() {
	waitForever := make(chan interface{})

	go func(){
		time.After
		panic("test panic")
	}()

	<- waitForever
}