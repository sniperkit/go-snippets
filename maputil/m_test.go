package maputil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestM(t *testing.T) {
	assert := require.New(t)
	input := M(map[string]interface{}{
		`first`: true,
		`second`: map[string]interface{}{
			`s1`:     `test`,
			`values`: []int{1, 2, 3, 4},
			`truthy`: `True`,
			`strnum`: `42`,
			`then`:   `2006-01-02`,
		},
		`now`:    time.Now(),
		`third`:  3.1415,
		`fourth`: 42,
	})

	assert.Equal(``, M(nil).String(`lol`))
	assert.False(M(nil).Bool(`lol`))
	assert.Equal(int64(0), M(nil).Int(`lol`))
	assert.Equal(float64(0), M(nil).Float(`lol`))
	assert.Len(M(nil).Slice(`lol`), 0)
	assert.Nil(M(nil).Auto(`second.strnum`))
	assert.Zero(M(nil).Time(`now`))

	assert.Equal(`test`, input.String(`second.s1`))
	assert.True(input.Bool(`first`))
	assert.True(input.Bool(`second.truthy`))
	assert.False(input.Bool(`second.s1`))
	assert.Equal(3.1415, input.Float(`third`))
	assert.Equal(int64(3), input.Int(`third`))
	assert.Equal(int64(42), input.Int(`fourth`))
	assert.Equal(int64(3), input.Int(`second.values.2`))
	assert.Equal(int64(0), input.Int(`second.values.99`))
	assert.Equal(float64(42), input.Float(`fourth`))
	assert.Len(input.Slice(`second.values`), 4)
	assert.Equal(int64(42), input.Auto(`second.strnum`))
	assert.Equal(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), input.Time(`second.then`))
}

func TestMSet(t *testing.T) {
	assert := require.New(t)
	input := M(nil)

	assert.Equal(``, input.String(`lol`))

	assert.Equal(`2funny4me`, input.Set(`lol`, `2funny4me`).String())

	assert.Equal(`2funny4me`, input.String(`lol`))
}
