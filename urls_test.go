package core

import (
	"github.com/datatogether/sql_datastore"
	"testing"
)

func TestListUrls(t *testing.T) {
	resetTestData(appDB, "urls")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Url{}); err != nil {
		t.Error(err.Error())
		return
	}

	urls, err := ListUrls(store, 20, 0)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(urls) != 9 {
		t.Errorf("urls length mismatch")
	}

	urls, err = ListUrls(store, 20, 10)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(urls) != 0 {
		t.Errorf("urls length mismatch")
	}
}
