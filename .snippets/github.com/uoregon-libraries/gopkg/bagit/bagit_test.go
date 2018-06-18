package bagit

import (
	"os"
	"sort"
	"testing"
)

func TestGenerateChecksums(t *testing.T) {
	var err error
	var path string
	var b *Bag

	path, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	b = New(path)
	err = b.GenerateChecksums()
	if err != nil {
		t.Fatalf("Should have generated checksums successfully in %q, but got error (%s)", b.root, err)
	}

	// Sort the list for easier comparison - filepath.Walk should already do
	// this, but I can't stand tests that stop passing just because of sorting
	sort.Slice(b.Checksums, func(i, j int) bool {
		return b.Checksums[i].Path < b.Checksums[j].Path
	})

	var expectedChecksums = []string{
		"60fa80b948a0acc557a6ba7523f4040a7b452736723df20f118d0aacb5c1901b", // another.txt's "sha256sum" value
		"55f8718109829bf506b09d8af615b9f107a266e19f7a311039d1035f180b22d4", // test.txt's "sha256sum" value
	}

	if len(b.Checksums) != len(expectedChecksums) {
		t.Errorf("Checksum list should have had %d items, but had %d", len(expectedChecksums), len(b.Checksums))
	}

	for i, ck := range b.Checksums {
		t.Logf("Checksum: %#v", ck)
		if ck.Checksum != expectedChecksums[i] {
			t.Errorf("Checksum for %q was %q, but should have been %q", ck.Path, ck.Checksum, expectedChecksums[i])
		}
	}
}
