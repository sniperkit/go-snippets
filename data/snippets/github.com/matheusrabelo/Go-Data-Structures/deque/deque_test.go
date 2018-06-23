package deque

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeque(t *testing.T) {
	assert := assert.New(t)
	d := New()

	d.PushFirst("a")
	d.PushFirst("b")
	d.PushLast("c")
	d.PushLast("d")

	assert.Equal("b", d.PopFirst())
	assert.Equal("a", d.First())
	assert.Equal("d", d.PopLast())
	assert.Equal("c", d.Last())
	assert.Equal(2, d.Size())
	assert.Equal(false, d.Empty())

}
