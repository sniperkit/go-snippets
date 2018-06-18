package leadership

import (
	"errors"
)

var (
	// No leader found.
	ErrNoLeader = errors.New("no leader found")
)
