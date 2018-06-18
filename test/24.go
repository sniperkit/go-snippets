package main

import (
	"fmt"
	"time"
)

// 24. Канал созданный таким образом имеет емкость?
// var chan int

// Ответы:
// не будет работать так как должен быть создан ч помощью вызова make()

func test24(ch chan int) {
	ch <- 1
}

func main() {
	var ch, quit chan int
	ch = make(chan int)
	quit = make(chan int)

	go test24(ch)
	go test24(ch)
	go test24(ch)

	go func() {
		for {
			select {
			case i := <- ch:
				fmt.Println("i", i)
			case <- quit:
				break
			default:

			}
		}
	}()

	time.Sleep(4 * time.Second)
	quit <- 1
}

