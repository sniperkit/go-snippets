package basic

import (
	"log"
	"os"
	"testing"
)

func setup() {
	log.Println("setup...")
}

func teardown() {
	log.Println("teardown...")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestHello(t *testing.T) {
	log.Println("test hello.")
}

func TestWorld(t *testing.T) {
	log.Println("test world.")
}
