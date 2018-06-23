package main

import (
	"fmt"
	"math/rand"
	"time"
)

var channel chan int64 = make(chan int64, 32)

func init() {
	go func() {
		var old int64
		for {
			o := rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
			if old != o {
				old = o
				select {
				case channel <- o:
				}
			}
		}
	}()
}

func RandInt64() (r int64) {
	select {
	case rand := <-channel:
		r = rand
	}
	return
}

func main() {
	var n int64
	n = RandInt64()
	fmt.Println(n)
}
