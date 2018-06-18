package timeutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDuration(t *testing.T) {
	assert := require.New(t)

	for in, out := range map[string]time.Duration{
		`4h`:    time.Duration(4 * time.Hour),
		`4H`:    time.Duration(4 * time.Hour),
		`1d`:    time.Duration(24 * time.Hour),
		`1day`:  time.Duration(24 * time.Hour),
		`1days`: time.Duration(24 * time.Hour),
		`5 years 4 weeks 3 days 2 hours 1 minute`: time.Duration(44546*time.Hour) + time.Minute,
		`1d1h`:  time.Duration(25 * time.Hour),
		`1d 1h`: time.Duration(25 * time.Hour),
	} {
		v, err := ParseDuration(in)
		assert.NoError(err)
		assert.Equal(out, v, fmt.Sprintf("in=%v", in))
	}
}
