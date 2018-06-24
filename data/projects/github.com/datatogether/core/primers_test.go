package core

import (
	"github.com/datatogether/sql_datastore"
	"testing"
)

func TestCountPrimers(t *testing.T) {
	// TODO - some test isn't cleaning up after itself.
	if err := resetTestData(appDB, "primers"); err != nil {
		t.Fatal(err)
	}

	c, err := CountPrimers(appDB)
	if err != nil {
		t.Error(err.Error())
	}

	if c != 6 {
		t.Errorf("wrong number of primers, expected %d, got: %d", 6, c)
	}
}

func TestListPrimers(t *testing.T) {
	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Primer{}); err != nil {
		t.Error(err.Error())
		return
	}

	primers, err := ListPrimers(store, 20, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(primers) != 6 {
		t.Errorf("primers length mismatch")
	}

	primers, err = ListPrimers(store, 20, 10)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(primers) != 0 {
		t.Errorf("primers length mismatch")
	}
}
