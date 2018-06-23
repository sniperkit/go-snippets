package main

import (
	"fmt"

	"github.com/asaskevich/EventBus"
)

const EVENT_1 = "demo:event-1"
const EVENT_2 = "demo:event-2"

// Demonstrate the use of event bus in Go
// using asaskevich/EventBus package
func main() {
	bus := EventBus.New()

	registerListeners(bus)

	// Trigger "event 1", no args
	bus.Publish(EVENT_1)

	// Trigger "event 2", with args
	bus.Publish(EVENT_2, "Hello", "World")

	// Remove (unsubsribe) listener from event
	bus.Unsubscribe(EVENT_1, eventOneListener1)
	fmt.Println("\nListener has been unsubsribed")

	bus.Publish(EVENT_1)

	// spew.Dump(bus)
}

func registerListeners(bus EventBus.Bus) {
	bus.Subscribe(EVENT_1, eventOneListener1)
	bus.Subscribe(EVENT_1, eventOneListener2)
	bus.Subscribe(EVENT_2, eventTwoListener)
}

func eventOneListener1() {
	fmt.Println("Listener 1 executed")
}

func eventOneListener2() {
	fmt.Println("Listener 2 executed")
}

func eventTwoListener(args ...interface{}) {
	fmt.Printf("Listener 3 executed. Args: %v", args)
}
