// Package nsq implements Publisher and Subscriber using NSQ library and JSON
// serialization
package nsq

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nsqio/go-nsq"

	"github.com/go-mixins/queue"
)

// Publisher implements queue.Publisher with NSQ
type Publisher struct {
	producer *nsq.Producer
	config
}

var _ queue.Publisher = (*Publisher)(nil)

// NewPublisher creates the Publisher
func NewPublisher(opts ...Option) (res *Publisher, err error) {
	res = new(Publisher)
	res.config = config{marshal: json.Marshal, options: nsq.NewConfig()}
	for _, opt := range opts {
		if err = opt(&res.config); err != nil {
			return
		}
	}
	if n := len(res.nsqds); n == 0 {
		res.nsqds = append(res.nsqds, "localhost:4150")
	} else if n > 1 {
		err = errors.New("only one queue allowed for Publisher")
		return
	}
	if res.producer, err = nsq.NewProducer(res.nsqds[0], res.options); err != nil {
		return
	}
	if res.logger != nil {
		res.producer.SetLogger(res, res.level)
	}
	err = res.producer.Ping()
	return
}

// Publish marshals object to JSON and sends it to the specified topic
func (p *Publisher) Publish(topic string, obj interface{}, options ...queue.Option) (err error) {
	data, err := p.marshal(obj)
	if err != nil {
		return
	}
	pf := p.producer.Publish
	for _, opt := range options {
		switch t := opt.(type) {
		case queue.Delay:
			pf = func(topic string, data []byte) error {
				return p.producer.DeferredPublish(topic, time.Duration(t), data)
			}
		}
	}
	return pf(topic, data)
}
