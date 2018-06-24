package core

import "testing"

func TestDatasetRead(t *testing.T) {
	a := NewDataset(trainSet)

	err := a.ReadFromFile()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	attrs := len(a.header)
	for _, tu := range a.data {
		if len(tu.Data) != attrs {
			t.Log(tu.String() + " not having correct attrs")
			t.FailNow()
		}
	}

}
