package main

import (
	"fmt"
)

type vehicle interface {
	engine() string
	wheels() int
	isElectric() bool
	displayInfo()
}

type car struct {
	color string
	make  string
	model string
}

func (c car) engine() string {
	return "standard engine"
}

func (c car) wheels() int {
	return 4
}

func (c car) isElectric() bool {
	return false
}

func (c car) displayInfo() {
	fmt.Println(" color: " + c.color + ", make: " + c.make + ", model: " + c.model)
}

type truck struct {
	color  string
	make   string
	pickup bool
}

func vehicleInfo(v vehicle) {
	v.displayInfo()
}

func main() {
	toyota := car{color: "Red", make: "Toyota", model: "2017"}
	car.displayInfo(toyota)
}
