package testing

import "testing"

func TestCheckingChannel(t *testing.T) {
	stop := make(chan bool)

	func (stop chan bool) {
		close(stop)
	}(stop)

	t.Log(1)
	_, ok := (<-stop)

	if ok {
		t.Error("Channel is not closed")
	}
}
