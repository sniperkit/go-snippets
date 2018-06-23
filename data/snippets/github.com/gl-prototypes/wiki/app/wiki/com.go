package wiki

import (
	"github.com/gliderlabs/com"
	"github.com/gliderlabs/stdcom/web/console"
)

// since this package was not made for reuse,
// we can init to register directly in the package
func init() {
	com.Register(&Component{}, "")
}

// Component is the wiki system component
type Component struct {
	// Console gives us a reference to the
	// component wiki is parimarily built around
	Console *console.Component `com:"singleton"`

	// Store is our datastore for wiki pages
	Store Store `com:"singleton"`
}
