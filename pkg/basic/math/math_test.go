package math

import (
	"log"
	"testing"
)

func TestIntDivide(t *testing.T) {
	const num = 256
	log.Printf("%T -- %v", num, (num+31)/32)
}
