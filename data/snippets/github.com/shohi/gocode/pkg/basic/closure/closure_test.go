package closure

import (
	"log"
	"sync"
	"testing"
)

func TestClosure(t *testing.T) {
	var aa int

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		aa = 10
		wg.Done()
	}()

	wg.Wait()

	log.Printf("value after called: %v", aa)
}
