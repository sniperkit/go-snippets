package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

/*
 * In the following application we have:
 * 1. a main process that gets events and instructs the Boss
 * 2. a Boss that managers the workers until:
 * 		A) Work is done
 * 		B) Process gives Boss new work
 * 3. workers that work until:
 * 		A) They are Finished
 * 		B) Boss tells them quit
 */

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	eventChan := make(chan string)
	done := make(chan struct{})

	//
	// Simulate response of polling worker feeding us events
	//
	go func() {
		for i := 0; i < 5; i++ {
			eventChan <- fmt.Sprintf("Event %d", i)
			time.Sleep(time.Second * time.Duration(randInt(1, 2)))
		}
	}()

	//
	// Simulate program exit or some failure
	//
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	//
	// Main Loop
	//
	var ctx context.Context
	var cancel context.CancelFunc

main:
	for {
		select {
		case event := <-eventChan:

			fmt.Println(event, "just arrived")

			// Stop previous work(ers)
			if cancel != nil {
				cancel()
			}

			// Start a new event context so we can close it if needed
			ctx, cancel = context.WithCancel(context.Background())

			// Give Boss orders
			go boss(ctx, event)

			// Since eventChan is non-buffered we loop back top and wait
			// while boss and workers are running

		case <-done:
			if cancel != nil {
				cancel()
			}
			break main
		case <-time.After(5 * time.Second):
			panic("System must have died, no events for a while...")
		}
	}

}

func boss(ctx context.Context, event string) {

	resultChan := make(chan interface{})
	defer close(resultChan)

	// Workers head off (and promise to watch their phones)
	go worker(ctx, event, 1, resultChan)
	go worker(ctx, event, 2, resultChan)

	// If the workers had enough time to finish
	// there will be something here
mainresult:
	for {
		select {

		// Ran out of time to process stuff
		case <-ctx.Done():
		// If we loose a connection or something and stop getting events this fires
		case <-time.After(time.Second * 5):
			break mainresult

		// Process worker result
		case r := <-resultChan:
			fmt.Printf("\t(RESULT) %s: %v\n", event, r)
		}
	}

}

// Does something important
func worker(ctx context.Context, event string, id int, resultChan chan interface{}) {

	var num int

	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\t%s (worker %d) only made it to %d when event cancel() was called\n", event, id, num)
			return
		default:
			// Do some work
			num++
			time.Sleep(time.Millisecond * time.Duration(randInt(50, 150)))
		}
	}

	// Finished
	resultChan <- fmt.Sprintf("worker-%d = %d", id, num)

}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
