package test

import (
	"fmt"
	"testing"
	"goworkshop/model"
)

func TestMaps(t *testing.T) {

	//1. How to declare?
	var ages map[string]int // mapping from strings to ints
	fmt.Printf("1. ages = %v\n", ages)

	//2. How to initialize?
	ages = make(map[string]int)
	// or
	ages = map[string]int{
		"alice": 31,
		"bob":   34,
	}
	fmt.Printf("2. ages = %v\n", ages)

	// 3. How to get the size?
	fmt.Printf("3. ages size = %v\n", len(ages))

	// 4. How to get an element
	fmt.Printf("4. ages[alice] = %d\n", ages["alice"])

	// 5. How to add an element?
	ages["carol"] = 30
	fmt.Printf("ages[carol] = %d\n", ages["carol"])

	// 6. How to traverse?
	fmt.Printf("6. traversing a map:\n")
	for name, age := range ages {
		fmt.Printf("\t ages[%s] = %d\n", name, age)
	}

	// 7. How to remove an element?
	delete(ages, "bob") // remove element ages["bob"]
	fmt.Printf("7. ages = %v\n", ages)

	authorsMap := make(map[string]model.Author)
	authorsMap["1"] = model.Author{
		UUID: "1",
	}

	if authorsMap["2"] == (model.Author{}){
		fmt.Println("There is no author with uuid 2")
	}
}
