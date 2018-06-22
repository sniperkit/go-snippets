package os

import (
	"log"
	"os"
	"testing"
)

func TestSetEnv(t *testing.T) {
	key := "_S_HELLO"
	os.Setenv(key, "world")
	log.Println(os.Getenv(key))
	os.Unsetenv(key)
	log.Println(os.Getenv(key))
}

func TestGetEnv(t *testing.T) {
	key := "_S_HELLO"
	log.Println(os.Getenv(key))
}
