package cove

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCoverageProfileDoesntDeleteNonProfileFilesInOutputDir(t *testing.T) {
	testdir, _ := ioutil.TempDir("", "coverage-profile")
	defer os.RemoveAll(testdir)

	bytes := []byte("foo")
	foofile := filepath.Join(testdir, "foo")
	ioutil.WriteFile(foofile, bytes, 0644)

	profile, _ := CoverageProfile(true, "count", testdir, "text/scanner")

	if profile == "" {
		t.Errorf("Could not create profile")
	}

	if _, err := os.Stat(foofile); os.IsNotExist(err) {
		t.Errorf("Foo file was deleted")
	}
}

func TestCoverageProfileCreatesOutputDirIfItDoesntExist(t *testing.T) {
	testdir, _ := ioutil.TempDir("", "coverage-profile")
	defer os.RemoveAll(testdir)
	outdir := filepath.Join(testdir, "outdir")

	profile, er := CoverageProfile(true, "set", outdir, "text/scanner")
	if er != nil {
		t.Errorf("%v", er)
	}

	if filepath.Dir(profile) != outdir {
		t.Errorf("Could not create profile %v in %v", profile, outdir)
	}
}
