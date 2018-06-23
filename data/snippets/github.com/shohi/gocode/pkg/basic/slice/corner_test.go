package slice

import (
	"log"
	"strings"
	"testing"
)

func TestInitialization(t *testing.T) {
	s := make([]int, 10)

	log.Println(strings.Repeat("*", 10)+"Before ===>", len(s))
	for k, v := range s {
		log.Printf("%d ===> %v", k, v)
	}

	//
	s = append(s, 10)
	log.Println(strings.Repeat("*", 10)+"after ===> ", len(s))
	for k, v := range s {
		log.Printf("%d ===> %v", k, v)
	}
}

func TestInitWithNil(t *testing.T) {
	var bs []byte
	bs = nil
	log.Printf("byte slice: %v", bs)
}
