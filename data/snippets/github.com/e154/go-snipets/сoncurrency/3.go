package main

import (
	"fmt"
	"time"
)

// Каналы обеспечивают возможность общения нескольких горутин друг с другом,
// чтобы синхронизировать их выполнение. Пример программы с использованием каналов:

func pinger(c chan string) {
	for i := 0;;i++ {
		c <- "ping"
//		time.Sleep(time.Second * 1)
	}
}

func printer(c chan string) {
	for {
		msg := <- c
		fmt.Println(msg)
		time.Sleep(time.Second * 1)
	}
}

func main() {

	var c chan string = make(chan string)

	go pinger(c)
	go printer(c)

	var input string
	fmt.Scanln(&input)
}

