package yaml

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mickep76/encdec"
)

var (
	testYAML = `test:
  abc:
    myTest: Test2
  myInt: 1
  myList:
  - Item 1
  - Item 2
  myTest: Test
`

	testMap = map[string]interface{}{
		"test": map[string]interface{}{
			"myTest": "Test",
			"abc": map[string]interface{}{
				"myTest": "Test2",
			},
			"myList": []interface{}{
				"Item 1",
				"Item 2",
			},
			"myInt": 1,
		},
	}

	testMapInterface = map[interface{}]interface{}{
		"test": map[interface{}]interface{}{
			"myTest": "Test",
			"abc": map[interface{}]interface{}{
				"myTest": "Test2",
			},
			"myList": []interface{}{
				"Item 1",
				"Item 2",
			},
			"myInt": 1,
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
	b, err := encdec.ToBytes("yaml", testMap)
	if err != nil {
		t.Error(err)
	}

	if string(b) != testYAML {
		t.Errorf("want:\n%s, got:\n%s", testYAML, string(b))
	}
}

func TestToByteWithIndent(t *testing.T) {
	if _, err := encdec.ToBytes("yaml", testMap, encdec.WithIndent("")); err != nil {
		t.Error(err)
	}
}

func TestToFile(t *testing.T) {
	if err := encdec.ToFile("yaml", tmpFile, testMap); err != nil {
		t.Error(err)
	}

	b, err := readFile(tmpFile)
	if err != nil {
		t.Error(err)
	}

	if string(b) != testYAML {
		t.Errorf("want:\n%s, got:\n%s", testYAML, string(b))
	}

	if err := os.Remove(tmpFile); err != nil {
		t.Error(err)
	}
}
