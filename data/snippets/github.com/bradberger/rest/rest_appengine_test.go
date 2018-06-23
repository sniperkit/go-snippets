// +build appengine

package rest

import (
	"github.com/bradberger/context"
	"google.golang.org/appengine/aetest"
)

// NewTestContext returns a new context suitable for testing. In the appengine
// environment, this is an aetest context
func NewTestContext() (context.Context, func(), error) {
	return aetest.NewContext()
}
