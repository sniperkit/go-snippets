package pdf

import (
	"fmt"
	"testing"
)

var climateChangeMeta = map[string]interface{}{
	"title": "What Climate Change Means for Massachusetts",
}

func TestMetadataForFile(t *testing.T) {
	meta, err := MetadataForFile("test_files/climate-change-ma.pdf")
	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Println(meta)

	if err := CompareObjects(climateChangeMeta, meta); err != nil {
		t.Error(err.Error())
	}
}

func CompareObjects(a, b map[string]interface{}) error {
	if a == nil || b == nil && (a != nil && b != nil) {
		return fmt.Errorf("nil map mismatch: %s != %s", a, b)
	}
	if len(a) != len(b) {
		return fmt.Errorf("map length mismatch: %d != %d", len(a), len(b))
	}

	for key, val := range a {
		if b[key] == nil {
			return fmt.Errorf("key '%s' missing", key)
		}

		switch val.(type) {
		case string:
		case int:
		case float32:
		case bool:
		default:
		}

	}
	return nil
}
