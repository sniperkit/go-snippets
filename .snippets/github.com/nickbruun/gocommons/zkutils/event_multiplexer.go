package zkutils

import (
	"github.com/samuel/go-zookeeper/zk"
	"sync"
)

var (
	eventMultiplexerSubscriberBuffer = 16
)

// Event multiplexer.
type EventMultiplexer struct {
	in   <-chan zk.Event
	outs []chan zk.Event
	lock sync.Mutex
}

// New event multiplexer.
func NewEventMultiplexer(eventChan <-chan zk.Event) *EventMultiplexer {
	m := &EventMultiplexer{
		in:   eventChan,
		outs: make([]chan zk.Event, 0),
	}

	go func() {
		for ev := range m.in {
			m.lock.Lock()

			remove := make([]int, 0)

			for i, out := range m.outs {
				select {
				case out <- ev:
				default:
					remove = append(remove, i)
				}
			}

			for j := len(remove) - 1; j >= 0; j-- {
				i := remove[j]
				m.outs = append(m.outs[:i], m.outs[i+1:]...)
			}

			m.lock.Unlock()
		}

		m.lock.Lock()
		for _, out := range m.outs {
			close(out)
		}
		m.lock.Unlock()
	}()

	return m
}

// Subscribe to events.
func (m *EventMultiplexer) Subscribe() <-chan zk.Event {
	ec := make(chan zk.Event, eventMultiplexerSubscriberBuffer)

	m.lock.Lock()
	m.outs = append(m.outs, ec)
	m.lock.Unlock()

	return ec
}
