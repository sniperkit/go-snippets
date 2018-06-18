package lazy

import (
	"reflect"
	"sync"
)

// Lazy is a lazy value wrapper
type Lazy interface {
	Value() interface{}
}

type lazy struct {
	fn     reflect.Value
	params []reflect.Value
	once   sync.Once
	result interface{}
}

// New ...
func New(op interface{}, params ...interface{}) Lazy {
	t := reflect.TypeOf(op)
	if t.Kind() != reflect.Func {
		return nil
	}

	fn := reflect.ValueOf(op)
	in := make([]reflect.Value, len(params))

	for i := 0; i < len(params); i++ {
		if fn.Type().In(i) != reflect.TypeOf(params[i]) {
			return nil
		}
		in[i] = reflect.ValueOf(params[i])
	}

	return &lazy{
		fn:     fn,
		params: in,
	}
}

func (lazy *lazy) Value() interface{} {
	lazy.once.Do(func() {
		lazy.result = lazy.fn.Call(lazy.params)[0].Interface()
	})
	return lazy.result
}
