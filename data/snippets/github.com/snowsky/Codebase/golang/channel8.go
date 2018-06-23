package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
}

func main() {
	runtime.GOMAXPROCS(2) //改成1就不会有任何输出了
	go say("world")       //开一个新的Goroutines执行
	for {
	}
}
