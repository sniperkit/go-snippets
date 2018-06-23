package commands

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/ngageoint/seed-common/util"
)

func init() {
	util.InitPrinter(util.PrintErr)
}

func TestSeedInit(t *testing.T) {
	cases := []struct {
		directory   string
		version     string
		expectedErr error
	}{
		{"../testdata/dummy-scratch/", "0.0.0", errors.New("This version of seed-cli does not have a sample manifest for version 0.0.0")},
		{"../testdata/dummy-scratch/", "1.0.0", nil},
		{"../testdata/complete/", "1.0.0",errors.New("Existing file left unmodified.")},
	}

	for _, c := range cases {
		err := SeedInit(c.directory, c.version)

		if c.expectedErr == nil && err != nil {
			t.Errorf("SeedInit(%q) == %v, expected %v", c.directory, err.Error(), c.expectedErr)
		}
		if err != nil {
			if !strings.Contains(err.Error(), c.expectedErr.Error()) {
				t.Errorf("SeedInit(%q) == %v, expected %v", c.directory, err.Error(), c.expectedErr.Error())
			}
		}
	}

	// Cleanup test file
	os.Remove("../testdata/dummy-scratch/seed.manifest.json")
}
