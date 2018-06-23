package main

import (
	"fmt"
	"time"
  "strconv"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
    fmt.Println(s)
	}
}

func say2(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(strconv.Itoa(i) + ":"+ s)
	}
}

func main() {
	go say2("world")
	say2("hello")
}
