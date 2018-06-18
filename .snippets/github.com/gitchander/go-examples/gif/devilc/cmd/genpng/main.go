package main

import (
	"log"

	"github.com/gitchander/go-examples/gif/devilc"
)

func main() {
	if err := devilc.MakePNG("test.png", 512, 512, 90, 100); err != nil {
		log.Fatal(err)
	}
}
