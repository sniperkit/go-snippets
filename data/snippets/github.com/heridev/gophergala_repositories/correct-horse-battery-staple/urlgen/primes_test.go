package urlgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPrimesUpTo(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(78498, len(getPrimesUpTo(1000000)))
	assert.Equal(41538, len(getPrimesUpTo(500000)))
	assert.Equal(60238, len(getPrimesUpTo(750000)))
	assert.Equal(67616, len(getPrimesUpTo(849996)))
	assert.Equal(67617, len(getPrimesUpTo(849997)))
}

func TestIsPrime(t *testing.T) {
	assert := assert.New(t)
	assert.True(isPrime(114809))
	assert.True(isPrime(339749))
	assert.True(isPrime(128599))
	assert.True(isPrime(865937))
	assert.True(isPrime(46817))

	assert.False(isPrime(114810))
	assert.False(isPrime(339750))
	assert.False(isPrime(128600))
	assert.False(isPrime(865938))
	assert.False(isPrime(46818))
}

func TestGetPrimeFactors(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, len(getPrimeFactors(114809)))
	assert.Equal(1, len(getPrimeFactors(339749)))
	assert.Equal(1, len(getPrimeFactors(128599)))
	assert.Equal(1, len(getPrimeFactors(865937)))
	assert.Equal(1, len(getPrimeFactors(46817)))

	assert.Equal(5, len(getPrimeFactors(114810)))
	assert.Equal(7, len(getPrimeFactors(339750)))
	assert.Equal(6, len(getPrimeFactors(128600)))
	assert.Equal(3, len(getPrimeFactors(865938)))
	assert.Equal(7, len(getPrimeFactors(46818)))
}

func TestGetUniquePrimeFactors(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, len(getUniquePrimeFactors(114809)))
	assert.Equal(1, len(getUniquePrimeFactors(339749)))
	assert.Equal(1, len(getUniquePrimeFactors(128599)))
	assert.Equal(1, len(getUniquePrimeFactors(865937)))
	assert.Equal(1, len(getUniquePrimeFactors(46817)))

	assert.Equal(5, len(getUniquePrimeFactors(114810)))
	assert.Equal(4, len(getUniquePrimeFactors(339750)))
	assert.Equal(3, len(getUniquePrimeFactors(128600)))
	assert.Equal(3, len(getUniquePrimeFactors(865938)))
	assert.Equal(3, len(getUniquePrimeFactors(46818)))
}
