package main

import (
	"testing"
)

var outputFilenameTests = []struct {
	in  string
	out string
}{
	{"test.jpg", "test_clean.jpg"},
	{"test.png", "test_clean.png"},
	{"thishasnoextension", "thishasnoextension_clean"},
}

func TestGetDefaultOutputFilename(t *testing.T) {
	for _, tt := range outputFilenameTests {
		s := getDefaultOutputFilename(tt.in)
		if s != tt.out {
			t.Errorf("getDefaultOutputFilename(%q) => %q, want %q", tt.in, s, tt.out)
		}
	}
}

var doCleaningTests = []struct {
	inputImage    string
	outputImage   string
	jsonFlag      bool
	cleanFlag     bool
	errorExpected bool
}{
	// Test program completes if input file is jpeg
	{"test/test.jpeg", "", false, true, false},
	// Test program completes if input file is png
	{"test/test.png", "", false, true, false},
	// Test program bails if no input file provided
	{"", "", false, true, true},
	// Test program bails if input file provided but it does not exist
	{"test/doesnotexist.jpeg", "", false, true, true},
	// Test program does not bail when JSON output is to be provided
	{"test/test.jpeg", "", true, true, false},
	// Test program does not bail when JSON output is to be provided for file without exif metadata
	{"test/test.png", "", true, true, false},
}

func TestDoCleaningSuccess(t *testing.T) {
	for _, tt := range doCleaningTests {
		s := doCleaning(tt.inputImage, tt.cleanFlag, tt.outputImage, tt.jsonFlag)
		if s != nil && tt.errorExpected == false {
			t.Errorf("doCleaning(%q, %q, %q, %q) => %q, expected errors: %q", tt.inputImage, tt.cleanFlag, tt.outputImage, tt.jsonFlag, s, tt.errorExpected)
		}
	}
}
