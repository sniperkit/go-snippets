package commands

import (
	"strings"
	"testing"

	common_const "github.com/ngageoint/seed-common/constants"
	"github.com/ngageoint/seed-common/util"
)

func init() {
	util.InitPrinter(util.PrintErr)
}

func TestValidate(t *testing.T) {
	cases := []struct {
		seedFileName     string
		expected         bool
		expectedErrorMsg string
	}{
		{"../examples/addition-job/seed.manifest.json", true, ""},
		{"../examples/extractor/seed.manifest.json", true, ""},
		{"../testdata/invalid-missing-job/seed.manifest.json",
			false, "job is required"},
		{"../testdata/invalid-missing-job-interface-inputs-files-name/seed.manifest.json",
			false, "name is required"},
		{"../testdata/invalid-reserved-name/seed.manifest.json",
			false, "Multiple Name values are assigned the same INPUT Name value. Each Name value must be unique."},
	}

	for _, c := range cases {
		name := util.GetFullPath(c.seedFileName, "")
		version := "1.0.0"
		err := ValidateSeedFile("", version, name, common_const.SchemaManifest)
		success := err == nil
		if success != c.expected {
			t.Errorf("ValidateSeedFile(%v, %v, %v, %v) == %v, expected %v", "", version, name, common_const.SchemaManifest, success, c.expected)
		}
		if err != nil {
			if !strings.Contains(err.Error(), c.expectedErrorMsg) {
				t.Errorf("ValidateSeedFile(%v, %v, %v, %v) == %v, expected %v", "", version, name, common_const.SchemaManifest, err.Error(), c.expectedErrorMsg)
			}
		}
	}
}
