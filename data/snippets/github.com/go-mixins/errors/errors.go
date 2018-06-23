// Package errors is a simple wrapper around excellent "github.com/pkg/errors"
// package by Dave Cheney.
//
// It add a notion of "error class" that is simply chain of strings. Functions
// and methods are provided to mark the error with such a chain.
//
// The error is said to belong to a certain class if the class chain
// prefix-matches error's chain. The typical workflow is something like:
//
//
//     import "github.com/go-mixins/errors"
//
//     HTTPErrors := errors.NewClass("http") // A general "HTTP" error class
//
//     FatalErrors := HTTPErrors.Sub("fatal") // Specific "bad" subclass
//
//
//     func makeRequest() error {
//         result, err := http.Do(request)
//         err = HTTPErrors.Wrap(err, "requesting service") // return error of class HTTPError if not nil
//         return result, err
//     }
//
//     func main() {
//         ...
//         err := makeRequest()
//         switch {
//                case FatalErrors.Contains(err):
//                        panic(err) // error belongs to the Fatal class
//                case HTTPErrors.Contains(err):
//                    log.Printf("%+v", err) // just some HTTP error
//         }
//         ...
//     }
//
package errors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// root is the class that contains all other classes
var root = Class(nil)

// Error has a place in error class hierarchy
type Error struct {
	class Class
	error
}

func (ec Error) Error() string {
	return ec.class.String() + ": " + ec.error.Error()
}

// Cause returns the underlying error
func (ec Error) Cause() error {
	return ec.error
}

// Class returns error's class
func (ec Error) Class() Class {
	return ec.class
}

// Class represents error hierarchy path
type Class []string

// NewClass creates independent new path class hierarchy
func NewClass(path ...string) Class {
	return root.Sub(path...)
}

// Cause is provided for compatibility with "github.com/pkg/errors" package.
// It simply calls the original errors.Cause function
func Cause(err error) error {
	return errors.Cause(err)
}

// New is provided for compatibility with the standard Go "errors" package
func New(message string) error {
	return root.New(message)
}

// Wrap is provided for compatibility with "github.com/pkg/errors" package.
// It wraps error into empty root class
func Wrap(err error, message string) error {
	return root.Wrap(err, message)
}

// Wrapf is provided for compatibility with "github.com/pkg/errors" package.
// It wraps error into empty root class
func Wrapf(err error, format string, args ...interface{}) error {
	return root.Wrapf(err, format, args...)
}

// Errorf is provided for compatibility with "github.com/pkg/errors" package.
// It returns an error formatted against supplied format in empty root class.
func Errorf(format string, args ...interface{}) error {
	return root.Errorf(format, args...)
}

// Sub creates subpath in the class hierarchy
func (c Class) Sub(path ...string) (res Class) {
	res = make([]string, len(c)+len(path))
	copy(res, c)
	copy(res[len(c):], path)
	return
}

func (c Class) String() string {
	return strings.Join(c, ".")
}

// Is returns true if the class belongs to specific parent class
func (c Class) Is(parent Class) bool {
	if len(parent) > len(c) {
		return false
	}
	for i := range parent {
		if parent[i] != c[i] {
			return false
		}
	}
	return true
}

type classer interface {
	Class() Class
}

type causer interface {
	Cause() error
}

// Contains is true if the error belongs to certain class
func (c Class) Contains(err error) bool {
	if e, ok := err.(classer); ok {
		return e.Class().Is(c)
	}
	if e, ok := err.(causer); ok {
		return c.Contains(e.Cause())
	}
	return false
}

type simpleError string

func (se simpleError) Error() string {
	return string(se)
}

// Wrap marks the error with certain class and wraps it using errors package
func (c Class) Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(Error{c, err}, message)
}

// Wrapf marks the error with certain class and wraps it using errors package
func (c Class) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(Error{c, err}, format, args...)
}

// New returns an error with the supplied message and class
func (c Class) New(message string) error {
	return errors.WithStack(Error{c, simpleError(message)})
}

// Errorf returns an error formatted against supplied format
func (c Class) Errorf(format string, args ...interface{}) error {
	return errors.WithStack(Error{c, fmt.Errorf(format, args...)})
}
