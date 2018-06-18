package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	log.SetFlags(0)

	// plant a seed, and watch it grow
	rand.Seed(time.Now().UnixNano())

	//
	const MaxRandomNumber = 1000
	const NumReceiver = 100

	ch := make(chan int, 100)

	var wg sync.WaitGroup
	wg.Add(NumReceiver)

	// one sender
	go func() {
		for {
			if value := rand.Intn(MaxRandomNumber); value == 0 {
				close(ch)
				return
			} else {
				ch <- value
			}
		}
	}()

	for i := 0; i < NumReceiver; i++ {
		go func() {
			defer wg.Done()
			for value := range ch {
				log.Println(value)
			}
		}()
	}

	wg.Wait()
}
