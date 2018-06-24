package core

import (
	"github.com/datatogether/sql_datastore"
	"testing"
)

func TestCollectionItemsSQLStorage(t *testing.T) {
	defer resetTestData(appDB, "collection_items", "collections", "urls")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Collection{}, &CollectionItem{}, &Url{}); err != nil {
		t.Error(err.Error())
		return
	}

	c := &Collection{Title: "test collection Item Count"}
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	count, err := c.ItemCount(store)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if count != 0 {
		t.Errorf("count mistmatch. expected %d, got %d", 0, count)
		return
	}

	saveitems := []*CollectionItem{
		&CollectionItem{Url: Url{Url: "http://test.0.com"}, Index: 0, Description: "zero"},
		&CollectionItem{Url: Url{Url: "http://test.1.com"}, Index: 1, Description: "one"},
		&CollectionItem{Url: Url{Url: "http://test.2.com"}, Index: 2, Description: "two"},
	}

	if err := c.SaveItems(store, saveitems); err != nil {
		t.Error(err.Error())
		return
	}

	if count, err := c.ItemCount(store); err != nil {
		t.Error(err.Error())
		return
	} else if count != len(saveitems) {
		t.Errorf("count mistmatch. expected %d, got %d", len(saveitems), count)
		return
	}

	readitems, err := c.ReadItems(store, "created DESC", 100, 0)
	if err != nil {
		t.Error(err.Error())
		return
	} else if len(readitems) != len(saveitems) {
		t.Errorf("count mistmatch. expected %d, got %d", len(saveitems), count)
		return
	}

	if err := c.DeleteItems(store, saveitems); err != nil {
		t.Error(err.Error())
		return
	}

	if count, err := c.ItemCount(store); err != nil {
		t.Error(err.Error())
		return
	} else if count != 0 {
		t.Errorf("count mistmatch. expected %d, got %d", 0, count)
		return
	}

}
