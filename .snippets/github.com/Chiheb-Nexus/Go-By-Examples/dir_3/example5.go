package main

import (
	"fmt"
)

type Position struct {
	Lat, Long float64
}

func main() {
	location := make(map[string]Position)
	location["Bell Lab"] = Position{40.68433, -74.39967}
	fmt.Printf("%v\n", location)
}
