package app

import (
	"errors"
	"launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) TestAppLifecycleError(c *gocheck.C) {
	e := AppLifecycleError{app: "myapp", Err: errors.New("failure in app")}
	expected := `gulpd failed to apply the lifecle to the app "myapp": failure in app`
	c.Assert(e.Error(), gocheck.Equals, expected)
}
