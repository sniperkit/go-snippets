package docindex

import (
	"bytes"
	"reflect"
	"testing"
)

func TestIndex_Open(t *testing.T) {
	type Metadata struct {
		Title string
	}

	idx, err := Open("testdata", Metadata{})
	if err != nil {
		t.Fatal("Open:", err)
	}

	wantFilenames := []string{"testdata/foo.txt", "testdata/subdir/bar.txt"}
	gotFilenames := idx.Filenames()
	if !reflect.DeepEqual(wantFilenames, gotFilenames) {
		t.Errorf("want filenames %v, got %v", wantFilenames, gotFilenames)
	}

	metaTests := []struct {
		filename string
		wantMeta Metadata
		wantData []byte
	}{
		{"testdata/foo.txt", Metadata{"foo"}, []byte("Hello from foo.\n")},
		{"testdata/subdir/bar.txt", Metadata{"bar"}, []byte("This is bar.\n")},
	}
	for _, test := range metaTests {
		var gotMeta Metadata
		gotData, err := idx.Doc(test.filename, &gotMeta)
		if err != nil {
			t.Errorf("Doc %s: %s", test.filename, err)
			continue
		}
		if !reflect.DeepEqual(test.wantMeta, gotMeta) {
			t.Errorf("Doc %s: want meta %+v, got %+v", test.filename, test.wantMeta, gotMeta)
		}
		if !bytes.Equal(test.wantData, gotData) {
			t.Errorf("Doc %s: want data %q, got %q", test.filename, test.wantData, gotData)
		}
	}

	wantAllMetadata := map[string]*Metadata{"testdata/foo.txt": {"foo"}, "testdata/subdir/bar.txt": {"bar"}}
	gotAllMetadata := make(map[string]*Metadata)
	idx.AllMetadata(gotAllMetadata)
	if !reflect.DeepEqual(wantAllMetadata, gotAllMetadata) {
		t.Errorf("all metadata: want %+v, got %+v", wantAllMetadata, gotAllMetadata)
	}
}
