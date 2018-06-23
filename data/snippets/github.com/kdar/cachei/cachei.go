package cachei

import (
	"fmt"
	"github.com/kdar/cachei/msgpack"
	//"github.com/kdar/cache/json"
)

type DataSource struct {
	Host  string
	Port  int
	Coder Coder
}

type Coder interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type CacheFunc func() (interface{}, error)

type Cacher interface {
	// Gets a cache value and instead of returning it, puts it into "out".
	// If no cache is found, the CacheFunc is called and put into cache.
	// It returns the CacheFunc error, and caching error
	OutSetFn(key string, expires int, out interface{}, cFunc CacheFunc) (funcErr error, cacheErr error)

	// Gets a cache value.
	// If no cache is found, the CacheFunc is called and put into cache.
	// It returns the CacheFunc error, and caching error
	GetSetFn(key string, expires int, cFunc CacheFunc) (ret interface{}, funcErr error, cacheErr error)

	Out(key string, out interface{}) (cacheErr error)
	Get(key string) (ret interface{}, cacheErr error)
	Set(key string, value interface{}, expires int) (cacheErr error)

	// Delete a key from the cache.
	// Returns a caching error.
	Delete(key string) (cacheErr error)

	Open() error
	Close() error
	Setup(DataSource) error
}

var wrappers = make(map[string]Cacher)

func Register(name string, driver Cacher) {
	if name == "" {
		panic("cache: Wrapper name cannot be nil.")
	}
	if _, ok := wrappers[name]; ok != false {
		panic("cache: Wrapper was already registered.")
	}
	wrappers[name] = driver
}

func Open(name string, settings DataSource) (Cacher, error) {
	if _, ok := wrappers[name]; ok == false {
		panic(fmt.Sprintf("cache: Unknown wrapper: %s.", name))
	}

	if settings.Coder == nil {
		settings.Coder = &msgpack.Coder{}
	}

	err := wrappers[name].Setup(settings)
	if err != nil {
		return nil, err
	}
	return wrappers[name], nil
}

func Wrapper(name string) Cacher {
	return wrappers[name]
}
