package discovery

// SortedMap is a unconcurrent-safe sorted map
type SortedMap struct {
	keys     []string
	length   int
	capacity int
	m        map[string]interface{}
}

// NewSortedMap create and init a new sorted map
func NewSortedMap() *SortedMap {
	return &SortedMap{
		keys:     make([]string, 4),
		length:   0,
		capacity: 4,
		m:        make(map[string]interface{}),
	}
}

func (sm *SortedMap) findPlace(key string) (index int) {
	low, high := 0, sm.length-1
	mid := 0

	for low <= high {
		mid = (low + high) / 2
		if sm.keys[mid] == key {
			return mid
		}

		if sm.keys[mid] < key {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low
}

func (sm *SortedMap) adjustKeys(newlen int) {
	if newlen > sm.capacity {
		sm.capacity *= 2
		if newlen > sm.capacity {
			sm.capacity = newlen * 2
		}

		newkeys := make([]string, sm.capacity)
		copy(newkeys, sm.keys)
		sm.keys = newkeys

	} else if newlen < sm.capacity/4 {
		sm.capacity /= 4
		newkeys := make([]string, sm.capacity)
		copy(newkeys, sm.keys)
		sm.keys = newkeys
	}
}

// AddKV ...
func (sm *SortedMap) AddKV(key string, value interface{}) (index int) {
	p := sm.findPlace(key)
	if sm.keys[p] == key {
		sm.m[key] = value
		return p
	}

	sm.adjustKeys(sm.length + 1)
	for i := sm.length; i > p; i-- {
		sm.keys[i] = sm.keys[i-1]
	}

	sm.keys[p] = key
	sm.length += 1
	sm.m[key] = value
	return p
}

// DeleteKV ...
func (sm *SortedMap) DeleteKV(key string) (value interface{}) {
	oldvalue, ok := sm.m[key]
	if !ok {
		return nil
	}

	delete(sm.m, key)

	p := sm.findPlace(key)
	if sm.keys[p] != key {
		return oldvalue
	}

	sm.adjustKeys(sm.length - 1)
	for i := p + 1; i < sm.length; i++ {
		sm.keys[i-1] = sm.keys[i]
	}

	sm.length -= 1
	return oldvalue
}

// SetKV ...
func (sm *SortedMap) SetKV(key string, value interface{}) (oldvalue interface{}) {
	oldvalue, ok := sm.m[key]
	if !ok {
		sm.AddKV(key, value)
		return nil
	}
	sm.m[key] = value
	return oldvalue
}

func (sm *SortedMap) GetKV(key string) (value interface{}) {
	return sm.m[key]
}

func (sm *SortedMap) GetIV(index int) (value interface{}) {
	return sm.m[sm.keys[index]]
}

func (sm *SortedMap) Len() int {
	return sm.length
}

func (sm *SortedMap) Cap() int {
	return sm.capacity
}

func (sm *SortedMap) Keys() []string {
	snapshot := make([]string, sm.length)
	copy(snapshot, sm.keys)
	return snapshot
}
