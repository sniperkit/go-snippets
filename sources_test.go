package core

import (
	"github.com/datatogether/sql_datastore"
	"testing"
)

func TestCountSources(t *testing.T) {
	// TODO - some test isn't cleaning up after itself.
	if err := resetTestData(appDB, "sources"); err != nil {
		t.Fatal(err)
	}

	c, err := CountSources(appDB)
	if err != nil {
		t.Error(err.Error())
	}

	if c != 4 {
		t.Errorf("wrong number of sources, expected %d, got: %d", 4, c)
	}
}

func TestListSources(t *testing.T) {
	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Source{}); err != nil {
		t.Error(err.Error())
		return
	}

	sources, err := ListSources(store, 20, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(sources) != 4 {
		t.Errorf("sources length mismatch")
	}

	sources, err = ListSources(store, 20, 10)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(sources) != 0 {
		t.Errorf("sources length mismatch")
	}
}

func TestCrawlingSources(t *testing.T) {
	s, err := CrawlingSources(appDB, 100, 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(s) != 3 {
		t.Errorf("wrong number of crawling sources, expected %d, got: %d", 3, len(s))
	}
}
