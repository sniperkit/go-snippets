package locking

// Lock provider.
type LockProvider interface {
	// Get a lock.
	GetLock(path string) Lock
}
