package main

import (
	"log"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// routines
	for i := 0; i < 2; i++ {
		go func(index int) {
			defer wg.Done()
			runtime.LockOSThread()
			ticker := time.NewTicker(time.Second * 1)
			defer ticker.Stop()
			var count int
			for {
				select {
				case <-ticker.C:
					log.Printf("[Goroutine %d] ticker cpu: %d", index, syscall.Gettid())
					if index == 1 {
						count++
						if count == 10 {
							log.Printf("[Goroutine %d] unlockOSThread original cpu: %d", index, syscall.Gettid())
							runtime.UnlockOSThread()
							log.Printf("[Goroutine %d] unlockOSThread now cpu: %d", index, syscall.Gettid())
						}
					}
				}
			}
		}(i)
	}

	log.Println("[Goroutine main] cpu: ", syscall.Gettid())
	wg.Wait()
}
