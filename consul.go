package quickstart

import (
	"context"
	"time"

	"github.com/hashicorp/consul/api"
)

// WatchHandleFunc defines the callback function type used for watch callbacks
type WatchHandleFunc func(kvPair *api.KVPair) error

// WithLockFunc defines the callback function type used while holding a lock
type WithLockFunc func() error

// ConsulAPI represents a consul API client
type ConsulAPI struct {
	*api.Client
}

// Watch monitors a given prefix for changes and executes the provided fn
func (c *ConsulAPI) Watch(ctx context.Context, prefix string, fn WatchHandleFunc) (err error) {
	currentIndex := uint64(0)

	for {
		// Prepare API query options
		opts := &api.QueryOptions{
			WaitIndex: currentIndex,
			WaitTime:  time.Minute,
		}

		// Invoke the List API call, passing in our context which allows the caller
		// to cancel this function call.
		kvPairs, meta, listErr := c.KV().List(prefix, opts.WithContext(ctx))

		if listErr != nil {
			// In case of an error break out of our loop
			err = listErr
			break
		}

		if currentIndex == meta.LastIndex {
			// Nothing changed, move on
			continue
		}

		// Iterate over all changed KV pairs
		for _, kvPair := range kvPairs {
			if kvPair.CreateIndex <= currentIndex && kvPair.ModifyIndex <= currentIndex {
				// This kvPair was not changed, but we are still receiving this
				// because a child has changed.
				// However, we ignore this change as we are only interested
				// in leafs.
				continue
			}

			if err = fn(kvPair); err != nil {
				// If the callback function returns an error abort the loop
				return
			}
		}

		// Update currentIndex, so the API only reports changes after the change
		// we just processed
		currentIndex = meta.LastIndex
	}

	return
}

// WithLock executes the provided WithLockFunc while holding the lock for the given key
func (c *ConsulAPI) WithLock(key string, fn WithLockFunc) (err error) {
	var lk *api.Lock

	// Obtain/create the lock key
	if lk, err = c.LockKey("locks/" + key); err != nil {
		return
	}
	// Ensure the lock is cleaned up when we leave this scope
	defer lk.Destroy()

	// Acquire the lock
	if _, err = lk.Lock(nil); err != nil {
		return
	}
	// Ensure we release the lock when leaving this scope
	defer lk.Unlock()

	// Call the provided function while we hold the lock
	err = fn()
	return
}

// NewConsulAPI returns a new consul API object
func NewConsulAPI(address string) (c *ConsulAPI, err error) {
	var consulClient *api.Client

	clientConfig := &api.Config{
		Address: address,
		Scheme:  "http",
	}

	// Construct the client
	if consulClient, err = api.NewClient(clientConfig); err != nil {
		return
	}

	// Construct the ConsulAPI object
	c = &ConsulAPI{
		Client: consulClient,
	}
	return
}
