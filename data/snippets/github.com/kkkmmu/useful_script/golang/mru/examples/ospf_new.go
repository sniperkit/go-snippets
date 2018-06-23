package main

import (
	"log"
	"ospf"
)

func main() {
	ins, err := ospf.New("ens33", "10.1.1.1", "0.0.0.1", 0)
	if err != nil {
		log.Fatal(err)
	}

	ins.Start()

	for {
	}
}
