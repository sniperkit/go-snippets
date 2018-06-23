package pubsub

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"
)

const (
	topic Topic = "test-topic"
)

func waitForMessagehandlingToBeComplete(done <-chan struct{}, t *testing.T) {
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("publish did not complete")
	}
}

func TestPublish(t *testing.T) {
	hub := NewSimplehub()
	done := hub.Publish(topic, nil)
	waitForMessagehandlingToBeComplete(done, t)
}

func TestSubscribe(t *testing.T) {
	var called bool
	hub := NewSimplehub()
	hub.Subscribe(topic, func(topic_ Topic, data interface{}) {
		if topic != topic_ {
			t.Fatalf("topic shuold equal %s", topic)
		}
		if data != nil {
			t.Fatalf("handler shuold receive a nil data")
		}
		called = true
	})
	done := hub.Publish(topic, nil)
	waitForMessagehandlingToBeComplete(done, t)
	if !called {
		t.Fatalf("publish handle not be called")
	}
}

func TestPublishCompleterWait(t *testing.T) {
	wait := make(chan struct{})
	hub := NewSimplehub()
	hub.Subscribe(topic, func(topic Topic, data interface{}) {
		<-wait
	})
	done := hub.Publish(topic, nil)

	select {
	case <-done:
		t.Fatalf("publish didn't wait")
	case <-time.After(time.Millisecond):
	}
	close(wait)
	waitForMessagehandlingToBeComplete(done, t)
}

func TestSubscriberExecInOrder(t *testing.T) {
	mutex := sync.Mutex{}
	hub := NewSimplehub()
	var calls []Topic
	hub.Subscribe(MatchRegexp("test.*"), func(topic Topic, data interface{}) {
		mutex.Lock()
		defer mutex.Unlock()
		calls = append(calls, topic)
	})

	var lastdone <-chan struct{}
	var rightCalls []Topic
	for index := 0; index < 10; index++ {
		topc := Topic(fmt.Sprintf("test.%d", index))
		rightCalls = append(rightCalls, topc)
		lastdone = hub.Publish(topc, nil)
	}
	waitForMessagehandlingToBeComplete(lastdone, t)
	if !reflect.DeepEqual(calls, rightCalls) {
		t.Fatalf("%#v != %#v", calls, rightCalls)
	}
}

func TestPublishNotBlockedByHandleFunc(t *testing.T) {
	wait := make(chan struct{})
	hub := NewSimplehub()
	hub.Subscribe(topic, func(t Topic, d interface{}) {
		<-wait
		log.Printf("subscriber receive data: %+v", d)
	})

	var lastdone <-chan struct{}
	for i := 0; i < 5; i++ {
		lastdone = hub.Publish(topic, i)
	}
	close(wait)
	waitForMessagehandlingToBeComplete(lastdone, t)
}
