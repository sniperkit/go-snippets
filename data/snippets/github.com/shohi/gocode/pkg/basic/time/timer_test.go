package time

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	log.Println("world")
	time.AfterFunc(200*time.Millisecond, func() {
		log.Println("hello")
		wg.Done()
	})

	wg.Wait()
}
