package main

import "fmt"
import "sync/atomic"
import "sync"

func main() {
	var v uint64 = 0
	fmt.Println("v is ", v)
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if atomic.CompareAndSwapUint64(&v, 0, 1) {
				fmt.Println("modify success")
			} else {
				fmt.Println("modify faild")
			}
		}()
	}
	wg.Wait()
}
