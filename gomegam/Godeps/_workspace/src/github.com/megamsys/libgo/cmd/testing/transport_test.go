package testing

import (
	"io/ioutil"
	"gopkg.in/check.v1"
	"net/http"
	"testing"
)

type S struct{}

var _ = check.Suite(S{})

func Test(t *testing.T) {
	check.TestingT(t)
}

func (S) TestTransport(c *check.C) {
	var t http.RoundTripper = &Transport{
		Message: "Ok",
		Status:  http.StatusOK,
		Headers: map[string][]string{"Authorization": {"something"}},
	}
	req, _ := http.NewRequest("GET", "/", nil)
	r, err := t.RoundTrip(req)
	c.Assert(err, check.IsNil)
	c.Assert(r.StatusCode, check.Equals, http.StatusOK)
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	c.Assert(string(b), check.Equals, "Ok")
	c.Assert(r.Header.Get("Authorization"), check.Equals, "something")
}

func (S) TestConditionalTransport(c *check.C) {
	var t http.RoundTripper = &ConditionalTransport{
		Transport: Transport{
			Message: "Ok",
			Status:  http.StatusOK,
		},
		CondFunc: func(req *http.Request) bool {
			return req.URL.Path == "/something"
		},
	}
	req, _ := http.NewRequest("GET", "/something", nil)
	r, err := t.RoundTrip(req)
	c.Assert(err, check.IsNil)
	c.Assert(r.StatusCode, check.Equals, http.StatusOK)
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	c.Assert(string(b), check.Equals, "Ok")
	req, _ = http.NewRequest("GET", "/", nil)
	r, err = t.RoundTrip(req)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "condition failed")
	c.Assert(r.StatusCode, check.Equals, http.StatusInternalServerError)
}
