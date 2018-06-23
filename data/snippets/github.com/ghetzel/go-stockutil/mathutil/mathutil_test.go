package mathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClamp(t *testing.T) {
	assert := assert.New(t)

	assert.True(Clamp(0, 0, 0) == 0)
	assert.True(Clamp(0, -1, 1) == 0)
	assert.True(Clamp(-1, -1, 1) == -1)
	assert.True(Clamp(-0.9999, -1, 1) == -0.9999)
	assert.True(Clamp(-1.0001, -1, 1) == -1)
	assert.True(Clamp(-2, -1, 1) == -1)
	assert.True(Clamp(1, -1, 1) == 1)
	assert.True(Clamp(0.9999, -1, 1) == 0.9999)
	assert.True(Clamp(1.0001, -1, 1) == 1)
	assert.True(Clamp(2, -1, 1) == 1)
}

func TestRound(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(float64(1), Round(1.499999999999999))
	assert.Equal(float64(2), Round(1.5))
	assert.Equal(float64(2), Round(1.999999999999999))
	assert.Equal(float64(2), Round(1.999999999999999))

	assert.Equal(float64(2), RoundPlaces(2.49, 0))
	assert.Equal(float64(2.5), RoundPlaces(2.49, 1))
	assert.Equal(float64(2.49), RoundPlaces(2.49, 2))
	assert.Equal(float64(2.49), RoundPlaces(2.490000000049, 10))
	assert.Equal(float64(2.4900000001), RoundPlaces(2.49000000005, 10))
}
