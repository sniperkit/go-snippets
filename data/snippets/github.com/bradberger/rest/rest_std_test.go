// +build !appengine,go1.7

package rest

import "github.com/bradberger/context"

// NewTestContext returns a new context suitable for testing. Outside of appengine,
// this is just a fresh background context
func NewTestContext() (context.Context, func(), error) {
	return context.Background(), func() {}, nil
}
