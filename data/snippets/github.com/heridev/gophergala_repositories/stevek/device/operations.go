package device

import (
	"time"

	"github.com/gophergala/stevek/event"
)

// As per http://blog.golang.org/pipelines
// Takes any number of ints (variadic fn) and puts them onto a channel
func Gen(nums ...int) <-chan message {
	out := make(chan message)
	go func() {
		time.Sleep(10 * time.Second)
		for _, n := range nums {
			var msg message
			msg.Device = "Source" // TODO - fix, shouldn't be hardcoded, make like Filter
			msg.Port = n
			msg.Action = "Generated"
			out <- msg
			event.Events <- event.Event{Data: msg, Type: "Transit event"}
		}
		close(out)
	}()
	return out
}

func Identity(msg message) message {
	return msg
}
func AllowAll(msg message, c config) bool {
	return true
}
