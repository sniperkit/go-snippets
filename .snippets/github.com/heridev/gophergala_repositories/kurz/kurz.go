package main

import (
	"fmt"
	"github.com/FGM/kurz/storage"
	"github.com/FGM/kurz/strategy"
	"log"
)

func main() {
	err := storage.Service.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Service.Close()

	fmt.Print(strategy.Statistics.Refresh(storage.Service))
}
