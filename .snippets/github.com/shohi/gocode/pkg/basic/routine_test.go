package basic

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestGoroutine(t *testing.T) {
	ch := make(chan interface{})
	go func() {
		ch <- nil
	}()

	info := <-ch
	fmt.Println(info)
}

func TestGoroutineSize(t *testing.T) {
	size := int(1e2)
	ch := make(chan struct{})

	var m0 runtime.MemStats
	var m1 runtime.MemStats
	runtime.GOMAXPROCS(1)

	// 1. before
	runtime.ReadMemStats(&m0)
	log.Println("before ==> ", m0.Sys/1024)
	for i := 0; i < size; i++ {
		go func() {
			ch <- struct{}{}
		}()
	}

	runtime.Gosched()
	runtime.GC()

	// 2. after
	runtime.ReadMemStats(&m1)
	log.Println("after ==> ", m1.Sys/1024)
	log.Println("avg ==> ", float64(m1.Sys-m0.Sys)/float64(size))
	log.Println("goroutines ==> ", runtime.NumGoroutine())

}

// ref, https://stackoverflow.com/questions/8509152/max-number-of-goroutines
func TestGoroutineSizeUsingStackoverflow(t *testing.T) {
	number := int(1e5)
	ch := make(chan byte)
	counter := 0

	f := func() {
		counter++
		<-ch
	}

	// limit the number of spare OS threads to just 1
	runtime.GOMAXPROCS(1)

	// Make a copy of MemStats
	var m0 runtime.MemStats
	runtime.ReadMemStats(&m0)

	t0 := time.Now().UnixNano()
	for i := 0; i < number; i++ {
		go f()
	}
	runtime.Gosched()
	t1 := time.Now().UnixNano()
	runtime.GC()

	// Make a copy of MemStats
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	if counter != number {
		fmt.Fprintf(os.Stderr, "failed to begin execution of all goroutines")
		os.Exit(1)
	}

	log.Printf("Number of goroutines: %d\n", number)
	log.Printf("Per goroutine:\n")
	log.Printf("  Memory: %.2f bytes\n", float64(m1.Sys-m0.Sys)/float64(number))
	log.Printf("  Time:   %f µs\n", float64(t1-t0)/float64(number)/1e3)

}
