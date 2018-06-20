package file

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/go-fsnotify/fsnotify"
	yaml "gopkg.in/yaml.v2"

	"github.com/go-mixins/loader"
)

// DebounceTimeout defines change event settle time in milliseconds
var DebounceTimeout = 500

// UnmarshalFunc parses provided data into object
type UnmarshalFunc func(data []byte, dest interface{}) error

// Loader implements loader.Loader for a generic file on disk
type Loader struct {
	name          string
	stop, changes chan struct{}
	result        chan error
	watcher       *fsnotify.Watcher
	err           error
	f             UnmarshalFunc
}

var _ loader.Loader = (*Loader)(nil)

// New creates Loader initialized with a file name
func New(name string, f UnmarshalFunc) (res *Loader) {
	res = &Loader{
		name:    name,
		f:       f,
		stop:    make(chan struct{}),
		changes: make(chan struct{}, 1),
		result:  make(chan error, 1),
	}
	if res.watcher, res.err = fsnotify.NewWatcher(); res.err != nil {
		res.err = loader.Errors.Wrap(res.err, "creating fsnotify watcher")
		return
	}
	go func() {
		defer close(res.changes)
		defer func() {
			res.result <- res.watcher.Close()
			close(res.result)
		}()
		for {
			select {
			case <-res.stop:
				return
			case <-res.watcher.Events:
			loop:
				for {
					// This loop will consume consecutive change events that
					// will come during DebounceTimeout. Only the last
					// one will be reported, after the timeout expires.
					select {
					case <-res.watcher.Events:
						break
					case <-time.After(time.Duration(DebounceTimeout) * time.Millisecond):
						break loop
					}
				}
				res.changes <- struct{}{}
			}
		}
	}()
	if res.err = res.watcher.Add(res.name); res.err != nil {
		res.err = loader.Errors.Wrapf(res.err, "adding %q to watcher", res.name)
		return
	}
	return
}

// JSON returns loader of JSON format
func JSON(name string) *Loader {
	return New(name, json.Unmarshal)
}

// YAML returns loader of YAML format
func YAML(name string) *Loader {
	return New(name, yaml.Unmarshal)
}

// Close stops background process and releases FS watcher
func (l *Loader) Close() error {
	close(l.stop)
	return loader.Errors.Wrap(<-l.result, "closing watcher")
}

// Changes provides source of config change events. For environment variables
// there will never be any.
func (l *Loader) Changes() <-chan struct{} {
	return l.changes
}

// Load target object from a file
func (l *Loader) Load(dest interface{}) error {
	if l.err != nil {
		return l.err
	}
	data, err := ioutil.ReadFile(l.name)
	if err != nil {
		return loader.Errors.Wrap(err, "read file")
	}
	return loader.Errors.Wrap(l.f(data, dest), "unmarshal data")
}
