package cdxj

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/datatogether/warc"
)

const eg = `!OpenWayback-CDXJ 1.0
(com,cnn,)/world> 2015-09-03T13:27:52Z response {"a":0,"b":"b","c":false}
(uk,ac,rpms,)/> 2015-09-03T13:27:52Z request {"frequency":241,"spread":3}
(uk,co,bbc,)/images> 2015-09-03T13:27:52Z response {"frequency":725,"spread":1}
`

const eg2 = `!OpenWayback-CDXJ 1.0
((com,reddit,)> 2015-09-03T13:27:52Z request {}`

func TestReader(t *testing.T) {
	cases := []struct {
		in     string
		expect []*Record
		err    string
	}{
		{eg, []*Record{
			&Record{"cnn.com/world", time.Date(2015, time.September, 3, 13, 27, 52, 0, time.UTC), warc.RecordTypeResponse, map[string]interface{}{"a": 0, "b": "b", "c": false}},
			&Record{"rpms.ac.uk/", time.Date(2015, time.September, 3, 13, 27, 52, 0, time.UTC), warc.RecordTypeRequest, map[string]interface{}{"frequency": 241, "spread": 3}},
			&Record{"bbc.co.uk/images", time.Date(2015, time.September, 3, 13, 27, 52, 0, time.UTC), warc.RecordTypeResponse, map[string]interface{}{"frequency": 725, "spread": 1}},
		}, ""},
		{eg2, []*Record{
			&Record{"reddit.com", time.Date(2015, time.September, 3, 13, 27, 52, 0, time.UTC), warc.RecordTypeRequest, map[string]interface{}{}},
		}, ""},
	}

	for i, c := range cases {
		records, err := NewReader(strings.NewReader(c.in)).ReadAll()
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err.Error())
			continue
		}

		if err := CompareRecordSlices(c.expect, records); err != nil {
			t.Errorf("case %d record slice mistmatch: %s", i, err.Error())
			continue
		}

	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		input string
		err   string
	}{
		{"", "invalid format, missing cdxj header"},
		{eg, ""},
	}

	for i, c := range cases {
		err := Validate(bytes.NewBufferString(c.input))
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
			continue
		}
	}
}
