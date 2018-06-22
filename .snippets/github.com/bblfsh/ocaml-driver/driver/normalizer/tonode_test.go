package normalizer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const fixtureDir = "fixtures"

func getFixture(name string) (interface{}, error) {
	path := filepath.Join(fixtureDir, name)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(f)
	var data interface{}
	if err := d.Decode(&data); err != nil {
		_ = f.Close()
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	return data, nil
}

func TestNativeToNode(t *testing.T) {
	f, err := getFixture("simple.ml.json")
	if err != nil {
		t.Fatal("open fixture")
	}
	n, err := ToNode.ToNode(f)
	if err != nil {
		t.Errorf("ToNode err: %s", err)
	}
	if n == nil {
		t.Error("n == nil")
	}
}
