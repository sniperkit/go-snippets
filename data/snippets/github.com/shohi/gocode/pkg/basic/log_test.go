package basic

import (
	"log"
	"testing"
)

func TestLogPrint(t *testing.T) {
	log.Print("hello\naaaa")
	log.Print("world")
}

func TestLogPrintf(t *testing.T) {
	log.Printf("test log: %s", "hello")
	log.Printf("test log: %s", "world")
}
