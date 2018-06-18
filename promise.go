package torrent

import (
	"errors"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type PromiseDelivery chan interface{}

type Promise struct {
	sync.RWMutex

	set     bool
	value   interface{}
	waiters []PromiseDelivery
}

func (p *Promise) Set(value interface{}) {
	p.Lock()
	defer p.Unlock()
	if p.set {
		return
	}

	p.value = value
	p.set = true
	for _, w := range p.waiters {
		locW := w
		go func() {
			locW <- value
		}()
	}
	// clean up the waiters
	p.waiters = []PromiseDelivery{}
}

func (p *Promise) Get() interface{} {
	p.Lock()
	// defer p.Unlock()     // can not do this, because return <-delivery will block the call
	if p.set {
		return p.value
	}

	delivery := make(PromiseDelivery)
	defer func() {
		close(delivery)
	}()
	p.waiters = append(p.waiters, delivery)
	p.Unlock()
	return <-delivery
}

func (p *Promise) GetTimeout(timeout time.Duration) (interface{}, error) {
	// Because the Get method return a value, not a channel. So do this
	valueTmp := make(chan interface{})
	defer close(valueTmp)
	// Do not use defer to close the valueTmp channel, with timeout
	// a message will be sent to closed channel.
	// NOTICE: use receover to catch the close channel error
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Debugln("promise get error:", err)
				return
			}
		}()
		valueTmp <- p.Get()
	}()
	select {
	case <-time.After(timeout):
		//		log.Debugln("Promise get timeout")
		return nil, errors.New("timeout")
	case value := <-valueTmp:
		return value, nil
	}
}

func (p *Promise) IsSet() bool {
	return p.set
}

func (p *Promise) Size() int {
	return len(p.waiters)
}

func NewPromise() *Promise {
	return &Promise{
		value:   nil,
		waiters: []PromiseDelivery{},
		set:     false,
	}
}
