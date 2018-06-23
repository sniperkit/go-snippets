package binaryTree

import (
	"testing"

	"strconv"

	"github.com/stretchr/testify/assert"
)

var s string

func TestBinaryTreeTraversal(t *testing.T) {
	assert := assert.New(t)

	btree := New(15)
	btree.PutLeft(66)
	btree.PutRight(234)
	btree.Left().PutLeft(128)
	btree.Left().PutRight(45)
	btree.Right().PutLeft(80)

	s = ""
	btree.PreOrder(addValueToString)
	assert.Equal("15|66|128|45|234|80|", s)

	s = ""
	btree.InOrder(addValueToString)
	assert.Equal("128|66|45|15|80|234|", s)

	s = ""
	btree.PostOrder(addValueToString)
	assert.Equal("128|45|66|80|234|15|", s)

}

func addValueToString(value interface{}) interface{} {
	s += strconv.Itoa(value.(int)) + "|"
	return nil
}
