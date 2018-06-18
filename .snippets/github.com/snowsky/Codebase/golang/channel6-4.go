package main

import (
	"fmt"
	"runtime"
	"sync"
)

// var a int = 1

type aMutex struct {
	a     int
	mutex sync.Mutex
}

var a = aMutex{a: 1, mutex: sync.Mutex{}}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go GO(c, i)
	}
	for i := 0; i < 10; i++ {
		<-c
	}
}

func GO(c chan bool, index int) {
	// a := 1
	a.mutex.Lock()
	a.a += 1
	fmt.Println(index, a)
	c <- true
	a.mutex.Unlock()

}
