package basic

import (
	"log"
	"testing"
)

// `defer` gotcha sample code
// mainly ref, https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01

func TestDeferGotchaNilFunc(t *testing.T) {
	var run func() = nil

	// run nil func will arise panic
	// use another defer to recover
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer run()

	log.Println("runs")
}
