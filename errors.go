package skiplist

import "errors"

var (
	ErrNilKey         = errors.New("nil key")
	ErrNilValue       = errors.New("nil value")
	ErrKeyExists      = errors.New("key already exists")
	ErrNilAction      = errors.New("nil action")
	ErrUnknownFromKey = errors.New("unknown from key")
	ErrUnknownToKey   = errors.New("unknown to key")
	ErrRemoveNilNode  = errors.New("can't remove nil node")
)
