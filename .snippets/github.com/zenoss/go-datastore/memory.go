package datastore

import (
	"encoding/json"
)

// dataitem is the underlying data that the MemoryContext stores ([]byte)
type dataitem []byte

// datamap finds a unique dataitem by a key
type datamap map[string]dataitem

// typemap finds a datamap for a unique data type
type typemap map[string]datamap

// MemoryContext implements a memory based StorageContext
type MemoryContext struct {
	values typemap
}

// Assert that MemoryContext implements Context interface
var _ Context = MemoryContext{}

// NewMemoryContext() create a new MemoryContext object
func NewMemoryContext() (context *MemoryContext, err error) {
	context = &MemoryContext{
		values: make(typemap),
	}
	return context, nil
}

// Put stores val in the MemoryContext
func (context MemoryContext) Put(val Storeable) error {

	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	type_, ok := context.values[val.Type()]
	if !ok {
		type_ = make(datamap)
		context.values[val.Type()] = type_
	}
	type_[val.Key()] = data
	return nil
}

// Exists() checks if the given val exists in the MemoryContext
func (context MemoryContext) Exists(val Storeable) (bool, error) {
	type_, ok := context.values[val.Type()]
	if !ok {
		return ok, nil
	}
	_, ok = type_[val.Key()]
	return ok, nil
}

// Get() retrieves a storable from the MemoryContext
func (context MemoryContext) Get(val Storeable) error {
	type_, ok := context.values[val.Type()]
	if !ok {
		return ErrNotFound
	}
	bytes, ok := type_[val.Key()]
	if !ok {
		return ErrNotFound
	}
	return json.Unmarshal(bytes, &val)
}

// Delete() removes an item (val) from the MemoryContext
func (context MemoryContext) Delete(val Storeable) error {
	type_, ok := context.values[val.Type()]
	if !ok {
		return nil
	}
	if _, found := type_[val.Key()]; found {
		delete(type_, val.Key())
	}
	return nil
}
