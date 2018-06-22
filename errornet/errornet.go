package errornet

import (
	"time"
)

type ErrorNet struct {
	// Handlers is a slice of Handler. All of which will be called during a panic or error.
	Handlers []Handler

	// errors is a channel of Error (panic and error values) in queue to be reported on.
	errors chan Error

	// exit is a channel used for cleanly exiting an application.
	exit chan bool
}

func (en *ErrorNet) AddHandler(h Handler) {
	en.Handlers = append(en.Handlers, h)
}

func (en *ErrorNet) Report(v interface{}) {
	// start a go routine so Report is never blocking
	go func() {
		// convert the incoming error or panic into errornet.Error
		err := Error{}

		// move the error to the queue
		en.errors <- err
	}()
}

func (en *ErrorNet) Exit(max time.Duration) {
}

func (en *ErrorNet) router() {
	select {
	case _, ok := <-en.exit:
		// try to cleanly exit

	case err := <-en.errors:
		// send error to each handler
		for _, h := range en.Handlers {
			if !h.Failover() {
				h.Handle(err)
			}
		}
	}
}

func New(dir string) *ErrorNet {
	en := ErrorNet{errors: make(chan Error), exit: make(chan bool)}

	go en.router()

	return en
}

type Handler interface {
	Handle(Error) error
	Failover() bool
}

type Error struct {
	Host     string
	Function string
	File     string
	Line     int
	Panic    bool
	Time     time.Time
}
