package cove

import (
	"os"
	"testing"
)

func TestToSingleAndMultiGoPath(t *testing.T) {
	oldPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldPath)

	os.Setenv("GOPATH", "/foo/bar:/local/foo/bar")
	first := GetFirstGoPath()
	if first != "/foo/bar" {
		t.Errorf("Didn't get the correct to path from to for a multi-path: %s", first)
	}
}
