package queue

import (
	"fmt"
	"time"
)

// Option is anything that can be passed to Publish and Subscribe
type Option interface{}

// Model helps specific implementations rehydrate messages from wire format
type Model interface{}

// Handler receives payload from the queue
type Handler func(msg interface{}) error

// Middleware applies transform to Handler
type Middleware func(Handler) Handler

// Delay is an option to defer Publish. Also can be returned by
// handlers to requeue the same data with delay.
type Delay time.Duration

func (e Delay) Error() string {
	return fmt.Sprintf("requeue after %v", e)
}

// Concurrency is passed to Subscribe and specifies level of parallelism.
type Concurrency int

// Publisher sends objects to topic
type Publisher interface {
	Publish(topic string, obj interface{}, options ...Option) error
}

// Subscriber connects Handler to topic through channel
type Subscriber interface {
	Subscribe(topic, channel string, h Handler, options ...Option) error
}

//go:generate moq -out mock/queue.go -pkg mock . Publisher Subscriber
