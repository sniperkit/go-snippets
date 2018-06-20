package genericop

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type binOpCase struct {
	x interface{}
	y interface{}
	r interface{}
}

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	cases := []binOpCase{
		{1, 2, 1 + 2},
		{1 + 2i, 3 + 4i, 1 + 2i + 3 + 4i},
		{12.3, 4.5, 12.3 + 4.5},
		{"xx", "yy", "xx" + "yy"},
		{1, 2 + 0i, nil},
		{1.0, 2, nil},
		{"1", 2, nil},
		{int(1), int8(2), nil},
	}

	for _, c := range cases {
		r, err := Add(c.x, c.y)
		if c.r == nil {
			assert.Error(err)
		} else {
			assert.NoError(err)
			assert.Equal(c.r, r, "%v + %v = %v", c.x, c.y, c.r)
		}
	}
}

func TestLt(t *testing.T) {
	assert := assert.New(t)
	cases := []binOpCase{
		{1, 2, true},
		{1.1, 0.0, false},
		{1.1, float32(0.0), nil},
		{"a", "b", true},
		{"x", "x", false},
		{1 + 1i, 2 + 2i, nil},
		{"a", 1, nil},
	}

	for _, c := range cases {
		r, err := Lt(c.x, c.y)
		if c.r == nil {
			assert.Error(err)
		} else {
			assert.NoError(err)
			assert.Equal(c.r, r, "%v < %v = %v", c.x, c.y, c.r)
		}
	}
}

func ExampleAdd() {
	fmt.Println(Add(1, 2))
	fmt.Println(Add("a", "b"))
	fmt.Println(Add(1+1i, 2+2i))
	fmt.Println(Add(1.1, "x"))
	// Output:
	// 3 <nil>
	// ab <nil>
	// (3+3i) <nil>
	// <nil> incompatible types: float64 + string
}

func ExampleMustInt() {
	MustInt(Add(1, 2))     // 3
	MustInt(Add(1.1, 2.2)) // panics, not int
	MustInt(Add("a", 2))   // panics, not addable
}
