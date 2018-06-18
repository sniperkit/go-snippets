package pool

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Object interface {
	Close() error
}

var ObjecterType reflect.Type = reflect.TypeOf((*Object)(nil)).Elem()

// channelPool implements the Pool interface based on buffered channels.
type channelPool struct {
	mu      sync.Mutex
	objects chan Object

	maxCap  int
	currNum int

	// Object generator
	factory Factory
}

// Factory is a function to create new Object.
type Factory func() (Object, error)

// NewChannelPool returns a new pool based on buffered channels with an initial
// capacity and maximum capacity. Factory is used when initial capacity is
// greater than zero to fill the pool. A zero initialCap doesn't fill the Pool
// until a new Borrow() is called. During a Borrow(), If there is no new Object
// available in the pool, a new Object will be created via the Factory()
// method.
func NewChannelPool(initialCap, maxCap int, factory Factory) (Pool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	c := &channelPool{
		objects: make(chan Object, maxCap),
		factory: factory,
		maxCap:  maxCap,
	}

	// create initial Objects, if something goes wrong,
	// just close the pool error out.
	for i := 0; i < initialCap; i++ {
		obj, err := factory()
		if err != nil {
			c.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.objects <- obj
	}
	c.currNum = initialCap

	return c, nil
}

func (c *channelPool) getObjects() chan Object {
	c.mu.Lock()
	objects := c.objects
	c.mu.Unlock()
	return objects
}

func (c *channelPool) Transaction(fn interface{}) ([]reflect.Value, error) {
	t := reflect.TypeOf(fn)

	if t.NumIn() != 1 {
		return nil, errors.New("pool transaction function argument must be one")
	}

	firstArgType := t.In(0)
	if firstArgType != ObjecterType {
		return nil, errors.New("pool transaction function argument must be a Object")
	}

	obj, err := c.Borrow()
	if err != nil {
		return nil, err
	}
	defer obj.Return()

	poolObj := obj.(*PoolObject)
	return reflect.ValueOf(fn).Call([]reflect.Value{reflect.ValueOf(poolObj.Object)}), nil
}

// Borrow implements the Pool interfaces Borrow() method. If there is no new
// Object available in the pool, a new Object will be created via the
// Factory() method.
func (c *channelPool) Borrow() (PoolObjecter, error) {
	objects := c.getObjects()
	if objects == nil {
		return nil, ErrClosed
	}

	// we always try to create a new Object if there is no one exists.
	// If failed block self until there is one available

	select {
	case obj := <-objects:
		if obj == nil {
			return nil, ErrClosed
		}
		return c.wrapObj(obj), nil
	default:
		// try to create a new one
		var obj Object
		var err error
		c.mu.Lock()
		if c.currNum < c.maxCap {
			obj, err = c.factory()
			if err != nil {
				c.mu.Unlock()
				return nil, err
			}
			c.currNum++
			c.mu.Unlock()
			return c.wrapObj(obj), nil
		}
		c.mu.Unlock()

		// let this one block until someone put back a Object
		select {
		case obj := <-objects:
			if obj == nil {
				return nil, ErrClosed
			}
			return c.wrapObj(obj), nil
		case <-time.After(time.Second * 5):
			return nil, errors.New("borrow object wait timeout")
		}
	}
}

func (c *channelPool) decrease() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.currNum > 0 {
		c.currNum--
	}
	return c.currNum
}

// put puts the Object back to the pool. If the pool is full or closed,
// Object is simply closed. A nil Object will be rejected.
func (c *channelPool) put(obj Object) error {
	if obj == nil {
		//we try to decrease currNum, let someone to create a new one if necessary
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.currNum > 0 {
			c.currNum--
		}
		return errors.New("Object is nil. rejecting")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.objects == nil {
		// pool is closed, close passed obj
		return obj.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case c.objects <- obj:
		return nil
	default:
		// pool is full, close passed obj
		return obj.Close()
	}
}

func (c *channelPool) Close() {
	c.mu.Lock()
	objects := c.objects
	c.objects = nil
	c.factory = nil
	c.currNum = 0
	c.mu.Unlock()

	if objects == nil {
		return
	}

	close(objects)
	for obj := range objects {
		obj.Close()
	}
}

func (c *channelPool) Len() int { return len(c.getObjects()) }
