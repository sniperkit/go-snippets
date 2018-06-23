package basic

import (
	"log"
	"testing"
)

func TestAppendWithInitCap(t *testing.T) {
	b := make([]int, 1024)
	b = append(b, 99)
	log.Println("len:", len(b), "cap:", cap(b))
}
