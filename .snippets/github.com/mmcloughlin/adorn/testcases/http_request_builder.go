package testcases

import "net/http"

// Builder is anything that can build an HTTP request.
type Builder interface {
	Build() (*http.Request, error)
}

// BuilderFunc is an adapter to allow ordinary functions to be used as Builder implementations.
type BuilderFunc func() (*http.Request, error)

// Build calls f.
func (f BuilderFunc) Build() (*http.Request, error) {
	return f()
}

// Compile time assertion that BuilderFunc implements Builder.
var _ Builder = new(BuilderFunc)
