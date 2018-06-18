package fileutil

import (
	"os"
	"testing"
)

// TestFind verifies that Find ... doesn't crash.  This needs a mock for the
// Readdir wrapper function so we can get actual high-level testing without
// relying on a completely unknown filesystem....
func TestFind(t *testing.T) {
	var _, err = Find(os.TempDir(), 1)
	if err != nil {
		t.Fatalf("Got an error trying to read the filesystem!  %s", err)
	}
}

func TestReaddir(t *testing.T) {
	var infos, err = ReaddirSorted(os.TempDir())
	if err != nil {
		t.Fatalf("Got an error trying to read the filesystem!  %s", err)
	}

	t.Logf("Found %d Files:", len(infos))
	for _, info := range infos {
		t.Log("  - " + info.Name())
	}
}
