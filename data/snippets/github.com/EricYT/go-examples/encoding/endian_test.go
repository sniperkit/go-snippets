package encoding

import "testing"

func TestEndian(t *testing.T) {
	goEndian := Endian()
	cEndian := CEndian()
	cEndianCon := CEndianConcision()
	if goEndian != cEndian {
		t.Fatalf("Endian result: %d CEndian result: %d not equal.", goEndian, cEndian)
	}

	if goEndian != cEndianCon {
		t.Fatalf("Endian result: %d CEndianConcision result: %d not equal.", goEndian, cEndianCon)
	}

}
