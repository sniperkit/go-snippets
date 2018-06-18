package lechan

import (
	"time"

	"sync/atomic"
)

type Lechan struct {
	size   int32
	items  int32
	queue  chan interface{}
	input  chan interface{}
	output chan interface{}
	close  chan struct{}
}

func New(size int) *Lechan {
	ch := &Lechan{
		size:  int32(size),
		queue: make(chan interface{}, size),
		close: make(chan struct{}, 1),
	}
	go ch.run()
	return ch
}

func (ch *Lechan) Put(value interface{}) {
	ch.queue <- value
	atomic.AddInt32(&ch.items, 1)
}

func (ch *Lechan) TryPut(value interface{}) bool {
	if atomic.LoadInt32(&ch.items) == ch.size {
		return false
	}
	ch.queue <- value
	atomic.AddInt32(&ch.items, 1)
	return true
}

func (ch *Lechan) TryCountedPut(value interface{}, count int) bool {
	for i := 0; i < count; i++ {
		if ch.TryPut(value) {
			return true
		}
	}
	return false
}

func (ch *Lechan) TryTimedPut(value interface{}, duration time.Duration) bool {
	for {
		select {
		case <-time.After(duration):
			return false
		default:
			if ch.TryPut(value) {
				return true
			}
		}
	}
}

func (ch *Lechan) TryLimitedPut(value interface{}, count int, delay time.Duration) bool {
	for i := 0; i < count; i++ {
		if ch.TryTimedPut(value, delay) {
			return true
		}
	}
	return false
}

func (ch *Lechan) Get() (value interface{}, ok bool) {
	if atomic.LoadInt32(&ch.items) == 0 {
		return nil, false
	}
	value = <-ch.queue
	atomic.AddInt32(&ch.items, -1)
	return value, true
}

func (ch *Lechan) In() chan<- interface{} {
	return ch.input
}

func (ch *Lechan) Out() <-chan interface{} {
	return ch.output
}

func (ch *Lechan) Close() {
	ch.close <- struct{}{}
}

func (ch *Lechan) run() {
	var next interface{}

	for {
		select {
		case ch.output <- next:

		default:
			select {
			case next = <-ch.queue:
				atomic.AddInt32(&ch.items, -1)

			default:
				select {
				case value := <-ch.input:
					ch.TryPut(value)

				case <-ch.close:
					close(ch.input)
					close(ch.output)
					close(ch.queue)
					close(ch.close)
					return
				}
			}
		}
	}
}
