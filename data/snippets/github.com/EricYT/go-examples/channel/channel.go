package channel

import (
	"fmt"
	"time"
)

// func pinger(c chan<- string){
func pinger(c chan string) {
	for i := 0; ; i++ {
		c <- "ping"
	}
}

func ponger(c chan string) {
	for i := 0; ; i++ {
		c <- "ponger"
	}
}

// func printer(c <-chan string){
func printer(c chan string) {
	count := 0
	for {
		str := <-c
		fmt.Println("receive channel: ", count, str)
		time.Sleep(time.Second * 1)
		if count == 10 {
			break
		} else {
			count += 1
		}
	}
}

/*
func main() {
	var c chan string = make(chan string)

	go pinger(c)
	go ponger(c)
	go printer(c)

	// chan
	var c1 chan string = make(chan string)
	var c2 chan string = make(chan string)

	go func() {
		for {
			c1 <- "ping 1"
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for {
			c2 <- "ping 2"
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			select {
			case msg1 := <-c1:
				fmt.Println("receive c1:", msg1)
			case msg2 := <-c2:
				fmt.Println("receive c2:", msg2)
			}
		}
	}()

	var input string
	fmt.Scanln(&input)

}
*/
