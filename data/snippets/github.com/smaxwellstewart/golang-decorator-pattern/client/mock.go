package client

import (
	"fmt"
	"net/http"
)

// MockClient is used to do pretend http requests.
// Notice MockClient has no knowledge of the other decorators that have orthogonal concerns.
type MockClient struct{}

// Do will do a pretend http request. Do be doo.
func (mc MockClient) Do(r *http.Request) (*http.Response, error) {
	fmt.Println("Mock request to: " + r.Host)
	return &http.Response{}, nil
}
