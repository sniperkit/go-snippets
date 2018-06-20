package concurrent

import "sync"

// Map that is safe for concurrent access.
type Map struct {
	data  map[string]interface{}
	mutex sync.Mutex
}

// NewMap creates a new concurrent map.
func NewMap() *Map {
	return &Map{
		data: make(map[string]interface{}),
	}
}

// PutIfAbsent if the key k is absent from the map, (k,v) is put into the map.
// The inserted value and true is returned. Otherwise, the map is not mutated
// and the existing value at k and false is returned.
func (m *Map) PutIfAbsent(k string, v interface{}) (interface{}, bool) {
	m.mutex.Lock()
	if v, ok := m.data[k]; ok {
		m.mutex.Unlock()
		return v, false
	}
	m.data[k] = v
	m.mutex.Unlock()
	return v, true
}
