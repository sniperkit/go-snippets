package engine

import (
	"errors"
)

// Latest version of our hash
// TODO: should auto calculate...
const LatestHashVersion = 1

// Setup interface
type Setup interface {
	Setup(eng *Engine, cfg Config) error
}

// StartStop interface
type StartStop interface {
	Start() error
	Stop() error
}

// Commander interface
type Commander interface {
	Command(cmd *Command) error
}

// ErrNoContainerRuntime ...
var ErrNoContainerRuntime = errors.New("no container runtime started")

// ErrNotAValidConfig ...
var ErrNotAValidConfig = errors.New("not a valid config")

// ErrUnknownEvent ...
var ErrUnknownEvent = errors.New("unknown event")

// ErrLoopDetected
var ErrLoopDetected = errors.New("loop detected")

// ErrNotAllowed ...
var ErrNotAllowed = errors.New("action not allowed")

// ErrInvalidPayload ...
var ErrInvalidPayload = errors.New("not a valid payload")

// ErrCircularDependency ...
var ErrCircularDependency = errors.New("circular dependecy detected")
