package main

import (
	"fmt"
	"strconv"
)

func main() {
	quit := make(chan int)
	output := sleepCounterWithExit("MySleeper", quit)

	for {

		select {
		case s := <-output:
			fmt.Println(s)
			if s == 100000 {
				quit <- -1
				fmt.Println("sending quit...")
			}
		case x := <-quit:
			fmt.Println("stop it now... " + strconv.Itoa(x))
			return
		}
	}

}

func sleepCounterWithExit(name string, quit chan int) <-chan int {

	c := make(chan int)

	go func() {
		for i := 1; i < 1000010; i++ {

			select {
			case c <- i:
			case x := <-quit:
				if x == -1 {
					quit <- i
				}
			}
		}
	}()

	return c

}
