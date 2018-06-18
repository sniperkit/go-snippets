package main

import (
	"fmt"
)

func main() {
	ford := Car{}
	ford.name = "Ford"
	ford.color = "black"
	ford.start()
	ford.move(40)
	ford.stop()
	fmt.Println("Is engine started?", ford.isStarted)

	ferrari := Car{}
	ferrari.name = "Ferrari"
	ferrari.color = "red"
	ferrari.start()

	stopAnyVehicle(&ford)
	stopAnyVehicle(&ferrari)

	// because `Vehicle` has start and stop method
	// which has pointer,
	// this won't work (see the err msg):
	// mustang := Car{}
	// stopAnyVehicle(mustang)

	// unless you do:
	// mustang := &Car{}
	// stopAnyVehicle(mustang)

	// or:
	// mustang := Car{}
	// stopAnyVehicle(&mustang)
}

type Vehicle interface {
	start()
	move(speed int)
	stop()
}

type Car struct {
	isStarted bool
	name      string
	color     string
}

// what's the diff between `func (c Car)` & `func (c *Car)`?
// `func (c Car)` will pass Car by value
// while `func (c *Car)` will pass it by reference.
// try change receiver of start method below from
// `c *Car` to `c Car`.
// see: https://golang.org/doc/faq#methods_on_values_or_pointers
func (c *Car) start() {
	c.isStarted = true
	fmt.Println("engine is started")
}

func (c Car) move(speed int) {
	if !c.isStarted {
		panic("Engine not started")
	}
	fmt.Println("the", c.name, "car which color is ", c.color, "is moving at", speed, "km/h")
}

func (c *Car) stop() {
	c.isStarted = false
	fmt.Println(c.name, "is stopped")
}

func stopAnyVehicle(v Vehicle) {
	v.stop()
}
