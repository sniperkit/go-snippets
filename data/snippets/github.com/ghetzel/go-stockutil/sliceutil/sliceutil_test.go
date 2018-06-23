package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	assert := require.New(t)
	var input interface{}

	input = []int{1, 3, 5}
	assert.True(Contains(input, 1))
	assert.True(Contains(input, 3))
	assert.True(Contains(input, 5))
	assert.False(Contains(input, -1))
	assert.False(Contains(input, 2))
	assert.False(Contains(input, -3))
	assert.False(Contains(input, 4))
	assert.False(Contains(input, -5))
	assert.False(Contains([]int{}, 1))
	assert.False(Contains([]int{}, 2))
	assert.False(Contains([]int{}, 0))

	input = []string{"one", "three", "five"}
	assert.True(Contains(input, "one"))
	assert.True(Contains(input, "three"))
	assert.True(Contains(input, "five"))
	assert.False(Contains(input, "One"))
	assert.False(Contains(input, "two"))
	assert.False(Contains(input, "Three"))
	assert.False(Contains(input, "four"))
	assert.False(Contains(input, "Five"))
	assert.False(Contains([]string{}, "one"))
	assert.False(Contains([]string{}, "two"))
	assert.False(Contains([]string{}, ""))
}

func TestAt(t *testing.T) {
	assert := require.New(t)
	var input interface{}
	var out interface{}
	var ok bool

	input = []int{1, 3, 5}
	out, ok = At(input, 0)
	assert.True(ok)
	assert.Equal(1, out)

	out, ok = At(input, 1)
	assert.True(ok)
	assert.Equal(3, out)

	out, ok = At(input, 2)
	assert.True(ok)
	assert.Equal(5, out)

	out, ok = At(input, 99999)
	assert.False(ok)
	assert.Nil(out)
}

func TestLen(t *testing.T) {
	assert := require.New(t)
	var input interface{}

	assert.Zero(Len(nil))
	assert.Zero(Len(input))
	input = []int{1, 3, 5}
	assert.Equal(3, Len(input))
	assert.Equal(3, Len(`123`))
}

func TestGet(t *testing.T) {
	assert := require.New(t)
	var input interface{}

	input = []int{1, 3, 5}
	assert.Equal(1, Get(input, 0))
	assert.Equal(3, Get(input, 1))
	assert.Equal(5, Get(input, 2))
	assert.Nil(Get(input, 99999))
	assert.Nil(Get(nil, 0))
}

func TestFirst(t *testing.T) {
	assert := require.New(t)
	var input interface{}

	assert.Nil(First(nil))
	assert.Nil(First(input))

	input = []int{}
	assert.Nil(First(input))

	input = []int{1, 3, 5}
	assert.Equal(1, First(input))
}

func TestRest(t *testing.T) {
	assert := require.New(t)
	var input interface{}

	assert.Nil(Rest(nil))
	assert.Nil(Rest(input))

	input = []int{1}
	assert.Nil(Rest(input))

	input = []int{1, 3, 5}
	assert.Equal([]interface{}{3, 5}, Rest(input))
}

func TestLast(t *testing.T) {
	assert := require.New(t)
	var input interface{}

	assert.Nil(Last(nil))
	assert.Nil(Last(input))

	input = []int{}
	assert.Nil(Last(input))

	input = []int{1, 3, 5}
	assert.Equal(5, Last(input))
}

func TestContainsString(t *testing.T) {
	assert := require.New(t)

	input := []string{"one", "three", "five"}

	assert.True(ContainsString(input, "one"))
	assert.True(ContainsString(input, "three"))
	assert.True(ContainsString(input, "five"))
	assert.False(ContainsString(input, "One"))
	assert.False(ContainsString(input, "two"))
	assert.False(ContainsString(input, "Three"))
	assert.False(ContainsString(input, "four"))
	assert.False(ContainsString(input, "Five"))
	assert.False(ContainsString([]string{}, "one"))
	assert.False(ContainsString([]string{}, "two"))
	assert.False(ContainsString([]string{}, ""))
}

func TestContainsAnyString(t *testing.T) {
	assert := require.New(t)

	input := []string{"one", "three", "five"}
	any := []string{"one", "two", "four"}

	assert.True(ContainsAnyString(input, any...))
	assert.False(ContainsAnyString(input))
	assert.False(ContainsAnyString([]string{}, "one"))
	assert.False(ContainsAnyString([]string{}, "two"))
	assert.False(ContainsAnyString([]string{}, ""))
	assert.False(ContainsAnyString(input, []string{"six", "seven"}...))
}

func TestContainsAllStrings(t *testing.T) {
	assert := require.New(t)

	input := []string{"one", "three", "five"}

	assert.True(ContainsAllStrings(input, "one"))
	assert.True(ContainsAllStrings(input, "three"))
	assert.True(ContainsAllStrings(input, "five"))
	assert.True(ContainsAllStrings(input, "one", "three"))
	assert.True(ContainsAllStrings(input, "one", "five"))
	assert.True(ContainsAllStrings(input, "one", "three", "five"))
	assert.False(ContainsAllStrings(input, "one", "four"))
	assert.True(ContainsAllStrings(input))
}

func TestCompact(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Compact(nil))

	assert.Equal([]interface{}{
		0, 1, 2, 3,
	}, Compact([]interface{}{
		0, 1, 2, 3,
	}))

	assert.Equal([]interface{}{
		1, 3, 5,
	}, Compact([]interface{}{
		nil, 1, nil, 3, nil, 5,
	}))

	assert.Equal([]interface{}{
		`one`, `three`, ` `, `five`,
	}, Compact([]interface{}{
		`one`, ``, `three`, ``, ` `, `five`,
	}))

	assert.Equal([]interface{}{
		[]int{1, 2, 3},
	}, Compact([]interface{}{
		nil, []string{}, []int{1, 2, 3}, map[string]bool{},
	}))
}

func TestCompactString(t *testing.T) {
	assert := require.New(t)

	assert.Nil(CompactString(nil))

	assert.Equal([]string{
		`one`, `three`, `five`,
	}, CompactString([]string{
		`one`, `three`, `five`,
	}))

	assert.Equal([]string{
		`one`, `three`, ` `, `five`,
	}, CompactString([]string{
		`one`, ``, `three`, ``, ` `, `five`,
	}))
}

func TestStringify(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Stringify(nil))

	assert.Equal([]string{
		`0`, `1`, `2`,
	}, Stringify([]interface{}{
		0, 1, 2,
	}))

	assert.Equal([]string{
		`0.5`, `0.55`, `0.555`, `0.555001`,
	}, Stringify([]interface{}{
		0.5, 0.55, 0.55500, 0.555001,
	}))

	assert.Equal([]string{
		`true`, `true`, `false`,
	}, Stringify([]interface{}{
		true, true, false,
	}))
}

func TestOr(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Or())
	assert.Nil(Or(nil))
	assert.Equal(1, Or(0, 1, 0, 2, 0, 3, 4, 5, 6))
	assert.Equal(true, Or(false, false, true))
	assert.Equal(`one`, Or(`one`))
	assert.Equal(4.0, Or(nil, ``, false, 0, 4.0))
	assert.Nil(Or(false, false, false))
	assert.Nil(Or(0, 0, 0))

	assert.Equal(`three`, Or(``, ``, `three`))

	type testStruct struct {
		name string
	}

	assert.Equal(testStruct{`three`}, Or(testStruct{}, testStruct{}, testStruct{`three`}))
}

func TestOrString(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, OrString())
	assert.Equal(``, OrString(``))

	assert.Equal(`one`, OrString(`one`))
	assert.Equal(`two`, OrString(``, `two`, ``, `three`))
}

func TestEach(t *testing.T) {
	assert := require.New(t)

	assert.Nil(Each(nil, nil))

	assert.Nil(Each([]string{`one`, `two`, `three`}, func(i int, v interface{}) error {
		return Stop
	}))

	count := 0
	assert.Nil(Each([]string{`one`, `two`, `three`}, func(i int, v interface{}) error {
		if v.(string) == `two` {
			return Stop
		} else {
			count += 1
			return nil
		}
	}))

	assert.Equal(1, count)
}

func TestUnique(t *testing.T) {
	assert := require.New(t)

	assert.Equal([]interface{}{`one`, `two`, `three`}, Unique([]string{`one`, `one`, `two`, `three`}))
	assert.Equal([]interface{}{1, 2, 3}, Unique([]int{1, 2, 2, 3}))
	assert.NotEqual([]interface{}{1, 2, 3}, Unique([]int64{1, 2, 2, 3}))
}

func TestMap(t *testing.T) {
	assert := require.New(t)

	assert.Equal(
		[]interface{}{10, 20, 30},
		Map([]int{1, 2, 3}, func(_ int, v interface{}) interface{} {
			return v.(int) * 10
		}),
	)

	assert.Equal(
		[]interface{}{true, false, true},
		Map([]bool{false, true, false}, func(_ int, v interface{}) interface{} {
			return !v.(bool)
		}),
	)
}

func TestMapString(t *testing.T) {
	assert := require.New(t)

	assert.Equal(
		[]string{`1-1thousand`, `2-1thousand`, `3-1thousand`},
		MapString([]int{1, 2, 3}, func(_ int, v string) string {
			return v + `-1thousand`
		}),
	)

	assert.Equal(
		[]string{`first`, `third`, `fifth`},
		CompactString(MapString([]string{`first`, `second`, `third`, `fourth`, `fifth`}, func(_ int, v string) string {
			switch v {
			case `second`, `fourth`:
				return ``
			default:
				return v
			}
		})),
	)
}

func TestChunks(t *testing.T) {
	assert := require.New(t)
	input := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23}

	assert.Equal([][]interface{}{
		[]interface{}{1},
		[]interface{}{3},
		[]interface{}{5},
		[]interface{}{7},
		[]interface{}{9},
		[]interface{}{11},
		[]interface{}{13},
		[]interface{}{15},
		[]interface{}{17},
		[]interface{}{19},
		[]interface{}{21},
		[]interface{}{23},
	}, Chunks(input, 1))

	assert.Equal([][]interface{}{
		[]interface{}{1, 3},
		[]interface{}{5, 7},
		[]interface{}{9, 11},
		[]interface{}{13, 15},
		[]interface{}{17, 19},
		[]interface{}{21, 23},
	}, Chunks(input, 2))

	assert.Equal([][]interface{}{
		[]interface{}{1, 3, 5},
		[]interface{}{7, 9, 11},
		[]interface{}{13, 15, 17},
		[]interface{}{19, 21, 23},
	}, Chunks(input, 3))

	assert.Equal([][]interface{}{
		[]interface{}{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23},
	}, Chunks(input, 1000))
}

func TestFlatten(t *testing.T) {
	assert := require.New(t)

	assert.Equal([]interface{}{`one`, `two`, `three`}, Flatten([]string{`one`, `two`, `three`}))
	assert.Equal([]interface{}{`one`, `two`, `three`}, Flatten([]interface{}{[]string{`one`, `two`}, `three`}))
	assert.Equal([]interface{}{`one`, `two`, `three`}, Flatten([]interface{}{[]string{`one`}, []string{`two`}, []string{`three`}}))
}
