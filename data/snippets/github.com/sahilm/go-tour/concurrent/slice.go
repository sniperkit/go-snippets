package concurrent

import (
	"sync"
)

// Slice that is safe for concurrent access.
type Slice struct {
	data  []interface{}
	mutex sync.Mutex
}

// NewSlice creates a new concurrent slice.
func NewSlice() *Slice {
	return &Slice{
		data: make([]interface{}, 0),
	}
}

// Append appends all of v to slice
func (s *Slice) Append(v ...interface{}) {
	s.mutex.Lock()
	s.data = append(s.data, v...)
	s.mutex.Unlock()
}

// View returns a copy of all the elements
// in the slice
func (s *Slice) View() []interface{} {
	var v []interface{}
	s.mutex.Lock()
	v = append(v, s.data...)
	s.mutex.Unlock()
	return v
}
