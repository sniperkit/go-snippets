package basic

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestSyncCondition(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	var conditionTrue = false
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		cond.L.Lock()
		for conditionTrue == false {
			log.Println("Wait for signal")
			cond.Wait()
			log.Println("Waked Up")
		}
		cond.L.Unlock()
	}()

	go func() {
		defer wg.Done()
		cond.L.Lock()
		time.AfterFunc(1*time.Second, func() {
			conditionTrue = true
			cond.Signal()
			log.Println("Send signal")
		})
		cond.L.Unlock()
	}()

	wg.Wait()
}

func TestSyncReentranceLock(t *testing.T) {
	var mu sync.Mutex
	mu.Lock()
	mu.Lock()

	mu.Unlock()
	mu.Unlock()
}
