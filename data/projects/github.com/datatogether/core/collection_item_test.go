package core

import (
	"fmt"
	"github.com/datatogether/sql_datastore"
	"github.com/ipfs/go-datastore"
	"testing"
)

// confirm CollectionItem conforms to sql datastore mode
var _ = sql_datastore.Model(&CollectionItem{})

func TestCollectionItemStorage(t *testing.T) {
	// store := datastore.NewLogDatastore(datastore.NewMapDatastore(), "CollectionItemTest")
	store := datastore.NewMapDatastore()

	c := &Collection{Title: "test collection"}
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	u := Url{Url: "http://example.url.test.com"}

	item := &CollectionItem{collectionId: c.Id, Url: u, Index: 0, Description: "Description"}
	if err := item.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	item2 := &CollectionItem{collectionId: c.Id, Url: item.Url}
	if err := item2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if err := CompareCollectionItems(item, item2); err != nil {
		t.Error(err.Error())
		return
	}

	item.Description = "Updated Description"
	if err := item.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	item2 = &CollectionItem{collectionId: c.Id, Url: item.Url}
	if err := item2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if err := CompareCollectionItems(item, item2); err != nil {
		t.Error(err.Error())
		return
	}

	if err := item.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestCollectionItemSQLStorage(t *testing.T) {
	// defer resetTestData(appDB, "collections")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Collection{}, &CollectionItem{}, &Url{}); err != nil {
		t.Error(err.Error())
		return
	}

	c := &Collection{Title: "test collection"}
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	item := &CollectionItem{collectionId: c.Id, Url: Url{Url: "http://test.url.two.example.test"}, Index: 0, Description: "item description"}
	if err := item.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	item.Description = "updated item description"
	if err := item.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	item2 := &CollectionItem{collectionId: c.Id, Url: Url{Id: item.Id}}
	if err := item2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if err := CompareCollectionItems(item, item2); err != nil {
		t.Error(err.Error())
		return
	}

	urlItem := &CollectionItem{collectionId: c.Id, Url: Url{Url: item.Url.Url}}
	if err := urlItem.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if err := CompareCollectionItems(item, urlItem); err != nil {
		t.Error(err.Error())
		return
	}

	if err := item.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}

func CompareCollectionItems(a, b *CollectionItem) error {
	// TODO - can't compare full urls b/c we don't always fill the whole thing out
	// if err := CompareUrls(&a.Url, &b.Url); err != nil {
	// 	return fmt.Errorf("url mismatch: %s", err.Error())
	// }

	if a.Url.Id != b.Url.Id {
		return fmt.Errorf("url id mismatch %s != %s", a.Url.Id, b.Url.Id)
	}

	if a.Url.Url != b.Url.Url {
		return fmt.Errorf("url mismatch %s != %s", a.Url.Url, b.Url.Url)
	}

	if a.collectionId != b.collectionId {
		return fmt.Errorf("collectionId mismatch: %d != %d", a.Index, b.Index)
	}
	if a.Index != b.Index {
		return fmt.Errorf("Index mismatch: %d != %d", a.Index, b.Index)
	}
	if a.Description != b.Description {
		return fmt.Errorf("Description mistmatch: %s != %s", a.Description, b.Description)
	}
	return nil
}
