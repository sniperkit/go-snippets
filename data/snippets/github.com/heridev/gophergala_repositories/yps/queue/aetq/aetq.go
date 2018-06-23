// Package aetq implements gophergala/yps/queue using AppEngine TaskQueue functionality
package aetq

import (
	"appengine"
	"appengine/taskqueue"

	"github.com/gophergala/yps/queue"
)

type (
	msg struct {
		original *taskqueue.Task
	}

	mq struct {
		name  string
		ctx   appengine.Context
		lease int
	}
)

func (m *msg) Original() interface{} {
	return m.original
}

func (m *msg) String() string {
	if m.original == nil {
		return ""
	}

	return string(m.original.Payload)
}

// NewMessage creates a wrapper message of queue.Message interface over an taskqueue.Task
func NewMessage(payload interface{}) queue.Message {
	var ms *taskqueue.Task
	if _, ok := payload.(*taskqueue.Task); ok {
		ms = payload.(*taskqueue.Task)
	} else {
		switch payload.(type) {
		case string:
			{
				ms = &taskqueue.Task{
					Payload: []byte(payload.(string)),
					Method:  "PULL",
				}
			}

		case []byte:
			{
				ms = &taskqueue.Task{
					Payload: payload.([]byte),
					Method:  "PULL",
				}
			}

		default:
			{
				return nil
			}
		}

	}

	return &msg{
		original: ms,
	}
}

func (q *mq) Add(message *queue.Message) (err error) {
	m := (*message).Original().(*taskqueue.Task)
	_, err = taskqueue.Add(q.ctx, m, q.name)
	return
}

func (q *mq) Fetch(count int) (messages []*queue.Message, err error) {
	var msgs []*taskqueue.Task

	if msgs, err = taskqueue.Lease(q.ctx, count, q.name, q.lease); err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		m := NewMessage(msg)
		messages = append(messages, &m)
	}

	return
}

func (q *mq) Confirm(message *queue.Message) error {
	return q.Delete(message)
}

func (q *mq) Delete(message *queue.Message) error {
	m := (*message).Original().(*taskqueue.Task)
	return taskqueue.Delete(q.ctx, m, q.name)
}

// NewQueue returns a new queue.Queue implementation using AppEngine TaskQueue as a backend
func NewQueue(context appengine.Context, queueName string, leaseTime int) queue.Queue {
	return &mq{
		name:  queueName,
		ctx:   context,
		lease: leaseTime,
	}
}
