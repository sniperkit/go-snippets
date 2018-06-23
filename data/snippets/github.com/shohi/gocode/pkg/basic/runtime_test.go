package basic

import (
	"log"
	"runtime"
	"sync"
	"testing"
	"time"
)

// ref, https://stackoverflow.com/questions/8509152/max-number-of-goroutines
func TestRuntime(t *testing.T) {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
	}()
	buff := make([]byte, 10000)
	stackSize := runtime.Stack(buff, true)
	log.Println(string(buff[0:stackSize]))

}
