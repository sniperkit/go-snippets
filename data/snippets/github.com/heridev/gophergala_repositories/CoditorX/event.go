// Copyright (c) 2015, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

const (
	// EvtCodeServerInternalError indicates an event: server internal error
	EvtCodeServerInternalError = iota
)

// Max length of queue.
const maxQueueLength = 10

// event represents an event.
type event struct {
	Code int         `json:"code"` // event code
	Sid  string      `json:"sid"`  // Coditor session id related
	Data interface{} `json:"data"` // event data
}

// Global event queue.
//
// Every event in this queue will be dispatched to each user event queue.
var EventQueue = make(chan *event, maxQueueLength)

// UserEventQueue represents a user event queue.
type userEventQueue struct {
	Sid      string         // Coditor session id related
	Queue    chan *event    // queue
	Handlers []eventHandler // event handlers
}

type queues map[string]*userEventQueue

// User event queues.
//
// <sid, *UserEventQueue>
var userEventQueues = queues{}

// Load initializes the event handling.
func Load() {
	go func() {
		for event := range EventQueue {
			logger.Debugf("Received a global event [code=%d]", event.Code)

			// dispatch the event to each user event queue
			for _, userQueue := range userEventQueues {
				event.Sid = userQueue.Sid

				userQueue.Queue <- event
			}
		}
	}()
}

// AddHandler adds the specified handlers to user event queues.
func (uq *userEventQueue) addHandler(handlers ...eventHandler) {
	for _, handler := range handlers {
		uq.Handlers = append(uq.Handlers, handler)
	}
}

// New initializes a user event queue with the specified Coditor session id.
func (ueqs queues) new(sid string) *userEventQueue {
	q := ueqs[sid]
	if nil != q {
		logger.Warnf("Already exist a user queue in session [%s]", sid)

		return q
	}

	q = &userEventQueue{
		Sid:   sid,
		Queue: make(chan *event, maxQueueLength),
	}

	ueqs[sid] = q

	go func() { // start listening
		for evt := range q.Queue {
			logger.Debugf("Session [%s] received an event [%d]", sid, evt.Code)

			// process event by each handlers
			for _, handler := range q.Handlers {
				handler.handle(evt)
			}
		}
	}()

	return q
}

// Close closes a user event queue with the specified Coditor session id.
func (ueqs queues) close(sid string) {
	q := ueqs[sid]
	if nil == q {
		return
	}

	delete(ueqs, sid)
}

// eventHandler represents an event handler.
type eventHandler interface {
	handle(evt *event)
}

// eventHandleFunc represents a handler function.
type eventHandleFunc func(evt *event)

// Default implementation of event handling.
func (fn eventHandleFunc) handle(evt *event) {
	fn(evt)
}
