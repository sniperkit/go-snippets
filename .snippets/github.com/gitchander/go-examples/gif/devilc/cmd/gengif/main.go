package main

import (
	"log"

	"github.com/gitchander/go-examples/gif/devilc"
)

func main() {
	if err := devilc.MakeGif("devil-curve.gif"); err != nil {
		log.Fatal(err)
	}
}
