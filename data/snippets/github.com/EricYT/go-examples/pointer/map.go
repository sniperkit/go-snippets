package main

import "fmt"

func modify(m map[string]string) {
	m["foo"] = "bar"
}

func main() {
	// m just a reference point to a map like slice
	//var m = map[string]string{} // initialize the map
	var m map[string]string
	//	m = make(map[string]string)

	m["key"] = "value"
	fmt.Printf("map:%+v\n", m)
	modify(m)
	fmt.Printf("map:%+v\n", m)
}
