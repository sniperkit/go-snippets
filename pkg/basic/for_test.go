package basic

import (
	"log"
	"testing"
)

func TestForRange(t *testing.T) {

	m := map[string]int{
		"NewYork":  1,
		"San Jose": 2,
		"Seattle":  3,
	}

	for range m {
		log.Println("hello")
	}
}

func TestForBreak(t *testing.T) {
loop:
	for i := 0; i < 10; i++ {
		select {
		default:
			log.Println(i)
			break loop
		}
	}
}

func TestForContinue(t *testing.T) {
	for i := 0; i < 10; i++ {
		select {
		default:
			log.Println(i)
			continue
		}
	}
}

func TestForAssignment(t *testing.T) {
	var fnList []func()
	for k := 0; k < 10; k++ {
		tmp := k
		fn := func() {
			log.Println(tmp)
		}
		fnList = append(fnList, fn)
	}

	for _, f := range fnList {
		f()
	}
}
