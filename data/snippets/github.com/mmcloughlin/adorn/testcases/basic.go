package testcases

// Basic is pretty simple.
type Basic interface {
	Method(int) string
}

// BasicFunc is an adapter to allow ordinary functions to be used as Basic implementations.
type BasicFunc func(int) string

// Method calls f.
func (f BasicFunc) Method(a0 int) string {
	return f(a0)
}

// Compile time assertion that BasicFunc implements Basic.
var _ Basic = new(BasicFunc)
