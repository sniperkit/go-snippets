package yaml

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"

	"github.com/mickep76/encdec"
)

func TestFromByte(t *testing.T) {
	var g interface{}
	if err := encdec.FromBytes("yaml", []byte(testYAML), &g); err != nil {
		t.Error(err)
	}

	if err := deep.Equal(g, testMapInterface); err != nil {
		t.Error(err)
	}
}

func TestFromByteWithMapString(t *testing.T) {
	var g interface{}
	if err := encdec.FromBytes("yaml", []byte(testYAML), &g, encdec.WithMapString()); err != nil {
		t.Error(err)
	}

	if err := deep.Equal(g, testMap); err != nil {
		t.Error(err)
	}
}

func TestFromByteWithMapStringError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	g := ""
	if err := encdec.FromBytes("yaml", []byte(testYAML), g, encdec.WithMapString()); err != nil {
		t.Error(fmt.Errorf("expected a panic! %v", err))
	}
}

func TestFromFile(t *testing.T) {
	var g interface{}
	if err := encdec.FromFile("yaml", "test.yaml", &g); err != nil {
		t.Error(err)
	}

	if err := deep.Equal(g, testMapInterface); err != nil {
		t.Error(err)
	}
}
