package engine

import (
	"fmt"
)

// Command object for communicating through the engine.
type Command struct {
	Name        string
	Source      *NameSpace
	Destination *NameSpace
	Payload     interface{}
}

// NameSpace object holds the namespace information for different modules.
type NameSpace struct {
	Type   string
	Module string
	ID     string
}

// String returns a string representation of a NameSpace.
func (ns *NameSpace) String() string {
	return fmt.Sprintf("%s:%s:%s", ns.Type, ns.Module, ns.ID)
}
