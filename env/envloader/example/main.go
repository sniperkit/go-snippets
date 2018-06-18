package main

import (
	"log"
	"os"
)

func main() {
	log.Println("This is on ex")
	log.Println(os.Getenv("SS_SOMETHING"))
}
