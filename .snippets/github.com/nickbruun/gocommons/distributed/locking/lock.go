package locking

import (
	"time"
)

// Lock.
type Lock interface {
	// Lock.
	//
	// Only returns an error if an unrecoverable error occurs. Otherwise,
	// returns a channel which receives an empty struct if the lock is lost.
	Lock() (<-chan struct{}, error)

	// Lock with timeout.
	//
	// Only returns an error if an unrecoverable error occurs or the lock is
	// not acquired before a certain timeout is reached. Otherwise, returns
	// a channel which receives an empty struct if the lock is lost.
	LockTimeout(time.Duration) (<-chan struct{}, error)

	// Unlock.
	//
	// Only returns an error if an unrecoverable error occurs.
	Unlock() error
}
