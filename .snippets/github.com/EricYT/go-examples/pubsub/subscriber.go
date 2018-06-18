package pubsub

import (
	"log"
	"sync"

	"github.com/juju/utils/deque"
)

type Unsubscriber interface {
	Unsubscribe()
}

type subscriber struct {
	id int

	topicMatcher TopicMatcher
	handler      func(topic Topic, data interface{})

	mutex   sync.Mutex
	pending *deque.Deque

	data   chan struct{}
	done   chan struct{}
	closed chan struct{}
}

func newSubscriber(topicMatcher TopicMatcher, handler func(Topic, interface{})) *subscriber {
	// A closed channel is used to provide an immediate route through a select call
	// in the loop function
	closed := make(chan struct{})
	close(closed)
	s := &subscriber{
		id:           0,
		topicMatcher: topicMatcher,
		handler:      handler,
		pending:      deque.New(),
		data:         make(chan struct{}, 1),
		done:         make(chan struct{}),
		closed:       closed,
	}
	go s.loop()
	log.Printf("created subscriber %p for %+v", s, topicMatcher)
	return s
}

func (s *subscriber) close() {
	log.Printf("subscriber ready to close %d", s.id)
	s.mutex.Lock()
	s.mutex.Unlock()
	for call, ok := s.pending.PopFront(); ok; call, ok = s.pending.PopFront() {
		call.(*handlerCallback).done()
	}
	close(s.done)
}

func (s *subscriber) Notify(call *handlerCallback) {
	log.Printf("subscriber %d get a notify %+v", s.id, call)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.pending.PushBack(call)
	if s.pending.Len() == 1 {
		s.data <- struct{}{}
	}
}

func (s *subscriber) popOne() (*handlerCallback, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.pending.PopFront()
	if !ok {
		return nil, true
	}
	empty := s.pending.Len() == 0
	return val.(*handlerCallback), empty
}

func (s *subscriber) loop() {
	log.Printf("subscriber %b loop", s.id)
	var next chan struct{}

	for {
		select {
		case <-s.done:
			log.Printf("subscriber %b done", s.id)
			return
		case <-s.data:
			log.Printf("subscriber %b receive a notify", s.id)
		case <-next:
		}

		call, empty := s.popOne()
		if empty {
			next = nil
		} else {
			next = s.closed
		}

		if call != nil {
			log.Printf("subscriber %p exec callback (%d) func %p", s, s.id, s.handler)
			s.handler(call.topic, call.data)
			call.done()
		}
	}
}
