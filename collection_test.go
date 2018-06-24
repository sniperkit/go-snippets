package core

import (
	"fmt"
	"github.com/datatogether/sql_datastore"
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestCollectionStorage(t *testing.T) {
	store := datastore.NewMapDatastore()

	c := &Collection{Title: "test collection"}
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c.Creator = "penelope"
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c2 := &Collection{Id: c.Id}
	if err := c2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if err := CompareCollections(c, c2); err != nil {
		t.Error(err.Error())
		return
	}

	if err := c.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestCollectionSQLStorage(t *testing.T) {
	// confirm collection conforms to sql datastore mode
	_ = sql_datastore.Model(&Collection{})

	defer resetTestData(appDB, "collections")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Collection{}); err != nil {
		t.Error(err.Error())
		return
	}

	c := &Collection{Title: "test collection", Url: "http://test.test"}
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c.Creator = "penelope"
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c2 := &Collection{Id: c.Id}
	if err := c2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if err := CompareCollections(c, c2); err != nil {
		t.Error(err.Error())
		return
	}

	c3 := &Collection{Url: c.Url}
	if err := c3.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if err := CompareCollections(c, c3); err != nil {
		t.Error(err.Error())
		return
	}

	if err := c.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}

func CompareCollections(a, b *Collection) error {
	if a.Id != b.Id {
		return fmt.Errorf("id mismatch: %s != %s", a.Id, b.Id)
	}

	if !a.Created.Equal(b.Created) {
		return fmt.Errorf("created doesn't match: %s != %s", a.Created.String(), b.Created.String())
	}

	if !a.Updated.Equal(b.Updated) {
		return fmt.Errorf("updated doesn't match: %s != %s", a.Updated.String(), b.Updated.String())
	}

	return nil
}
