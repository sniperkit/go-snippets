package commands

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ngageoint/seed-common/objects"
	"github.com/ngageoint/seed-common/util"
)

func init() {
	util.InitPrinter(util.PrintErr)
}

func TestProcessDirectory(t *testing.T) {
	cases := []struct {
		batchDir         string
		outDir           string
		manifestFile     string
		expected         string
		expectedErrorMsg string
	}{
		{"../testdata", "../testdata/test-extract", "../examples/extractor/seed.manifest.json",
			"[{[ZIP=../testdata/batch-test.csv] [] ../testdata/test-extract/batch-test.csv} " +
				"{[ZIP=../testdata/empty-batch.csv] [] ../testdata/test-extract/empty-batch.csv} " +
				"{[ZIP=../testdata/missing-keys.csv] [] ../testdata/test-extract/missing-keys.csv} " +
				"{[ZIP=../testdata/seed-scale.zip] [] ../testdata/test-extract/seed-scale.zip}]",
			""},
		{"../testdata", "../testdata/test-multiple", "../testdata/multiple-required-inputs/seed.manifest.json",
			"[]", "ERROR: Multiple required inputs are not supported when batch processing directories."},
		{"../testdata", "../testdata/test-no-inputs", "../testdata/no-inputs/seed.manifest.json",
			"[]", "ERROR: Could not determine which input to use from Seed manifest."},
	}

	for _, c := range cases {
		os.Mkdir(c.outDir, os.ModePerm)
		defer os.Remove(c.outDir)
		seed := objects.SeedFromManifestFile(c.manifestFile)
		out, err := ProcessDirectory(seed, c.batchDir, c.outDir)
		outstr := fmt.Sprintf("%v", out)
		if outstr != c.expected {
			t.Errorf("ProcessDirectory(%q, %q, %q) == %v, expected %v", c.manifestFile, c.batchDir, c.outDir, outstr, c.expected)
		}
		if err != nil {
			if !strings.Contains(err.Error(), c.expectedErrorMsg) {
				t.Errorf("ProcessDirectory(%q, %q, %q) == %v, expected %v", c.manifestFile, c.batchDir, c.outDir, err.Error(), c.expectedErrorMsg)
			}
		}
	}
}

func TestProcessBatchFile(t *testing.T) {
	cases := []struct {
		batchFile        string
		outDir           string
		manifestFile     string
		expected         string
		expectedErrorMsg string
	}{
		{"../testdata/batch-test.csv", "../testdata/test-extract-file", "../examples/extractor/seed.manifest.json",
			"[{[ZIP=/home/jtobe/go/src/github.com/ngageoint/seed-cli/testdata/test1.zip] [] ../testdata/test-extract-file/1-test1.zip} " +
				"{[ZIP=/home/jtobe/go/src/github.com/ngageoint/seed-cli/testdata/test2.zip] [] ../testdata/test-extract-file/2-test2.zip} " +
				"{[ZIP=/home/jtobe/go/src/github.com/ngageoint/seed-cli/testdata/test3.zip] [] ../testdata/test-extract-file/3-test3.zip}]",
			""},
		{"../testdata/empty-batch.csv", "../testdata/test-empty", "../testdata/multiple-required-inputs/seed.manifest.json",
			"[]", "ERROR: Empty batch file"},
		{"../testdata/missing-keys.csv", "../testdata/missing-keys", "../testdata/no-inputs/seed.manifest.json",
			"[]", "ERROR: Empty keys list on first line of batch file."},
	}

	for _, c := range cases {
		os.Mkdir(c.outDir, os.ModePerm)
		defer os.Remove(c.outDir)
		seed := objects.SeedFromManifestFile(c.manifestFile)
		out, err := ProcessBatchFile(seed, c.batchFile, c.outDir)
		outstr := fmt.Sprintf("%v", out)
		fmt.Println(outstr)
		fmt.Println(c.expected)
		if outstr != c.expected {
			t.Errorf("ProcessFile(%q, %q, %q) == %v, expected %v", c.manifestFile, c.batchFile, c.outDir, outstr, c.expected)
		}
		if err != nil {
			if !strings.Contains(err.Error(), c.expectedErrorMsg) {
				t.Errorf("ProcessFile(%q, %q, %q) == %v, expected %v", c.manifestFile, c.batchFile, c.outDir, err.Error(), c.expectedErrorMsg)
			}
		}
	}
}
