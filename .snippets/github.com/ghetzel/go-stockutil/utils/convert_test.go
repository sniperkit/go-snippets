package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectConvertType(t *testing.T) {
	assert := require.New(t)

	assert.Equal(Nil, DetectConvertType(nil))
	assert.Equal(Nil, DetectConvertType(``))

	assert.Equal(String, DetectConvertType(`07753`))
	assert.Equal(String, DetectConvertType(`test`))

	assert.Equal(Time, DetectConvertType(`2015-05-01 00:15:16`))
	assert.Equal(Time, DetectConvertType(`Fri May 1 00:15:16 UTC 2015`))
	assert.Equal(Time, DetectConvertType(`Fri May 01 00:15:16 +0000 2015`))
	assert.Equal(Time, DetectConvertType(`01 May 15 00:15 UTC`))
	assert.Equal(Time, DetectConvertType(`01 May 15 00:15 +0000`))
	assert.Equal(Time, DetectConvertType(`Friday, 01-May-15 00:15:16 UTC`))

	assert.Equal(Integer, DetectConvertType(`1136239445`))
	assert.Equal(Integer, DetectConvertType(`1136239445000`))
	assert.Equal(Integer, DetectConvertType(`1136239445000000`))
	assert.Equal(Integer, DetectConvertType(`0`))
	assert.Equal(Integer, DetectConvertType(`1`))
	assert.Equal(Integer, DetectConvertType(`17753`))

	assert.Equal(Float, DetectConvertType(`0.0`))
	assert.Equal(Float, DetectConvertType(`3.1415`))
	assert.Equal(Float, DetectConvertType(`3.0001`))
	assert.Equal(Float, DetectConvertType(`3.1000`))
}
