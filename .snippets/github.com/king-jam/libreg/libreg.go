package libreg

import (
	"fmt"
	"sort"
	"strings"

	"github.com/king-jam/libreg/registry"
)

// Initialize creates a new Store object, initializing the client
type Initialize func(addrs []string, options *registry.Config) (registry.Registry, error)

var (
	// Backend initializers
	initializers = make(map[registry.Backend]Initialize)

	supportedBackend = func() string {
		keys := make([]string, 0, len(initializers))
		for k := range initializers {
			keys = append(keys, string(k))
		}
		sort.Strings(keys)
		return strings.Join(keys, ", ")
	}()
)

// NewRegistry creates an instance of registry
func NewRegistry(backend registry.Backend, addrs []string, options *registry.Config) (registry.Registry, error) {
	if init, exists := initializers[backend]; exists {
		return init(addrs, options)
	}

	return nil, fmt.Errorf("%s %s", registry.ErrBackendNotSupported.Error(), supportedBackend)
}

// AddRegistry adds a new registry backend to libreg
func AddRegistry(backend registry.Backend, init Initialize) {
	initializers[backend] = init
}
