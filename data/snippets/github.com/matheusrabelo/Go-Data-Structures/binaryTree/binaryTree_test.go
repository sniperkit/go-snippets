package binaryTree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryTree(t *testing.T) {
	assert := assert.New(t)

	btree := New(20)
	left := btree.PutLeft(10)
	right := btree.PutRight(44)
	leftleft := left.PutLeft(5)
	leftright := left.PutRight(30)

	assert.Equal(10, left.Value())
	assert.Equal(44, right.Value())

	assert.Equal(5, leftleft.Value())
	assert.Equal(30, leftright.Value())

	assert.Equal(5, btree.Left().Left().Value())
	assert.Equal(10, btree.Left().Value())
	assert.Equal(30, btree.Left().Right().Value())
	assert.Equal(20, btree.Value())
	assert.Equal(44, btree.Right().Value())
}
