package basic

import (
	"log"
	"testing"
)

func TestBetweenOperator(t *testing.T) {
	a := 10

	if 1 < a && a < 20 {
		log.Println(a)
	}

}

func TestParallelAssign(t *testing.T) {
	a, b := 10, 20

	log.Println(a, b)
}
