package bytes

import "testing"

var fmtTests = []struct {
	Value     int64
	Base      int
	Precision int
	Result    string
}{
	{1000, 10, 0, "1KB"},
	{1024, 10, 2, "1.02KB"},
	{1024, 2, 1, "1.0KiB"},
	{56000, 10, 2, "56.00KB"},        // 1 sec of 56K modem
	{1024000, 2, 2, "1000.00KiB"},    // A "1MB" floppy disk is 1000KiB
	{52428800, 2, 0, "50MiB"},        // A "50MB" file on a CD is 50MiB
	{1300000000, 2, 2, "1.21GiB"},    // A 1.3GB file on a DVD is 1.3GB
	{300000000000, 10, 1, "300.0GB"}, // A 300GB hard disk is 300GB
	{8589934592, 2, 2, "8.00GiB"},    // "8GB" of RAM is 8GiB
}

func TestFormat(t *testing.T) {
	for _, test := range fmtTests {
		ts := Format(test.Value, test.Base, test.Precision)
		if ts != test.Result {
			t.Errorf("%d base %d precision %d should be %s, was %s", test.Value,
				test.Base, test.Precision, test.Result, ts)
		}
	}
}

var parseTests = []struct {
	Value  string
	Result int64
}{
	{"1KB", 1000},
	{"1KiB", 1024},
	{"6GiB", 6 * 1024 * 1024 * 1024},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		x, err := ParseBytes(test.Value)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if x != test.Result {
			t.Errorf("failed to parse %s correctly", test.Value)
		}
	}
}

func TestSplit(t *testing.T) {
	x, y := split("10MB")
	if x != "10" || y != "MB" {
		t.Errorf("failed to split 10MB")
	}
	x, y = split("5 GiB")
	if x != "5" || y != "GiB" {
		t.Errorf("failed to split 5 GiB")
	}
}

func TestErrors(t *testing.T) {
	_, err := ParseBytes("1ZB")
	if err == nil {
		t.Errorf("failed to flag range error")
	}
}
