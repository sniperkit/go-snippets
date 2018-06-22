package basic

import (
	"log"
	"testing"
)

func TestVariableScope(t *testing.T) {

	var aa int
	{
		// `aa` has been shadowed
		aa, bb := 10, 20
		log.Printf("aa: %v, bb: %v", aa, bb)
	}

	log.Printf("aa: %v", aa)
}
