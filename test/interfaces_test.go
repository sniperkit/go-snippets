package test

import "testing"

func TestInterfaces(t *testing.T) {

	// 1. How to define?
	type Reader interface {
		Read(p []byte) (n int, err error)
	}

	type Closer interface {
		Close() error
	}

	// 2. How to compose interfaces
	type ReadCloser interface {
		Reader
		Closer
	}
}
