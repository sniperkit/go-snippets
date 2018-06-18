package main

import (
	"fmt"
	"time"
)

var limit chan int = make(chan int, 3)

func main() {

	var count int = 100
	work := make([]func(), count)
	for i := 0; i < count; i++ {
		work[i] = func() {
			time.Sleep(time.Second * 1)
			fmt.Println("index:", i)
		}
	}

	for _, w := range work {
		go func(w func()) {
			limit <- 1
			w()
			<-limit
		}(w)
	}

	fmt.Println("hold on ...............")
	select {}
}
