package core

import (
	"github.com/datatogether/sql_datastore"
	"testing"
)

func TestListCustomCrawls(t *testing.T) {
	resetTestData(appDB, "custom_crawls")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&CustomCrawl{}); err != nil {
		t.Error(err.Error())
		return
	}

	custom_crawls, err := ListCustomCrawls(store, 20, 0)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(custom_crawls) != 1 {
		t.Errorf("custom_crawls length mismatch")
	}

	custom_crawls, err = ListCustomCrawls(store, 20, 10)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(custom_crawls) != 0 {
		t.Errorf("custom_crawls length mismatch")
	}
}
