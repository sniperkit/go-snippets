package main

import (
	"fmt"
	"sync"
)

type SyncMap struct {
	sync.RWMutex // embedded.  see http://golang.org/ref/spec#Struct_types
	hm           map[string]string
}

// Mue's function
func NewSyncMap() *SyncMap {
	return &SyncMap{hm: make(map[string]string)}
}

func (m *SyncMap) Put(k, v string) {
	m.Lock()    // note lock.Lock stutter is gone
	m.hm[k] = v // , true is from a very old version of Go, and even then it is not needed in this case.
	m.Unlock()  // I wouldn't bother with defer in such a simple function.
}

func (m *SyncMap) Get(k string) string {
	m.RLock()
	v := m.hm[k]
	m.RUnlock()
	return v
}

func main() {
	sm := NewSyncMap()
	sm.Put("kTest", "vTest")
	fmt.Println(sm.Get("kTest"))
}
