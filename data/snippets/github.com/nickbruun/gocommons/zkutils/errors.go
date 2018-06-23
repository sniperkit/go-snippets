package zkutils

import (
	"errors"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	// Node is not a match.
	ErrNodeNotMatch = errors.New("node is not a match")

	// Timeout error.
	ErrTimeout = errors.New("timeout")
)

// Test if a ZooKeeper error is recoverable.
//
// Takes a conservative approach, and only considers authentication failures
// etc. as unrecoverable.
func IsErrorRecoverable(err error) bool {
	switch err {
	case zk.ErrNoAuth, zk.ErrNoChildrenForEphemerals, zk.ErrNotEmpty, zk.ErrInvalidACL, zk.ErrAuthFailed:
		return false

	default:
		return true
	}
}
