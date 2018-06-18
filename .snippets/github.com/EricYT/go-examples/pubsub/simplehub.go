package pubsub

import (
	"log"
	"sync"
)

type simplehub struct {
	mutex       sync.Mutex
	subscribers []*subscriber
	idx         int
}

func NewSimplehub() *simplehub {
	hub := &simplehub{}
	return hub
}

func (s *simplehub) Publish(topic Topic, data interface{}) <-chan struct{} {
	log.Printf("simplehub publish topic: %s data: %+v", topic, data)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	done := make(chan struct{})
	var wg sync.WaitGroup
	for _, subscriber := range s.subscribers {
		if subscriber.topicMatcher.Match(topic) {
			wg.Add(1)
			handler := &handlerCallback{topic, data, &wg}
			subscriber.Notify(handler)
		}
	}
	// wait all subscribers done
	go func() {
		wg.Wait()
		close(done)
	}()

	return done
}

func (s *simplehub) Subscribe(matcher TopicMatcher, handler func(Topic, interface{})) Unsubscriber {
	log.Printf("simplehub subscribe handler %p matcher %+v", handler, matcher)
	if matcher == nil || handler == nil {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	subscriber := newSubscriber(matcher, handler)
	subscriber.id = s.idx
	s.subscribers = append(s.subscribers, subscriber)
	s.idx++
	return &handle{s, subscriber.id}
}

func (s *simplehub) unsubscribe(id int) {
	log.Printf("simplehub unsubscribe id %b", id)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for index, subscriber := range s.subscribers {
		if subscriber.id != id {
			continue
		}
		subscriber.close()
		s.subscribers = append(s.subscribers[0:index], s.subscribers[index+1:]...)
	}
}

type handle struct {
	hub *simplehub
	id  int
}

func (h *handle) Unsubscribe() {
	h.hub.unsubscribe(h.id)
}

type handlerCallback struct {
	topic Topic
	data  interface{}
	wg    *sync.WaitGroup
}

func (h *handlerCallback) done() {
	h.wg.Done()
}
