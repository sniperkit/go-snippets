package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Plugin foo
type Plugin interface {
	DoWork(ctx context.Context, event Event, id int, resultChan chan interface{})
}

// Global plugin store
var plugins = []Plugin{}

// Event for processing
type Event struct {
	Name string
	Data map[string]string
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	// Install our plugins
	plugins = append(plugins, &MathPlugin{})
	plugins = append(plugins, &NumberPlugin{})
	// etc...

	eventChan := make(chan Event)
	done := make(chan struct{})

	//
	// Simulate response of polling worker feeding us events
	//
	go func() {
		for i := 0; i < 5; i++ {
			eventChan <- Event{Name: fmt.Sprintf("Event %d", i)}
			// eventChan <- fmt.Sprintf("Event %d", i)
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
		case e := <-eventChan:

			fmt.Println(e.Name, "just arrived")

			// Stop previous work(ers)
			if cancel != nil {
				cancel()
			}

			// Start a new event context so we can close it if needed
			ctx, cancel = context.WithCancel(context.Background())

			// Give Boss orders
			go boss(ctx, e)

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

func boss(ctx context.Context, e Event) {

	resultChan := make(chan interface{})
	defer close(resultChan)

	// Workers head off (and promise to watch their phones)
	for id, p := range plugins {
		go p.DoWork(ctx, e, id, resultChan)
	}

	// go worker(ctx, e, 1, resultChan)
	// go worker(ctx, e, 2, resultChan)

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
			fmt.Printf("\t(RESULT) %s: %v\n", e.Name, r)
		}
	}

}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
