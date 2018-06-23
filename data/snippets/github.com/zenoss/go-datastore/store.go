package datastore

import (
	"errors"
)

type Storeable interface {
	Type() string
	Key() string
	Validate(context Context) error
}

type Context interface {
	Get(Storeable) error
	Put(Storeable) error
	Delete(Storeable) error
	Exists(Storeable) (bool, error)
}

type DataStore struct {
	schema string
}

var (
	ErrUnimplemented = errors.New("unimplemented")
	ErrNotFound      = errors.New("not found")
)

// Get()s a value from the datastore; s is a pointer to a storable object who's Type()
// and Key() methods return the appropiate type and primary key, respectively.
// The retrieved value is stored in s.
func Get(context Context, s Storeable) error {
	return context.Get(s)
}

// Put()s a value from the datastore; s is a pointer to a storable object who's Type()
// and Key() methods return the appropiate type and primary key, respectively
func Put(context Context, s Storeable) error {
	if err := s.Validate(context); err != nil {
		return err
	}
	return context.Put(s)
}

// Delete()s an object from the datastore; s is a pointer to a storable object who's Type()
// and Key() methods return the appropiate type and primary key, respectively
func Delete(context Context, s Storeable) error {
	return context.Delete(s)
}

// Exists() determines if a storable exists in the storagecontext
func Exists(context Context, s Storeable) (bool, error) {
	return context.Exists(s)
}

// TODO: implement search
