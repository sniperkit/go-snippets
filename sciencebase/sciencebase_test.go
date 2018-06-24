package sciencebase

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestItemJson(t *testing.T) {
	cases := []struct {
		path string
		err  error
	}{
		{"testdata/55b943ade4b09a3b01b65d78.json", nil},
		{"testdata/4f4e476ae4b07f02db47e13b.json", nil},
	}

	for i, c := range cases {
		data, err := ioutil.ReadFile(c.path)
		if err != nil {
			t.Errorf("case %d error opening file: %s", i, err.Error())
			continue
		}

		item := &Item{}
		if err := json.Unmarshal(data, item); err != nil {
			t.Errorf("case %d unmarshal error: %s", i, err.Error())
			continue
		}

		data, err = json.MarshalIndent(item, "", "  ")
		if err != nil {
			t.Errorf("case %d marshal error: %s", i, err.Error())
			continue
		}

		if err := ioutil.WriteFile(c.path[:len(c.path)-len(".json")]+"_parsed.json", data, 0777); err != nil {
			t.Errorf("case %d error writing file: %s", err.Error())
		}
	}
}
