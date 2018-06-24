package core

import (
	"github.com/datatogether/sql_datastore"
	"testing"
)

func TestListCollections(t *testing.T) {
	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Collection{}); err != nil {
		t.Error(err.Error())
		return
	}

	collections, err := ListCollections(store, 20, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(collections) != 4 {
		t.Errorf("collections length mismatch")
	}

	collections, err = ListCollections(store, 20, 1)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(collections) != 3 {
		t.Errorf("collections length mismatch")
	}
}

func TestCollectionsByCreator(t *testing.T) {
	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Collection{}); err != nil {
		t.Error(err.Error())
		return
	}

	collections, err := CollectionsByCreator(store, "EDGI_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc", "created DESC", 20, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(collections) != 2 {
		t.Errorf("collections length mismatch: %d != %d", len(collections), 2)
	}
}
