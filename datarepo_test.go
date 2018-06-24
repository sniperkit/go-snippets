package core

import (
	"github.com/datatogether/sql_datastore"
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestDataRepoStorage(t *testing.T) {
	store := datastore.NewMapDatastore()

	c := &DataRepo{Title: "test dataRepo"}
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c.Description = "test description!"
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c2 := &DataRepo{Id: c.Id}
	if err := c2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if !c2.Created.Equal(c.Created) {
		t.Errorf("created doesn't match: %s != %s", c2.Created.String(), c.Created.String())
	}

	if !c2.Updated.Equal(c.Updated) {
		t.Errorf("updated doesn't match: %s != %s", c2.Updated.String(), c.Updated.String())
	}

	if err := c.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestDataRepoSQLStorage(t *testing.T) {
	defer resetTestData(appDB, "dataRepos")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&DataRepo{}); err != nil {
		t.Error(err.Error())
		return
	}

	c := &DataRepo{Title: "test dataRepo"}
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c.Description = "test description!"
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c2 := &DataRepo{Id: c.Id}
	if err := c2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if !c2.Created.Equal(c.Created) {
		t.Errorf("created doesn't match: %s != %s", c2.Created.String(), c.Created.String())
	}

	if !c2.Updated.Equal(c.Updated) {
		t.Errorf("updated doesn't match: %s != %s", c2.Updated.String(), c.Updated.String())
	}

	if err := c.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}
