package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	log.SetFlags(0)
	log.Println("N senders and M receivers go")

	// plant a seed, and watch it grow
	rand.Seed(time.Now().UnixNano())

	//
	const MaxRandomNumber = 1000000
	const NumReceiver = 10
	const NumSender = 1000

	//
	var wg sync.WaitGroup
	wg.Add(NumReceiver)

	//...
	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})

	toStop := make(chan string, 1)

	var stoppedBy string

	// moderator
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// senders
	for i := 0; i < NumSender; i++ {
		go func(id string) {
			for {
				value := rand.Intn(MaxRandomNumber)
				if value == 0 {
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}

				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// receivers
	for i := 0; i < NumReceiver; i++ {
		go func(id string) {
			defer wg.Done()
			for {
				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == MaxRandomNumber-1 {
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	wg.Wait()
	log.Println("stopper is: ", stoppedBy)
}
