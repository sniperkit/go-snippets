package rewrite

import (
	"github.com/datatogether/warc"
	"testing"
)

func TestWarcRecordRewriter(t *testing.T) {
	cases := []struct {
		in, out *warc.Record
		err     string
	}{}

	rw := NewWarcRecordRewriter("http://a.com")
	for i, c := range cases {
		got, err := rw.RewriteRecord(c.in)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d RewriteRecord mismatch: expected: %s, got: %s", i, c.err, err)
			continue
		}

		if c.out.Type.String() != got.Type.String() {
			t.Errorf("case %d returned type mismatch. expected: %s, got: %s", i, c.out.Type, got.Type)
			continue
		}
	}
}
