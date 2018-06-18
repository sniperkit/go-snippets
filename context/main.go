package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	for i := 0; i <= 3; i++ {

		// Create some work
		ctx, cancel := context.WithCancel(context.Background())

		// Worker heads off and promises to watch his pager
		go worker(ctx, i)

		// Pretent we are waiting for another event
		time.Sleep(time.Second)

		// Cancel this workload!
		cancel()
	}

}

func worker(ctx context.Context, i int) {

	var num int

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%d: Made it to %d when cancel() was called\n", i, num)
			return
		default:
			// Do some work
			num++
		}
	}

}
