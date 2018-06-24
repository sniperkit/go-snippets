package sqlutil

import (
	"testing"
)

func TestInitTestSuite(t *testing.T) {
	cases := []struct {
		o   *TestSuiteOpts
		ts  *TestSuite
		err error
	}{}

	for i, c := range cases {
		ts, err := InitTestSuite(c.o)
		if err != c.err {
			t.Errorf("case %d error mismatch: %s != %s", i, c.err, err)
			return
		}
		if ts == nil && c.ts != nil || ts != nil && c.ts == nil {
			t.Errorf("case %d TestSuite nil mismatch: %s != %s", i, c.ts, ts)
			return
		}
	}
}
