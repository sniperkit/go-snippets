package ffi

import (
	"io"
	"testing"
)

func errorMismatch(err error, expect string) bool {
	return (err == nil && expect != "") || (err != nil && expect == "") || (err != nil && err.Error() != expect)
}

func TestFilenameFromUrlString(t *testing.T) {
	cases := []struct {
		in, out, err string
	}{
		{"https://archivers.space/files/not/a/path/image.jpg", "image.jpg", ""},
		{"https://www.epa.gov/sites/production/files/2017-02/tscainv_feb2017_csv.zip", "tscainv_feb2017_csv.zip", ""},
		{"https://www.census.gov/topics/business/retail-trade.html", "retail-trade.html", ""},
		{"https://www.census.gov/topics/education/educational-attainment.html", "educational-attainment.html", ""},
		{"https://www.epa.gov/caddis-vol2", "caddis-vol2", ""},
		{"https://www.epa.gov/caddis-vol4", "caddis-vol4", ""},
		{"https://www.census.gov/topics/population/language-use.html", "language-use.html", ""},
		{"https://www.census.gov/topics/business/classification-codes.html", "classification-codes.html", ""},
		{"https://www.census.gov/topics/business/business-help.html", "business-help.html", ""},
		{"https://www.epa.gov/sites/production/files/2017-02/", "2017-02", ""},
		{"https://www.epa.gov/caddis-vol1", "caddis-vol1", ""},
		{"https://www.epa.gov/sites/production/files/2017-02", "2017-02", ""},
		{"https://www.census.gov/topics/population/age-and-sex.html", "age-and-sex.html", ""},
		{"https://census.gov/library/publications/2009/compendia/statab/129ed/population.html#NAV_1529793603_0_accd", "population.html", ""},
		{"https://www.census.gov/topics/education/school-enrollment.html?query=param", "school-enrollment.html", ""},
		{"https://www.census.gov/topics/housing/housing-vacancies.html", "housing-vacancies.html", ""},
		// {"https://qri.io/files/not/a/path/image", "", nil},
	}

	for i, c := range cases {
		got, err := FilenameFromUrlString(c.in)
		if errorMismatch(err, c.err) {
			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
		}
		if got != c.out {
			t.Errorf("case %d output mismatch. expected: %s, got: %s", i, c.out, got)
		}
	}
}

func TestFilenameMimeType(t *testing.T) {
	cases := []struct {
		in, out, err string
	}{
		{"", "", "unrecognized MIME-Type: ''"},
	}

	for i, c := range cases {
		got, err := FilenameMimeType(c.in)
		if errorMismatch(err, c.err) {
			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
		}
		if got != c.out {
			t.Errorf("case %d output mismatch. expected: %s, got: %s", i, c.out, got)
		}
	}
}

func TestSetExtension(t *testing.T) {
	cases := []struct {
		fn, mt, out, err string
	}{
		{"apples.html", "application/pdf", "apples.pdf", ""},
		{"apples.html", "invalid/blech", "apples.html", "unrecognized MIME-Type: 'invalid/blech'"},
	}

	for i, c := range cases {
		got, err := SetExtension(c.fn, c.mt)
		if errorMismatch(err, c.err) {
			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
		}
		if got != c.out {
			t.Errorf("case %d output mismatch. expected: %s, got: %s", i, c.out, got)
		}
	}
}

func TestMimeTypeExtension(t *testing.T) {
	cases := []struct {
		in, out, err string
	}{
		{"", "", "unrecognized MIME-Type: ''"},
	}

	for i, c := range cases {
		got, err := MimeTypeExtension(c.in)
		if errorMismatch(err, c.err) {
			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
		}
		if got != c.out {
			t.Errorf("case %d output mismatch. expected: %s, got: %s", i, c.out, got)
		}
	}
}

func TestExtensionMimeType(t *testing.T) {
	cases := []struct {
		in, out, err string
	}{
		{"", "", "unrecognized extension: ''"},
	}

	for i, c := range cases {
		got, err := ExtensionMimeType(c.in)
		if errorMismatch(err, c.err) {
			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
		}
		if got != c.out {
			t.Errorf("case %d output mismatch. expected: %s, got: %s", i, c.out, got)
		}
	}
}

func TestMimeType(t *testing.T) {
	cases := []struct {
		in  io.Reader
		out string
		err error
	}{
	// {"", "", nil},
	}

	for i, c := range cases {
		got, err := MimeType(c.in)
		if err != c.err {
			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
		}
		if got != c.out {
			t.Errorf("case %d output mismatch. expected: %s, got: %s", i, c.out, got)
		}
	}
}
