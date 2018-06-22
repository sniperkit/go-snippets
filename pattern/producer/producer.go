package producer

import (
	"log"

	"golang.org/x/net/context"
)

// dataCh is a channel for storing data producer generates
var dataCh = make(chan int)

// fbCh is a channel for holding feedback from consumer
// if consumer consume a data, then sends a feedback.
// when Producer receives the feedback, then start producing again
var fbChan = make(chan struct{})

func produce(ctx context.Context) {
	var sn = 0
	// produce one first to start whole process
	// otherwise all goroutine hangs
	dataCh <- sn
	for {
		select {
		case <-ctx.Done():
			log.Println("context is done. Produce exit....")

			// read last item if needs
			go func() {
				<-fbChan
			}()
			return
		case <-fbChan:
			dataCh <- sn
			sn++
		}
	}
}

func consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("context is done. Consume exit....")

			// read last item if needed
			go func() {
				<-dataCh
			}()

			return
		case data := <-dataCh:
			log.Printf("consume %v\n", data)
			fbChan <- struct{}{}
		}
	}
}
