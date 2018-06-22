package hello

import "testing"

// HelloWorld
func TestHelloWorld(t *testing.T) {
	if hello() != "hello world!!" {
		t.Fatal("test")
	}
}
