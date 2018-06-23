package basic

import (
	"log"
	"testing"
)

func TestSyntaxAssignment(t *testing.T) {
	a := 20
	if 20 > 10 {
		a, b := 10, 20
		log.Println(a, b)
	}

	log.Println(a)

}

func TestSyntaxReassignment(t *testing.T) {

	a := 20
	a, b := 15, 20

	if b > a {
		log.Println(a, b)
	}

	log.Println(a, b)
}
