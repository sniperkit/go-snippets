package consul

import (
	"time"

	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"

	"github.com/go-mixins/loader"
	"github.com/go-mixins/loader/libkv"
)

// Loader implements Consul Loader
type Loader struct {
	*libkv.Loader
	err error
}

// New creates Consul loader initialized with specific prefix and endpoints
func New(prefix string, endpoints ...string) (res *Loader) {
	res = new(Loader)
	kv, err := consul.New(
		endpoints,
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	if err != nil {
		res.err = loader.Errors.Wrap(err, "creating Consul source")
		return
	}
	res.Loader, res.err = libkv.New(prefix, kv)
	return
}

// Load loads the target from Consul source
func (l *Loader) Load(dest interface{}) error {
	if l.err != nil {
		return l.err
	}
	return l.Loader.Load(dest)
}
