package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	head := make(chan string)
	tail := make(chan string)

	go func() { head <- "heads" }()
	go func() { tail <- "tails" }()

	rand.Seed(time.Now().UTC().UnixNano())
	r := rand.Intn(2)

	if r == 1 {
		head = nil
	} else {
		tail = nil
	}

	select {
	case s := <-head:
		fmt.Println(s)
	case s := <-tail:
		fmt.Println(s)
	}

}
