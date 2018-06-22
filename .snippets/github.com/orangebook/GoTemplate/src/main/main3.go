package main

import "fmt"

func Count(ch chan int) {
	fmt.Println("Counting")
	ch <- 1
}

func SeleceHandler(chan1, chan2 chan int) {
	select {
	case <-chan1:
		fmt.Println("chan1")
	case <-chan2:
		fmt.Println("chan2")
	default:
		fmt.Println("default")
	}
}

func main() {
	//chs := make([]chan int, 10)
	//fmt.Println("length", cap(chs));
	//for i := 0; i < 10; i++ {
	//	chs[i] = make(chan int)
	//	go Count(chs[i])
	//}
	//
	//for _, ch := range (chs) {
	//	<-ch
	//}

	ch := make(chan int, 1)
	select {
	case ch <- 0:
		fmt.Println("select input running")
	case ch <- 1:
		fmt.Println("select input running")
	}
	//i := <-ch
	select {
	case <-ch:
		fmt.Println("select read running")
	}
}