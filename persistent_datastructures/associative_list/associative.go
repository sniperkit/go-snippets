package associative_list

type Interface interface {
	Value(k interface{}) interface{}
	Set(k, v interface{}) Interface

	// Iterate calls f with all key-value pairs in the map
	Iterate(f func(k, v interface{}))
}

type empty struct{}

func (empty) Value(_ interface{}) interface{} {
	return nil
}

func (empty) Set(k interface{}, v interface{}) Interface {
	return pair{k, v, empty{}}
}

func (e empty) Iterate(f func(k, v interface{})) {
}

type pair struct {
	k, v   interface{}
	parent Interface
}

func (p pair) Value(k interface{}) interface{} {
	if k == p.k {
		return p.v
	}
	return p.parent.Value(k)
}

func (p pair) Set(k, v interface{}) Interface {
	// Same key will be hidden by the last one
	return pair{k, v, p}
}

func (p pair) Iterate(f func(k, v interface{})) {
	f(p.k, p.v)
	p.parent.Iterate(f)
}

type Map struct {
	Interface
}

func (m Map) Value(k interface{}) interface{} {
	return m.Interface.Value(k)
}

func (m Map) Set(k, v interface{}) Map {
	return Map{pair{k, v, m.Interface}}
}

func (m Map) Iterate(f func(k, v interface{})) {
	seen := make(map[interface{}]bool)
	m.Interface.Iterate(func(k, v interface{}) {
		if !seen[k] {
			f(k, v)
			seen[k] = true
		}
	})
}
