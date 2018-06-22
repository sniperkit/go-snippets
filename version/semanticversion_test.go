package version

import (
	"testing"
)

func prstr(s string) PreReleaseValue {
	return PreReleaseValue{s, 0, false}
}

func prnum(i uint64) PreReleaseValue {
	return PreReleaseValue{"", i, true}
}

type formatTest struct {
	v      *SemanticVersion
	result string
}

var formatTests = []formatTest{
	{&SemanticVersion{1, 2, 3, nil, nil}, "1.2.3"},
	{&SemanticVersion{0, 0, 1, nil, nil}, "0.0.1"},
	{&SemanticVersion{0, 0, 1, []PreReleaseValue{prstr("alpha"), prstr("preview")}, []string{"123", "456"}}, "0.0.1-alpha.preview+123.456"},
	{&SemanticVersion{1, 2, 3, []PreReleaseValue{prstr("alpha"), prnum(1)}, []string{"123", "456"}}, "1.2.3-alpha.1+123.456"},
	{&SemanticVersion{1, 2, 3, []PreReleaseValue{prstr("alpha"), prnum(1)}, nil}, "1.2.3-alpha.1"},
	{&SemanticVersion{1, 2, 3, nil, []string{"123", "456"}}, "1.2.3+123.456"},
	// Prereleases and build metadata hyphens
	{&SemanticVersion{1, 2, 3, []PreReleaseValue{prstr("alpha"), prstr("b-eta")}, []string{"123", "b-uild"}}, "1.2.3-alpha.b-eta+123.b-uild"},
	{&SemanticVersion{1, 2, 3, nil, []string{"123", "b-uild"}}, "1.2.3+123.b-uild"},
	{&SemanticVersion{1, 2, 3, []PreReleaseValue{prstr("alpha"), prstr("b-eta")}, nil}, "1.2.3-alpha.b-eta"},
}

var tolerantFormatTests = []formatTest{
	{&SemanticVersion{1, 2, 3, nil, nil}, "v1.2.3"},
	{&SemanticVersion{1, 2, 3, nil, nil}, "	1.2.3 "},
	{&SemanticVersion{1, 2, 0, nil, nil}, "1.2"},
	{&SemanticVersion{1, 0, 0, nil, nil}, "1"},
}

func TestStringer(t *testing.T) {
	for _, test := range formatTests {
		if res := test.v.String(); res != test.result {
			t.Errorf("Stringer, expected %q but got %q", test.result, res)
		}
	}
}

func TestParse(t *testing.T) {
	for _, test := range formatTests {
		if v, err := Parse(test.result); err != nil {
			t.Errorf("Error parsing %q: %q", test.result, err)
		} else if comp := v.Compare(test.v); comp != 0 {
			t.Errorf("Parsing, expected %q but got %q, comp: %d ", test.v, v, comp)
		} else if err := v.Validate(); err != nil {
			t.Errorf("Error validating parsed version %q: %q", test.v, err)
		}
	}
}

func TestParseTolerant(t *testing.T) {
	for _, test := range tolerantFormatTests {
		if v, err := ParseTolerant(test.result); err != nil {
			t.Errorf("Error parsing %q: %q", test.result, err)
		} else if comp := v.Compare(test.v); comp != 0 {
			t.Errorf("Parsing, expected %q but got %q, comp: %d ", test.v, v, comp)
		} else if err := v.Validate(); err != nil {
			t.Errorf("Error validating parsed version %q: %q", test.v, err)
		}
	}
}

func TestMustParse(t *testing.T) {
	_ = MustParse("32.2.1-alpha")
}

func TestMustParse_panic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Errorf("Should have panicked")
		}
	}()
	_ = MustParse("invalid version")
}

func TestValidate(t *testing.T) {
	for _, test := range formatTests {
		if err := test.v.Validate(); err != nil {
			t.Errorf("Error validating %q: %q", test.v, err)
		}
	}
}

type compareTest struct {
	v1     Version
	v2     Version
	result int
}

var compareTests = []compareTest{
	{&SemanticVersion{1, 0, 0, nil, nil}, &SemanticVersion{1, 0, 0, nil, nil}, 0},
	{&SemanticVersion{2, 0, 0, nil, nil}, &SemanticVersion{1, 0, 0, nil, nil}, 1},
	{&SemanticVersion{0, 1, 0, nil, nil}, &SemanticVersion{0, 1, 0, nil, nil}, 0},
	{&SemanticVersion{0, 2, 0, nil, nil}, &SemanticVersion{0, 1, 0, nil, nil}, 1},
	{&SemanticVersion{0, 0, 1, nil, nil}, &SemanticVersion{0, 0, 1, nil, nil}, 0},
	{&SemanticVersion{0, 0, 2, nil, nil}, &SemanticVersion{0, 0, 1, nil, nil}, 1},
	{&SemanticVersion{1, 2, 3, nil, nil}, &SemanticVersion{1, 2, 3, nil, nil}, 0},
	{&SemanticVersion{2, 2, 4, nil, nil}, &SemanticVersion{1, 2, 4, nil, nil}, 1},
	{&SemanticVersion{1, 3, 3, nil, nil}, &SemanticVersion{1, 2, 3, nil, nil}, 1},
	{&SemanticVersion{1, 2, 4, nil, nil}, &SemanticVersion{1, 2, 3, nil, nil}, 1},

	// Spec Examples #11
	{&SemanticVersion{1, 0, 0, nil, nil}, &SemanticVersion{2, 0, 0, nil, nil}, -1},
	{&SemanticVersion{2, 0, 0, nil, nil}, &SemanticVersion{2, 1, 0, nil, nil}, -1},
	{&SemanticVersion{2, 1, 0, nil, nil}, &SemanticVersion{2, 1, 1, nil, nil}, -1},

	// Spec Examples #9
	{&SemanticVersion{1, 0, 0, nil, nil}, &SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("alpha")}, nil}, 1},
	{&SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("alpha")}, nil}, &SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("alpha"), prnum(1)}, nil}, -1},
	{&SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("alpha"), prnum(1)}, nil}, &SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("alpha"), prstr("beta")}, nil}, -1},
	{&SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("alpha"), prstr("beta")}, nil}, &SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("beta")}, nil}, -1},
	{&SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("beta")}, nil}, &SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("beta"), prnum(2)}, nil}, -1},
	{&SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("beta"), prnum(2)}, nil}, &SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("beta"), prnum(11)}, nil}, -1},
	{&SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("beta"), prnum(11)}, nil}, &SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("rc"), prnum(1)}, nil}, -1},
	{&SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("rc"), prnum(1)}, nil}, &SemanticVersion{1, 0, 0, nil, nil}, -1},

	// Ignore Build metadata
	{&SemanticVersion{1, 0, 0, nil, []string{"1", "2", "3"}}, &SemanticVersion{1, 0, 0, nil, nil}, 0},
}

func TestCompare(t *testing.T) {
	for _, test := range compareTests {
		if res := test.v1.Compare(test.v2); res != test.result {
			t.Errorf("Comparing %q : %q, expected %d but got %d", test.v1, test.v2, test.result, res)
		}
		//Test counterpart
		if res := test.v2.Compare(test.v1); res != -test.result {
			t.Errorf("Comparing %q : %q, expected %d but got %d", test.v2, test.v1, -test.result, res)
		}
	}
}

type wrongformatTest struct {
	v   Version
	str string
}

var wrongformatTests = []wrongformatTest{
	{nil, ""},
	{nil, "."},
	{nil, "1."},
	{nil, ".1"},
	{nil, "a.b.c"},
	{nil, "1.a.b"},
	{nil, "1.1.a"},
	{nil, "1.a.1"},
	{nil, "a.1.1"},
	{nil, ".."},
	{nil, "1.."},
	{nil, "1.1."},
	{nil, "1..1"},
	{nil, "1.1.+123"},
	{nil, "1.1.-beta"},
	{nil, "-1.1.1"},
	{nil, "1.-1.1"},
	{nil, "1.1.-1"},
	// giant numbers
	{nil, "20000000000000000000.1.1"},
	{nil, "1.20000000000000000000.1"},
	{nil, "1.1.20000000000000000000"},
	{nil, "1.1.1-20000000000000000000"},
	// Leading zeroes
	{nil, "01.1.1"},
	{nil, "001.1.1"},
	{nil, "1.01.1"},
	{nil, "1.001.1"},
	{nil, "1.1.01"},
	{nil, "1.1.001"},
	{nil, "1.1.1-01"},
	{nil, "1.1.1-001"},
	{nil, "1.1.1-beta.01"},
	{nil, "1.1.1-beta.001"},
	{&SemanticVersion{0, 0, 0, []PreReleaseValue{prstr("!")}, nil}, "0.0.0-!"},
	{&SemanticVersion{0, 0, 0, nil, []string{"!"}}, "0.0.0+!"},
	// empty prversion
	{&SemanticVersion{0, 0, 0, []PreReleaseValue{prstr(""), prstr("alpha")}, nil}, "0.0.0-.alpha"},
	// empty build meta data
	{&SemanticVersion{0, 0, 0, []PreReleaseValue{prstr("alpha")}, []string{""}}, "0.0.0-alpha+"},
	{&SemanticVersion{0, 0, 0, []PreReleaseValue{prstr("alpha")}, []string{"test", ""}}, "0.0.0-alpha+test."},
}

func TestWrongFormat(t *testing.T) {
	for _, test := range wrongformatTests {

		if res, err := Parse(test.str); err == nil {
			t.Errorf("Parsing wrong format version %q, expected error but got %q", test.str, res)
		}

		if test.v != nil {
			if err := test.v.Validate(); err == nil {
				t.Errorf("Validating wrong format version %q (%q), expected error", test.v, test.str)
			}
		}
	}
}

var wrongTolerantFormatTests = []wrongformatTest{
	{nil, "1.0+abc"},
	{nil, "1.0-rc.1"},
}

func TestWrongTolerantFormat(t *testing.T) {
	for _, test := range wrongTolerantFormatTests {
		if res, err := ParseTolerant(test.str); err == nil {
			t.Errorf("Parsing wrong format version %q, expected error but got %q", test.str, res)
		}
	}
}

func TestCompareHelper(t *testing.T) {
	v := &SemanticVersion{1, 0, 0, []PreReleaseValue{prstr("alpha")}, nil}
	v1 := &SemanticVersion{1, 0, 0, nil, nil}

	if !v.Equals(v) {
		t.Errorf("%q should be equal to %q", v, v)
	}
	if v1.Equals(v) {
		t.Errorf("%q should not be equal to %q", v1, v)
	}
	if v.Compare(v) < 0 {
		t.Errorf("%q should be greater than or equal to %q", v, v)
	}
	if v.Compare(v) > 0 {
		t.Errorf("%q should be less than or equal to %q", v, v)
	}
	if v.Compare(v1) >= 0 {
		t.Errorf("%q should be less than %q", v, v1)
	}
	if !(v.Compare(v1) <= 0) {
		t.Errorf("%q should be less than or equal %q", v, v1)
	}
	if !(v1.Compare(v) > 0) {
		t.Errorf("%q should be greater than %q", v1, v)
	}
	if !(v1.Compare(v) >= 0) {
		t.Errorf("%q should be greater than or equal %q", v1, v)
	}
}

func TestPreReleaseVersions(t *testing.T) {
	p1, err := NewPreReleaseValue("123")
	if !p1.IsNumeric() {
		t.Errorf("Expected numeric prversion, got %q", p1)
	}
	if p1.numValue != 123 {
		t.Error("Wrong prversion number")
	}
	if err != nil {
		t.Errorf("Not expected error %q", err)
	}
	p2, err := NewPreReleaseValue("alpha")
	if p2.IsNumeric() {
		t.Errorf("Expected non-numeric prversion, got %q", p2)
	}
	if p2.strValue != "alpha" {
		t.Error("Wrong prversion string")
	}
	if err != nil {
		t.Errorf("Not expected error %q", err)
	}
}

func TestBuildMetaDataVersions(t *testing.T) {
	_, err := NewBuildValue("123")
	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}

	_, err = NewBuildValue("build")
	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}

	_, err = NewBuildValue("test?")
	if err == nil {
		t.Error("Expected error, got none")
	}

	_, err = NewBuildValue("")
	if err == nil {
		t.Error("Expected error, got none")
	}
}

func TestNewHelper(t *testing.T) {
	v, err := Parse("1.2.3")
	if err != nil {
		t.Fatalf("Unexpected error %q", err)
	}

	// New returns pointer
	if v == nil {
		t.Fatal("Version is nil")
	}
	if v.Compare(&SemanticVersion{1, 2, 3, nil, nil}) != 0 {
		t.Fatal("Unexpected comparison problem")
	}
}

func TestMakeHelper(t *testing.T) {
	v, err := Parse("1.2.3")
	if err != nil {
		t.Fatalf("Unexpected error %q", err)
	}
	if v.Compare(&SemanticVersion{1, 2, 3, nil, nil}) != 0 {
		t.Fatal("Unexpected comparison problem")
	}
}

func BenchmarkParseSimple(b *testing.B) {
	const VERSION = "0.0.1"
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(VERSION)
	}
}

func BenchmarkParseComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(VERSION)
	}
}

func BenchmarkParseAverage(b *testing.B) {
	l := len(formatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(formatTests[n%l].result)
	}
}

func BenchmarkParseTolerantAverage(b *testing.B) {
	l := len(tolerantFormatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ParseTolerant(tolerantFormatTests[n%l].result)
	}
}

func BenchmarkStringSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkStringLarger(b *testing.B) {
	const VERSION = "11.15.2012"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkStringComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkStringAverage(b *testing.B) {
	l := len(formatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = formatTests[n%l].v.String()
	}
}

func BenchmarkValidateSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Validate()
	}
}

func BenchmarkValidateComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Validate()
	}
}

func BenchmarkValidateAverage(b *testing.B) {
	l := len(formatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		formatTests[n%l].v.Validate()
	}
}

func BenchmarkCompareSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Compare(v)
	}
}

func BenchmarkCompareComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Compare(v)
	}
}

func BenchmarkCompareAverage(b *testing.B) {
	l := len(compareTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		compareTests[n%l].v1.Compare((compareTests[n%l].v2))
	}
}
