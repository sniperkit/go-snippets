package main

import "fmt"

type Color struct {
	r, g, b int
}

func main() {
	//	var m map[string]int
	m := make(map[string]int)
	m["k1"] = 1
	m["k2"] = 1
	fmt.Println("map before delete ", m)
	delete(m, "k3")
	fmt.Println("map delete ", m)

	// test
	if value, ok := m["k2"]; ok {
		fmt.Println(value, ok)
	}

	for key := range m {
		fmt.Println("Key:", key)
	}

	if v1, o1 := m["k3"]; o1 {

	} else {
		fmt.Println(v1, o1)
	}

	// map
	//	m = map[string]Color{
	//		"Red":   Color{255, 0, 0},
	//		"Green": Color{0, 255, 0},
	//		"Blue":  Color{0, 0, 255},
	//	}
	//	key, value := m["PeachPuff"]

}
