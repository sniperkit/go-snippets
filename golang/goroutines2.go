package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup

func ready(w string, sec int) {
	time.Sleep(time.Duration(sec) * time.Second)
	fmt.Println(w, "is ready!")
	defer wg.Done()
}

func main() {
	wg.Add(1)
	runtime.GOMAXPROCS(1)
	go ready("Tea", 2)
	wg.Add(1)
	go ready("Coffee", 1)
	fmt.Println("I'm waiting")
	wg.Wait()
}
