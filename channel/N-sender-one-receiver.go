package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	log.SetFlags(0)
	log.Println("N sender and one receiver go")

	// plant a seed and watch it grow
	rand.Seed(time.Now().UnixNano())

	//
	const MaxRandomNumber = 1000
	const SenderNumber = 1000

	//
	var wg sync.WaitGroup
	wg.Add(1)

	dataCh := make(chan int, 100)
	stopCh := make(chan bool)

	for i := 0; i < SenderNumber; i++ {
		go func() {
			for {
				value := rand.Intn(MaxRandomNumber)
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}()
	}

	go func() {
		defer wg.Done()
		for {
			for value := range dataCh {
				if value == MaxRandomNumber-1 {
					// It is safe to close stopCh here
					close(stopCh)
					return
				} else {
					log.Println(value)
				}
			}
		}
	}()

	wg.Wait()
}
