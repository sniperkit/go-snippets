package rest

import "net/http"

var (
	_ http.RoundTripper = (*RoundTripper)(nil)
)

// RoundTripper provides a easy way to implement the http.RoundTripper interface
// for simulating HTTP responses without having to make an HTTP request
type RoundTripper struct {
	Error    error
	Response *http.Response
}

// RoundTrip implements the http.RoundTripper interface
func (r RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return r.Response, r.Error
}
