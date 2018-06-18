// Package pool implements a pool of Object interfaces to manage and reuse them.
package pool

import (
	"errors"
	"reflect"
)

var (
	ErrClosed = errors.New("pool is closed")
)

// Pool interface describes a pool implementation. A pool should have maximum
// capacity. An ideal pool is threadsafe and easy to use.
type Pool interface {
	Transaction(fn interface{}) ([]reflect.Value, error)

	Borrow() (PoolObjecter, error)

	Close()

	// Len returns the current number of connections of the pool.
	Len() int
}
