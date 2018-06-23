package rest

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTripper(t *testing.T) {
	errIn := errors.New("TEST")
	respIn := &http.Response{}
	rt := RoundTripper{Response: respIn, Error: errIn}
	respOut, errOut := rt.RoundTrip(&http.Request{})
	assert.Equal(t, respIn, respOut)
	assert.Equal(t, errIn, errOut)
}
