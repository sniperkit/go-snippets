package main

import "fmt"

var ch chan int = make(chan int)

func foo() {
	ch <- 0 // 向ch中加数据，如果没有其他goroutine来取走这个数据，那么挂起foo, 直到main函数把0这个数据拿走
}

func main() {
	go foo()
	x := <-ch // 从ch取数据，如果ch中还没放数据，那就挂起main线，直到foo函数中放数据为止
	fmt.Println(x)
}
