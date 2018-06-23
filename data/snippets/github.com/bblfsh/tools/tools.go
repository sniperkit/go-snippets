package tools

import "gopkg.in/bblfsh/sdk.v1/uast"

// Tooler is an interface which can be implemented by any supported tool.
// When implemented, the Exec method will be called with a UAST root node.
type Tooler interface {
	// Exec will be called with a UAST root node. The error will be passed
	// to the command handler
	Exec(*uast.Node) error
}
