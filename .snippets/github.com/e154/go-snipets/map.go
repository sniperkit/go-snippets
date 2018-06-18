package main

import "fmt"

type Vertex struct {
	Lat, Long float64
}


var m map[string]Vertex
var	n = map[string]Vertex{
		"Bell Labs": Vertex {40.68433, -74.39967,},
		"Google": Vertex {37.42202, -122.08408,},
	}

var	o = map[string]Vertex{
	"Bell Labs": {40.68433, -74.39967,},
	"Google": {37.42202, -122.08408,},
}

func main() {
	m = make(map[string]Vertex)
	m["Bell Labs"] = Vertex{
		40.68433, -74.39967,
	}
	fmt.Println(m["Bell Labs"])

	//////////////////
	//Map literals
	//////////////////
	fmt.Println("")
	fmt.Println("Map literals")
	fmt.Println(n)

	//////////////////
	//Map lierals continued
	//////////////////
	fmt.Println("")
	fmt.Println("Map literals continued")
	fmt.Println(o)

	//////////////////
	//Mutating Maps
	//////////////////
	fmt.Println("")
	fmt.Println("Mutating Maps")
	m := make(map[string]int)

	m["Answer"] = 42
	fmt.Println("The value:", m["Answer"])

	m["Answer"] = 48
	fmt.Println("The value:", m["Answer"])

	delete(m, "Answer")
	fmt.Println("The value:", m["Answer"])

	v, ok := m["Answer"]
	fmt.Println("The value:", v, "Present?", ok)

	l := map[string]string{}
	l["key"] = "value"
	for key, value := range l {
		fmt.Println("The key:", key, "value:", value)
	}

	// How to check if a map contains a key in go?
	if val, ok := l["key2"]; ok {
		fmt.Printf("key exist, %s\r\n", val)
	}

	fmt.Print(`//////////////////
// links
//////////////////
`)

	var m1, m2 map[string]interface{}
	m1 = make(map[string]interface{})
	m2 = make(map[string]interface{})
	m1["m1"] = "m1"
	m2["m2"] = "m2"

	fmt.Println(m1)
	fmt.Println(m2)

	m1 = m2
	m2["m3"] = "m3"
	m1["m4"] = "m4"
	fmt.Println("m1",m1, &m1)
	fmt.Println("m2",m2, &m2)
}
