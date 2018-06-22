package cmd

import (
	"bytes"
	ttesting "github.com/megamsys/libgo/cmd/testing"
	"gopkg.in/check.v1"
	"net/http"
)

func (s *S) TestShouldSetCloseToTrue(c *check.C) {
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, check.IsNil)
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: "OK",
	}
	client := NewClient(&http.Client{Transport: &transport}, nil, manager)
	client.Do(request)
	c.Assert(request.Close, check.Equals, true)
}





