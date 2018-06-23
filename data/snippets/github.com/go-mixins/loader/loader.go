package loader

import "github.com/go-mixins/errors"

// Errors defines error class that all returned errors belong to
var Errors = errors.NewClass("config.loader")

// Loader updates fields in the target object
type Loader interface {
	// Load the specified object's fields
	Load(dest interface{}) error
	// The optional source of change events
	Changes() <-chan struct{}
}

//go:generate moq -out mock/loader.go -pkg mock . Loader
