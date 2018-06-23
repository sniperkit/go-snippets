package world

import (
	"log"
	"testing"
)

// SayWorld
func SayWorld() {
	log.Println("world")
	// can not hello.SayHello()
}

func TestSayWorld(t *testing.T) {
	SayWorld()
}
