package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	assert := assert.New(t)
	q := New()

	q.Push(44)
	q.Push("testing")
	q.Push(128)
	q.Push(9086)

	assert.Equal(44, q.Pop())
	assert.Equal("testing", q.Pop())
	assert.Equal(2, q.Size())
	assert.Equal(false, q.Empty())

	q.Push(50)

	assert.Equal(128, q.First())
}
