package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("vim-go")

	tick := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case t := <-tick.C:
				fmt.Println("time is ", t)
			}
		}
	}()

	// block the main process
	select {}
}
