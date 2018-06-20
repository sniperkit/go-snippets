package yaml

import (
	"fmt"
	"testing"

	"github.com/mickep76/encdec"
)

func inList(a string, list []string) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}
	return false
}

func TestEncodings(t *testing.T) {
	encs := encdec.Encodings()

	if !inList("yaml", encs) {
		t.Error(fmt.Errorf("encoding not in list"))
	}

	if ok := encdec.Registered("yaml"); !ok {
		t.Error("not registered")
	}
}
