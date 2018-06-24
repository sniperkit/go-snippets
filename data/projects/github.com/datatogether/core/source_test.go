package core

import (
	"github.com/datatogether/sql_datastore"
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestSourceStorage(t *testing.T) {

	store := datastore.NewMapDatastore()

	c := &Source{Url: "youtube.com", Primer: &Primer{Id: "5b1031f4-38a8-40b3-be91-c324bf686a87"}, Crawl: true}
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c.Crawl = false
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c2 := &Source{Id: c.Id}
	if err := c2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if c2.Crawl != c.Crawl {
		t.Errorf("crawl doesn't match: %t != %t", c2.Crawl, c.Crawl)
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

func TestSourceSQLStorage(t *testing.T) {
	defer resetTestData(appDB, "crawl_urls", "sources", "primers")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Source{}); err != nil {
		t.Error(err)
		return
	}

	c := &Source{Url: "youtube.com", Primer: &Primer{Id: "5b1031f4-38a8-40b3-be91-c324bf686a87"}, Crawl: true}
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c.Crawl = false
	if err := c.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	c2 := &Source{Url: "youtube.com"}
	if err := c2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if c2.Crawl != c.Crawl {
		t.Errorf("crawl doesn't match: %t != %t", c2.Crawl, c.Crawl)
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

func TestSourceUndescribedContent(t *testing.T) {
	// TODO - some test isn't cleaning up after itself
	if err := resetTestData(appDB, "sources"); err != nil {
		t.Fatal(err.Error())
	}

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Source{}); err != nil {
		t.Error(err)
		return
	}

	c := &Source{Url: "www.census.gov"}
	if err := c.Read(store); err != nil {
		t.Error(err.Error())
		return
	}
	urls, err := c.UndescribedContent(appDB, 100, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(urls) != 1 {
		t.Errorf("UnescribedContent Fail:")
		for _, u := range urls {
			t.Error(u)
		}
	}
}

func TestSourceDescribedContent(t *testing.T) {
	// TODO - some test isn't cleaning up after itself
	if err := resetTestData(appDB, "sources"); err != nil {
		t.Fatal(err.Error())
	}

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Source{}); err != nil {
		t.Error(err)
		return
	}

	c := &Source{Url: "www.census.gov"}
	if err := c.Read(store); err != nil {
		t.Error(err.Error())
		return
	}
	urls, err := c.DescribedContent(appDB, 100, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(urls) != 1 {
		t.Errorf("DescribedContent Fail:")
		t.Error(urls)
	}
}
