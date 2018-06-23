package io

import "errors"

// some codes inspired by go source code package io

var EOF error = errors.New("EOF")

type Writer interface {
	Write(b []byte) (n int, err error)
}

type Reader interface {
	Read(b []byte) (n int, err error)
}
