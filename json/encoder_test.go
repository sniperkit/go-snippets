package json

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mickep76/encdec"
)

var (
	testJSON = `{"test":{"abc":{"myTest":"Test2"},"myTest":"Test"}}
`

	testJSONIndent = `{
  "test": {
    "abc": {
      "myTest": "Test2"
    },
    "myTest": "Test"
  }
}
`

	testMap = map[string]interface{}{
		"test": map[string]interface{}{
			"myTest": "Test",
			"abc": map[string]interface{}{
				"myTest": "Test2",
			},
		},
	}

	tmpFile = "test.tmp"
)

func readFile(fn string) ([]byte, error) {
	fp, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	if err := fp.Close(); err != nil {
		return nil, err
	}

	return b, nil
}

func TestToByte(t *testing.T) {
	b, err := encdec.ToBytes("json", testMap)
	if err != nil {
		t.Error(err)
	}

	if string(b) != testJSON {
		t.Errorf("want:\n%s, got:\n%s", testJSON, string(b))
	}
}

func TestToByteWithIndent(t *testing.T) {
	b, err := encdec.ToBytes("json", testMap, encdec.WithIndent("  "))
	if err != nil {
		t.Error(err)
	}

	if string(b) != testJSONIndent {
		t.Errorf("want:\n%s, got:\n%s", testJSONIndent, string(b))
	}
}

func TestToFile(t *testing.T) {
	if err := encdec.ToFile("json", tmpFile, testMap); err != nil {
		t.Error(err)
	}

	b, err := readFile(tmpFile)
	if err != nil {
		t.Error(err)
	}

	if string(b) != testJSON {
		t.Errorf("want:\n%s, got:\n%s", testJSONIndent, string(b))
	}

	if err := os.Remove(tmpFile); err != nil {
		t.Error(err)
	}
}
