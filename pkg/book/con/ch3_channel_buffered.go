package main

import (
	"fmt"
	"strconv"

	"github.com/shohi/gocode/util"
)

func main() {
	var ch chan string
	ch = make(chan string, 2)
	// ch = make(chan string)
	begin := make(chan interface{})
	go func() {
		defer close(begin)
		defer close(ch)
		for k := 0; k < 3; k++ {
			fmt.Printf("Goroutine ID ==> %d, value ==> %s\n", util.GoID(), strconv.Itoa(k))
			ch <- strconv.Itoa(k)
		}
	}()

	<-begin
	for val := range ch {
		fmt.Printf("Goroutine ID ==> %d, value ==> %s\n", util.GoID(), val)
	}
}
