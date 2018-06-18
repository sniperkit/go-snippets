package locking

import (
	"errors"
)

var (
	// Lock is not currently held.
	ErrNotLocked = errors.New("lock is not currently held")

	// Lock is already held.
	ErrLocked = errors.New("lock is already held")

	// Lock acquisition timed out.
	ErrTimeout = errors.New("lock acquisition timed out")

	// Lock cancelled.
	ErrCancelled = errors.New("lock acquisition cancelled")
)
