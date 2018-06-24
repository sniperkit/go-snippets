package core

import (
	"github.com/datatogether/sql_datastore"
	"testing"
	"time"
)

func TestSnapshotStorge(t *testing.T) {
	defer resetTestData(appDB, "snapshots")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Url{}); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()
	u := &Url{
		Url:          "http://www.epa.gov",
		LastGet:      &now,
		Status:       200,
		DownloadTook: 20,
		Headers:      []string{"test", "header"},
		Hash:         "thisshouldbeahash",
	}

	if err := WriteSnapshot(store, u); err != nil {
		t.Error(err.Error())
		return
	}
}
