package channel

import (
	"log"
	"testing"
)

func TestCloseChannel(t *testing.T) {
	ch := make(chan error, 10)
	close(ch)
	select {
	case err, ok := <-ch:
		log.Println(err, ok)
	}

}
