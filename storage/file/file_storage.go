package file

import (
	"errors"
	"sync"

	"github.com/golang/protobuf/ptypes"

	"github.com/homebot/insight"
	insightV1 "github.com/homebot/protobuf/pkg/api/insight/v1"
)

// TODO(ppacher): this is more of a PoC. Make this more efficient
// TODO(ppacher): actually save messages to file

// Storage is a insight.Storage implementation that saves
// log messages to files
type Storage struct {
	rw       sync.RWMutex
	messages []*insightV1.LogMessage
}

// Save saves a log message
func (store *Storage) Save(msg *insightV1.LogMessage) error {
	store.rw.Lock()
	defer store.rw.Unlock()

	store.messages = append(store.messages, msg)

	return nil
}

// Get returns the logging message with the given ID
func (store *Storage) Get(id string) (*insightV1.LogMessage, error) {
	store.rw.RLock()
	defer store.rw.RUnlock()

	for _, m := range store.messages {
		if m.Id == id {
			return m, nil
		}
	}

	return nil, errors.New("id not found")
}

// List lists all log messages sorted by time
func (store *Storage) List() ([]*insightV1.LogMessage, error) {
	store.rw.RLock()
	defer store.rw.RUnlock()

	return sortMessages(store.messages, func(a, b *insightV1.LogMessage) bool {
		t1, _ := ptypes.Timestamp(a.CreatedTime)
		t2, _ := ptypes.Timestamp(b.CreatedTime)

		return t1.After(t2)
	}), nil
}

// Delete deletes the log message with the given ID
func (store *Storage) Delete(id string) error {
	store.rw.Lock()
	defer store.rw.Unlock()

	var copy []*insightV1.LogMessage

	for _, msg := range store.messages {
		if msg.Id != id {
			copy = append(copy, msg)
		}
	}

	store.messages = copy
	return nil
}

// Search searches for log messages that match all criteria of filter
func (store *Storage) Search(filter insight.Filter) ([]*insightV1.LogMessage, error) {
	store.rw.RLock()
	defer store.rw.RUnlock()

	var res []*insightV1.LogMessage

	for _, msg := range store.messages {
		if filter.Match(msg) {
			res = append(res, msg)
		}
	}

	return sortMessages(res, func(a, b *insightV1.LogMessage) bool {
		t1, _ := ptypes.Timestamp(a.CreatedTime)
		t2, _ := ptypes.Timestamp(b.CreatedTime)

		return t1.After(t2)
	}), nil
}

func sortMessages(list []*insightV1.LogMessage, less func(a, b *insightV1.LogMessage) bool) []*insightV1.LogMessage {
	copy := make([]*insightV1.LogMessage, len(list))

	for idx, msg := range list {
		copy[idx] = msg

		if idx > 0 {
			if less(copy[idx], copy[idx-1]) {
				copy[idx], copy[idx-1] = copy[idx-1], copy[idx]
			}
		}
	}

	for {
		swapped := false

		for idx := range copy {
			if idx > 0 {
				if less(copy[idx], copy[idx-1]) {
					copy[idx], copy[idx-1] = copy[idx-1], copy[idx]
					swapped = true
				}
			}
		}

		if swapped == false {
			break
		}
	}

	return copy
}
