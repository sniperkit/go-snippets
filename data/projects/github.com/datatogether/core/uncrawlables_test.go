package core

import (
	"github.com/datatogether/sql_datastore"
	"testing"
)

func TestListUncrawlables(t *testing.T) {
	resetTestData(appDB, "uncrawlables")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Uncrawlable{}); err != nil {
		t.Error(err.Error())
		return
	}

	uncrawlables, err := ListUncrawlables(store, 20, 0)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(uncrawlables) != 1 {
		t.Errorf("uncrawlables length mismatch")
	}

	uncrawlables, err = ListUncrawlables(store, 20, 10)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(uncrawlables) != 0 {
		t.Errorf("uncrawlables length mismatch")
	}
}
