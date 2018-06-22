package cache

// Package cache provides interfaces for cache functions

// Cache is the common interface implemented by all cache functions
type Cache interface {
	// Add an element to the cache
	Add()
	// Delete an element from the cache
	Delete()
	// Lookup an element in the cache
	Lookup()
}
