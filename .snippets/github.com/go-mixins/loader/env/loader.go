package env

import (
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/go-mixins/loader"
)

// Loader implements loader.Loader
type Loader struct {
	prefix  string
	changes chan struct{}
}

var _ loader.Loader = (*Loader)(nil)

// Load loads the target from environment
func (l *Loader) Load(dest interface{}) error {
	return loader.Errors.Wrap(envconfig.Process(l.prefix, dest), "load from environment")
}

// Close closes underlying changes channel
func (l *Loader) Close() error {
	close(l.changes)
	return nil
}

// Changes provides source of config change events. For environment variables
// there will never be any.
func (l *Loader) Changes() <-chan struct{} {
	return l.changes
}

// New creates loader initialized with environment variable prefix
func New(prefix string) *Loader {
	return &Loader{
		prefix:  strings.ToUpper(prefix),
		changes: make(chan struct{}),
	}
}
