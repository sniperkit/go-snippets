package insight

import (
	insightV1 "github.com/homebot/protobuf/pkg/api/insight/v1"
)

type MessageFilter interface {
	// Match checks if the current filter matches on the message
	Match(*insightV1.LogMessage) bool

	String() string
}

type Filter interface {
	// Match checks if a given message matches all filters in the chain
	Match(*insightV1.LogMessage) bool

	// First returns the first message filter in the filter chain
	First() MessageFilter

	// Len returns the length of the filter chain
	Len() int
}

type Storage interface {
	// Save saves a new log message
	Save(*insightV1.LogMessage) error

	// List lists all log messages
	List() ([]*insightV1.LogMessage, error)

	// Delete deletes a log message
	Delete(id string) error

	// Get returns the message with the given ID
	Get(id string) (*insightV1.LogMessage, error)

	// Search searches for messages
	Search(iterator Filter) ([]*insightV1.LogMessage, error)
}
