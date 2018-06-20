package toml

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/mickep76/encdec"
)

func TestFromByte(t *testing.T) {
	g := Messages{}
	if err := encdec.FromBytes("toml", []byte(testTOML), &g); err != nil {
		t.Error(err)
	}

	if err := deep.Equal(g, testStruct); err != nil {
		t.Error(err)
	}
}

func TestFromByteWithMapString(t *testing.T) {
	g := Messages{}
	if err := encdec.FromBytes("toml", []byte(testTOML), &g, encdec.WithMapString()); err != nil {
		t.Error(err)
	}
}

func TestFromFile(t *testing.T) {
	g := Messages{}
	if err := encdec.FromFile("toml", "test.toml", &g); err != nil {
		t.Error(err)
	}

	if err := deep.Equal(g, testStruct); err != nil {
		t.Error(err)
	}
}
