package basic

import (
	"log"
	"testing"
)

func TestFunctionWithNil(t *testing.T) {
	var test func(func())
	test = func(f func()) {
		log.Printf("%T", f)
		log.Println(f == nil)
	}

	test(nil)
	test(func() { log.Println("hello") })
}
