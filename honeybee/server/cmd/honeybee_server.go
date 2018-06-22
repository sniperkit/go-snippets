package main

import (
	"fmt"
	"github.com/gophergala/honeybee/server"
)

//TODO save data in DB

func main() {
	nodes := make(map[string]string)
	edges := make(map[string][]string)
	fmt.Println("Starting honeybee server")
	go server.ListenPB(&nodes, &edges)
	server.Run(&nodes, &edges)
}
