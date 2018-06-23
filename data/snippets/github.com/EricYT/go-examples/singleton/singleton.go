package main

import (
	"fmt"
	"sync"
	"time"
)

type Foo struct {
}

var _instance *Foo
var _once sync.Once

func GetInstance() *Foo {
	_once.Do(func() {
		fmt.Println("get instance initialize")
		_instance = &Foo{}
	})
	return _instance
}

func main() {
	fmt.Println("main function")

	var count int = 100
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			wg.Done()
			GetInstance()
		}()
	}
	wg.Wait()

	time.Sleep(time.Second * 10)
}
