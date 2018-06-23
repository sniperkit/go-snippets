package device

import "github.com/gophergala/stevek/event"

var Registry = make(map[string]Device)

type config string

// should typedef the filter function signature
type Device struct {
	Name      string
	filterFn  func(message, config) bool
	Transform func(message) message
	config    config
	Transit   <-chan message
}

func New(name, conf string) (dev Device) {
	// Should perform some error checking
	dev = Registry[name]
	dev.config = config(conf)
	return
}

func (d Device) Filter(in <-chan message) <-chan message {
	out := make(chan message)
	go func() {
		for msg := range in {
			msg.Device = d.Name
			if d.filterFn(msg, d.config) {
				msg.Action = "Permitted"
			} else {
				msg.Ghost = true
				msg.Action = "Denied"
			}
			out <- msg
			event.Events <- event.Event{Data: msg, Type: "Transit event"}
		}
		close(out)
	}()
	return out
}
