// +build appengine

package rest

import (
	"net/http"

	"google.golang.org/appengine/urlfetch"

	"github.com/bradberger/context"
)

// TestClient creates a http.Client which will return the given response
func TestClient(ctx context.Context, resp *http.Response, err error) *http.Client {
	client := urlfetch.Client(ctx)
	client.Transport = RoundTripper{Response: resp, Error: err}
	return client
}
