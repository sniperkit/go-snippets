package stringutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	assert := require.New(t)
	input := []byte{
		0x01, 0x02, 0x03, 0x01,
		0x02, 0x03, 0x01, 0x02,
		0x03, 0x01, 0x02, 0x03,
		0x01, 0x02, 0x03, 0x01,
	}

	uuid, err := UuidFromBytes(input)

	assert.NoError(err)

	assert.Equal(`01020301-0203-0102-0301-020301020301`, uuid.String())
	assert.Equal(input, uuid.Bytes())
}
