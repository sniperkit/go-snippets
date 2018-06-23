// +build !appengine

package rest

import (
	"net/http"

	"github.com/bradberger/context"
)

// TestClient creates a http.Client which will return the given response
func TestClient(ctx context.Context, resp *http.Response, err error) *http.Client {
	return &http.Client{Transport: RoundTripper{Response: resp, Error: err}}
}
