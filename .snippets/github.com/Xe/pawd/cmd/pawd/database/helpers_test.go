package database

import (
	"context"
	"os"
	"testing"

	"github.com/Xe/uuid"
	"github.com/asdine/storm"
)

type testingDB struct {
	*storm.DB
	fname string
}

func (tdb *testingDB) Close() {
	tdb.DB.Close()
	os.Remove(tdb.fname)
}

func newTestingDB(t *testing.T) (ctx context.Context, tdb *testingDB) {
	fname := uuid.New() + ".db"

	db, err := storm.Open(fname)
	if err != nil {
		t.Fatal(err)
	}

	return context.Background(), &testingDB{DB: db, fname: fname}
}
