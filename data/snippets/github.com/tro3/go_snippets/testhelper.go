package snippets

import "testing"

type TestHelper struct {
	t *testing.T
}

func NewTestHelper(t *testing.T) TestHelper {
	return TestHelper{t}
}

func (h *TestHelper) Assert(cond bool, msg string) {
	if !cond {
		h.t.Fatalf(msg + "\n")
	}
}

func (h *TestHelper) Assertf(cond bool, format string, args ...interface{}) {
	if !cond {
		h.t.Fatalf(format, args)
	}
}
