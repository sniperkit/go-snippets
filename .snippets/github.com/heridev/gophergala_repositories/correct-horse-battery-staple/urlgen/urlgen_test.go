package urlgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlGen(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(len(alphabet), len(shuffledAlphabet))
}

func TestEncodeDecode(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(4, decode(encode(4, 1)))
	assert.Equal(4, decode(encode(4, 3)))
	assert.Equal(3, len(encode(4, 3)))
}

func TestTokenizing(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(GetTokenFromId(4), GetTokenFromId(4))
	assert.Equal(4, GetIdFromToken(GetTokenFromId(4)))
	assert.NotEqual(GetTokenFromId(4), GetTokenFromId(5))
}
