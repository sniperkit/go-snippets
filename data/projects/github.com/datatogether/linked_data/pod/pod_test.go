package pod

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestCatalogJson(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/epa.json")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	cat := &Catalog{}

	if err := json.Unmarshal(data, cat); err != nil {
		t.Errorf("unmarshal error: %s", err.Error())
		return
	}

	data, err = json.MarshalIndent(cat, "", "  ")
	if err != nil {
		t.Errorf("marshal error: %s", err.Error())
		return
	}

	// if err := ioutil.WriteFile("testdata/epa_parsed.json", data, 0777); err != nil {
	// 	t.Errorf("error writing file: %s", err.Error())
	// }
}
