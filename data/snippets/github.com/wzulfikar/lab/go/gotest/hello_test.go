// file: yo_test.go
// run test: go test .
package hello

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	fmt.Println("test 1")
	SayHi() // ← undefined
}

func Test2(t *testing.T) {
	fmt.Println("test 2")
	SayHi() // ← undefined
}
