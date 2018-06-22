package main

import "fmt"
import "unsafe"
import "sync"
import "time"

func uintptrExample() {
	a := 1
	var p uintptr
	p = uintptr(unsafe.Pointer(&a))
	fmt.Printf("%v\n", p)
}

func main() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		fmt.Printf("Removed from queue, before len: %d\n", len(queue))
		queue = queue[1:]
		fmt.Printf("Removed from queue, current len: %d\n", len(queue))
		// c.Signal()
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}

		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}
