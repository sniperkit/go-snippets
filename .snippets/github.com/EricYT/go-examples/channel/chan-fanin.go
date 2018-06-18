package main

import "fmt"
import "time"

func fanIn(input1, input2 <-chan byte) <-chan byte {
	tmpCh := make(chan byte)

	go func() {
		for {
			select {
			case res := <-input1:
				tmpCh <- res
				time.Sleep(time.Duration(3))
			case res := <-input2:
				tmpCh <- res
			}
		}
	}()

	//	go func() {
	//		for {
	//			fmt.Println("1 receive input")
	//			time.Sleep(1 * time.Second)
	//			tmpCh <- <-input1
	//		}
	//	}()
	//	go func() {
	//		for {
	//			fmt.Println("2 receive input")
	//			time.Sleep(1 * time.Second)
	//			tmpCh <- <-intput2
	//		}
	//	}()
	//
	/*
		defer func() {
			fmt.Println("function end")
			close(tmpCh)
		}()
	*/
	return tmpCh
}

func boring(msg string) <-chan byte {
	tmpCh := make(chan byte, len(msg))

	for i := 0; i < len(msg); i++ {
		tmpCh <- msg[i]
	}

	close(tmpCh)
	return tmpCh
}

func main() {
	resCh := fanIn(boring("hello"), boring("world"))

	for i := 0; i < 10; i++ {
		res := <-resCh

		fmt.Printf("main loop receive index:%d res:%d\n", i, res)
	}
}
