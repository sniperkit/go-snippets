package main

import (
	"sync"
)

var (
	myMap map[string]string
	Mutex sync.Mutex
)
func main() {
	myMap = make(map[string]string)

	for i := 0; i < 1000; i++ {
		go func() {
			Mutex.Lock()
			// defer Mutex.Unlock()
			myMap["key"] = "value"
			Mutex.Unlock()
		}()
	}

	//for {
	//	time.Sleep(time.Second)
	//}
}
