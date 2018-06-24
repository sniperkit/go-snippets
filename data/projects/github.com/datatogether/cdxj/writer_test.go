package cdxj

import (
	"bytes"
	"testing"
	"time"

	"github.com/datatogether/warc"
)

func TestWriter(t *testing.T) {
	records := []*Record{
		&Record{"cnn.com/world", time.Date(2015, time.September, 3, 13, 27, 52, 0, time.UTC), warc.RecordTypeResponse, map[string]interface{}{"a": 0, "b": "b", "c": false}},
		&Record{"rpms.ac.uk/", time.Date(2015, time.September, 3, 13, 27, 52, 0, time.UTC), warc.RecordTypeRequest, map[string]interface{}{"frequency": 241, "spread": 3}},
		&Record{"bbc.co.uk/images", time.Date(2015, time.September, 3, 13, 27, 52, 0, time.UTC), warc.RecordTypeResponse, map[string]interface{}{"frequency": 725, "spread": 1}},
	}
	buf := &bytes.Buffer{}

	w := NewWriter(buf)
	for i, rec := range records {
		if err := w.Write(rec); err != nil {
			t.Errorf("error writing record %d: %s", i, err.Error())
			return
		}
	}

	if err := w.Close(); err != nil {
		t.Errorf("close error: %s", err.Error())
		return
	}

	if buf.String() != eg {
		t.Errorf("result mismatch expected:\n%s\ngot:\n%s", eg, buf.String())
	}
}
