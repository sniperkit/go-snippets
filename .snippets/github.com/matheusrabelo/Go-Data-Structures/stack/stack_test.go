package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	assert := assert.New(t)
	s := New()

	s.Push(23)
	s.Push(3123)
	s.Push(100)
	s.Push(33)

	assert.Equal(33, s.Pop())
	assert.Equal(100, s.Pop())
	assert.Equal(2, s.Size())
	assert.Equal(false, s.Empty())

	s.Push("some string")

	assert.Equal("some string", s.Top())
}
