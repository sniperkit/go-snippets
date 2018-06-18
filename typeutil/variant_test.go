package typeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVariant(t *testing.T) {
	assert := require.New(t)

	assert.Equal(`test`, Variant{`test`}.String())
	assert.True(Variant{`True`}.Bool())
	assert.True(Variant{`true`}.Bool())
	assert.True(Variant{`TRUE`}.Bool())
	assert.True(Variant{`1`}.Bool())
	assert.False(Variant{`False`}.Bool())
	assert.False(Variant{`false`}.Bool())
	assert.False(Variant{`0`}.Bool())
	assert.False(Variant{`dennis`}.Bool())
	assert.Equal(int64(1), Variant{1}.Int())
	assert.Equal(int64(1), Variant{1.9}.Int())
	assert.Equal(float64(1.9), Variant{1.9}.Float())
	assert.True(time.Unix(1500000000, 0).Equal(Variant{1500000000}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`1500000000`}.Time()))
	assert.True(time.Unix(1500000000, 0).Equal(Variant{`2017-07-14T02:40:00Z`}.Time()))
	assert.Equal([]byte{0x74, 0x65, 0x73, 0x74}, Variant{`test`}.Bytes())

	assert.Equal(map[Variant]Variant{
		V(`test`):  V(1),
		V(`other`): V(2.4),
	}, V(map[string]interface{}{
		`test`:  1,
		`other`: 2.4,
	}).Map())
}
