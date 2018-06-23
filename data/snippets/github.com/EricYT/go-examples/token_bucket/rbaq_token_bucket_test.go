package token_bucket

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestEnqueue(t *testing.T) {
	q := NewDefaultRateBasedTokenBucketQueue()

	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(time.Second * 2):
				err := q.Enqueue(struct{}{}, 1)
				if err != nil {
					log.Printf("test enqueu: 2 enqueue err: %s\n", err)
				}
			case <-done:
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		var index int
		for {
			time.Sleep(time.Second * 1)
			_, err := q.Dequeue(time.Second * 3)
			if err != nil {
				log.Printf("test enqueu: dequeue err: %s\n", err)
				index++
			}
			log.Printf("test enqueu: dequeue normal \n")
			if index > 100 {
				close(done)
				return
			}
		}
	}()
	wg.Wait()
}
