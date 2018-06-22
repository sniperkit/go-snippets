package container

// Action a container can be in
type Action int

// List of actions
const (
	Done Action = iota
	Error
	Creating
	Starting
	Stopping
	Removing
	Needing
)
