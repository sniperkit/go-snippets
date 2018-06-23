package optional

// Optional is a monad Maybe
type Optional interface {
	Apply(func(interface{}) Optional) Optional
	HasValue() bool
	Or(interface{}) interface{}
	Get() interface{}
}

// Wrap returns a value wrapped to monad
func Wrap(value interface{}) Optional {
	return &Just{value}
}

// None is a non-value way
type None struct{}

// Just is a value way
type Just struct {
	value interface{}
}

// Apply on None do nothing
func (n *None) Apply(f func(interface{}) Optional) Optional {
	return n
}

// HasValue for None is alwaya false
func (n *None) HasValue() bool {
	return false
}

// Or for None always return a param
func (n *None) Or(value interface{}) interface{} {
	return value
}

// Get for None always return a nil
func (n *None) Get() interface{} {
	return nil
}

// Apply executes f on monad value
func (j *Just) Apply(f func(interface{}) Optional) Optional {
	return f(j.value)
}

// HasValue is true for non-nil value
func (j *Just) HasValue() bool {
	return j.value != nil
}

// Or returns a param if monad value is nil
func (j *Just) Or(value interface{}) interface{} {
	if j.HasValue() {
		return j.value
	}
	return value
}

// Get for Just return a value
func (j *Just) Get() interface{} {
	return j.value
}
