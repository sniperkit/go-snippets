package main

import (
	"fmt"
	"time"
)

type Person struct {
	Name    string
	Surname string
	Age     int
}

func timeMap(my_interface interface{}) {
	updated_interface, status := my_interface.(map[string]interface{})
	if status {
		updated_interface["updated_at"] = time.Now()
	}
}

func main() {
	person := Person{Name: "Nexus", Surname: "Corpse", Age: 22}
	my_map := map[string]interface{}{"Nexus": person}

	timeMap(my_map)
	fmt.Println(my_map)
}
