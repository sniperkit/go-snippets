package util

import (
	"fmt"
	"testing"
)

func TestGoID(t *testing.T) {
	fmt.Printf("Goroutine ID ==> %d\n", GoID())
}
