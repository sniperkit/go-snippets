package mapdecor

import "sync"

var registry struct {
	decorators []Decor
	mutex      sync.Mutex
}

// RegistryCount - number of registered functions
func RegistryCount() int {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	return len(registry.decorators)
}

// Register - register new decor functions
func Register(decors ...Decor) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	for _, decor := range decors {
		registry.decorators = append(registry.decorators, decor)
	}
}

// ResetRegistry - reset currently registered functions
func ResetRegistry() {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	registry.decorators = []Decor{}
}
