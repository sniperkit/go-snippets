package main

import (
	"fmt"
	"sort"

	"github.com/gophergala/correct-horse-battery-staple/common"
)

func ExampleServerUpdateClientsSort() {
	var msg = common.ServerUpdate{
		Clients: common.ClientStates{
			{Id: 2, Name: "Second"},
			{Id: 1, Name: "First"},
		},
	}

	for _, cs := range msg.Clients {
		fmt.Print(cs.Id, " ", cs.Name, ",")
	}
	fmt.Println()

	sort.Sort(msg.Clients)

	for _, cs := range msg.Clients {
		fmt.Print(cs.Id, " ", cs.Name, ",")
	}
	fmt.Println()

	// Output:
	// 2 Second,1 First,
	// 1 First,2 Second,
}
