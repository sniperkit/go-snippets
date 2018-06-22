package function

import (
	"log"
	"testing"
)

func Hello() {

}

func TestFunction(t *testing.T) {
	h := func() {}
	log.Printf("%p ==> %p", Hello, h)
}
