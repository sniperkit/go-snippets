package main

import "testing"

func TestParseTemplatesDir(t *testing.T) {
	cases := []struct {
		dir        string
		validCount int
		throwsErr  bool
	}{
		{"./testing/templates", 2, false},
	}

	for _, test := range cases {
		templates, err := ParseTemplatesDir(test.dir)
		if test.validCount != len(templates) {
			t.Errorf("ParseTemplatesDir expected %d valid templates but found %d in dir: %s", test.validCount, len(templates), test.dir)
		}

		if test.throwsErr {
			if err == nil {
				t.Errorf("ParseTemplatesDir expected an err but didn't for dir: %s", test.dir)
			}
		} else {
			if err != nil {
				t.Errorf("ParseTemplatesDir should not have thrown an error but did for directory: %s, %q", test.dir, err.Error())
			}
		}
	}
}
