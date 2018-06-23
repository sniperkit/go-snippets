package multierr

import "bytes"

// Multierr struct allows to have multiple errors together
type Multierr struct {
	errors []error
}

// NewMultierr creates a new pointer to Multierr
func NewMultierr(errs ...error) *Multierr {
	m := Multierr{}

	m.errors = append(m.errors, errs...)

	return &m
}

// Add adds error to the collection of errors
func (m *Multierr) Add(err error) {
	m.errors = append(m.errors, err)
}

// Len return a number of errors in the collection
func (m *Multierr) Len() int {
	return len(m.errors)
}

// Error error interface implementation
func (m *Multierr) Error() string {
	var buffer bytes.Buffer

	for _, err := range m.errors {
		buffer.WriteString(err.Error())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
