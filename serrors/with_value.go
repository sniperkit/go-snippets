package serrors

import (
	"fmt"
	"io"
)

//This package is for Structured Errors (hence `serrors`) which adds key-value pairs to an error.
//see https://github.com/pkg/errors/issues/34 and the closed/rejected PR here https://github.com/pkg/errors/pull/127

// WithValue returns a copy of parent in which the value associated with key is
// val.
func WithValue(err error, key string, val interface{}) error {
	if err == nil {
		return nil
	}
	return &withValue{err, key, val}
}

type withValue struct {
	cause error
	key   string
	val   interface{}
}

func (w *withValue) Error() string { return fmt.Sprintf("[%s=%+v] ", w.key, w.val) + w.cause.Error() }
func (w *withValue) Cause() error  { return w.cause }

func (w *withValue) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, fmt.Sprintf("%s=%+v", w.key, w.val))
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

//UniqueValues returns the unique+ordered keys and values map
func UniqueValues(err error) (orderedKeys []string, values map[string]interface{}) {
	type causer interface {
		Cause() error
	}

	type keyVal struct {
		k string
		v interface{}
	}

	pairs := []*keyVal{}
	for err != nil {
		if w, ok := err.(*withValue); ok {
			pairs = append(pairs, &keyVal{w.key, w.val})
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	orderedKeys = []string{}
	m := make(map[string]interface{})

	// Looping backward since assuming we want the first added value.
	// This is because each time we add the same key (with `WithValue(`) it will wrap the original.
	for i := len(pairs) - 1; i >= 0; i-- {
		pair := pairs[i]
		if _, alreadyInMap := m[pair.k]; !alreadyInMap {
			orderedKeys = append(orderedKeys, pair.k)
			m[pair.k] = pair.v
		}
	}

	return orderedKeys, m
}
