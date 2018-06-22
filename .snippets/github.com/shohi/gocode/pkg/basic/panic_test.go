package basic

// ref, https://medium.com/@thedevsaddam/go-101-defer-panic-and-recover-65a40ee7dcb4

import (
	"log"
	"testing"
)

func TestPanicRecover(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	panic("unable to run program")
}
