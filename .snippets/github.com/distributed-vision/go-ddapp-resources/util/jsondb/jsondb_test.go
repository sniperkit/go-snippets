package jsondb

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/distributed-vision/go-resources/util"
)

func reset(file string) (string, error) {
	err := os.Remove(file)

	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	return file, nil
}

func TestFreelist(t *testing.T) {
	dbfile := filepath.Join(os.TempDir(), "test-file-fl.db")

	var blen uint = 1000

	file, err := reset(dbfile)

	if err != nil {
		t.Fatal("Reset failed:", err)
		return
	}

	var db = NewJsonDb(util.NewFileStorage(file, os.ModePerm), false)

	err = util.AwaitError(db.Open())

	if err != nil {
		t.Fatal("Open failed:", err)
		return
	}

	for nextBlockSize(blen) < uint(len(db.freelists)) {
		blen *= 2
	}

	var data = ""
	for i := uint(0); i < blen; i++ {
		data += "a"
	}

	var dbdata struct {
		D string
	}

	dbdata.D = data

	db.Set("a", dbdata)
	db.Set("b", dbdata)
	db.Set("a", dbdata)
	db.Delete("b")

	if db.Len() != 1 {
		t.Errorf("db.Len failed expected: 1 got: %d", db.Len())
	}

	// there was also a bug on re-loading freelists
	// for large block sizes, so let's read the DB back
	err = util.AwaitError(db.Sync())

	if err != nil {
		t.Fatal("Sync failed:", err)
		return
	}

	var db2 = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db2.Open())

	if err != nil {
		t.Fatal("db2.Open failed:", err)
		return
	}

	if db2.Len() != db.Len() {
		t.Errorf("db.Len failed expected: %d got: %d", db.Len(), db2.Len())
	}

	util.AwaitError(db.Close())

	util.AwaitError(db2.Clear())

	if db2.Len() != 0 {
		t.Errorf("db.Len failed expected: 0 got: %d", db2.Len())
	}

	util.AwaitError(db.Open())

	if db2.Len() != db.Len() {
		t.Errorf("db.Len failed expected: %d got: %d", db.Len(), db2.Len())
	}
}

func TestOpenWriteGet(t *testing.T) {
	dbfile := filepath.Join(os.TempDir(), "test-file-owg.db")

	file, err := reset(dbfile)

	if err != nil {
		t.Fatal("Reset failed:", err)
		return
	}

	//console.log("reset");
	var db = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db.Open())

	if err != nil {
		t.Fatal("Open failed:", err)
		return
	}
	//console.log("open");
	util.AwaitError(db.Set("hello", "world"))

	if value, ok := db.Get("hello"); !ok || value != "world" {
		t.Errorf("Get failed: expected %s got:%v", "world", value)
	}

	if !db.Has("hello") {
		t.Errorf("Has failed")
	}

	db.Set("hello-2", "world-2")

	if !db.Has("hello-2") {
		t.Errorf("Has 2 failed")
	}

	if db.Len() != 2 {
		t.Errorf("db.Len failed expected: 2 ot: %d", db.Len())
	}

	db.ForEach(func(key string, value interface{}) {
		if !db.Has(key) {
			t.Errorf("Undexpected key: %s in ForEach", key)
		}
		if val, ok := db.Get(key); !ok || val != value {
			t.Errorf("Undexpected value in ForEach: expected %s got:%v", "world", value)
		}
	})

	var db2 = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db2.Open())

	if err != nil {
		t.Fatal("Open failed:", err)
		return
	}

	if value, ok := db2.Get("hello-2"); !ok || value != "world-2" {
		t.Error("db2 Get failed")
	}

	if db2.Len() != db.Len() {
		t.Errorf("db.Len failed expected: %d got: %d", db.Len(), db2.Len())
	}
}

func TestDel(t *testing.T) {
	dbfile := filepath.Join(os.TempDir(), "test-file-del.db")

	file, err := reset(dbfile)

	if err != nil {
		t.Fatal("Reset failed:", err)
		return
	}

	var db = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db.Open())

	if err != nil {
		t.Fatal("Open failed:", err)
		return
	}

	db.Set("hello", "world")

	if value, ok := db.Get("hello"); !ok || value != "world" {
		t.Error("Get failed")
	}

	util.AwaitError(db.Delete("hello"))

	if _, ok := db.Get("hello"); ok {
		t.Error("Delete failed")
	}

	if db.Len() != 0 {
		t.Error("db.Len should be 0")
	}

	util.AwaitError(db.Sync())

	var db2 = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db2.Open())
	if err != nil {
		t.Fatal("db2 Open failed:", err)
		return
	}

	if _, ok := db.Get("hello"); ok {
		t.Error("db2 should not contain 'hello'")
	}

	if db2.Len() != db.Len() {
		t.Errorf("db2.Len failed expected: %d got: %d", db.Len(), db2.Len())
	}
}

func TestMultipleWrites(t *testing.T) {
	dbfile := filepath.Join(os.TempDir(), "test-file-mw.db")

	file, err := reset(dbfile)

	if err != nil {
		t.Fatal("Reset failed:", err)
		return
	}

	var db = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db.Open())

	if err != nil {
		t.Fatal("Open failed:", err)
		return
	}

	for i := 0; i < 20; i++ {
		db.Set(fmt.Sprintf("hello-%d", i), fmt.Sprintf("world-%d", i))
	}

	err = util.AwaitError(db.Sync())

	if err != nil {
		t.Fatal("Sync failed:", err)
		return
	}
	//time.Sleep(500 * time.Millisecond)
	var db2 = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db2.Open())
	if err != nil {
		t.Fatal("db2 Open failed:", err)
		return
	}

	for i := 0; i < 20; i++ {
		if value, ok := db2.Get(fmt.Sprintf("hello-%d", i)); !ok || value != fmt.Sprintf("world-%d", i) {
			t.Errorf("Get %d failed", i)
		}
	}

	db2.ForEach(func(key string, value interface{}) {
		if !db.Has(key) {
			t.Errorf("Undexpected key: %s in ForEach", key)
		}
		if val, ok := db.Get(key); !ok || val != value {
			t.Errorf("Undexpected value in ForEach: expected %s got:%v", "world", value)
		}
	})

}

func TestLastWriteWins(t *testing.T) {
	dbfile := filepath.Join(os.TempDir(), "test-file-lww.db")

	file, err := reset(dbfile)

	if err != nil {
		t.Fatal("Reset failed:", err)
		return
	}

	var db = NewJsonDb(util.NewFileStorage(file, os.ModePerm), false)

	err = util.AwaitError(db.Open())

	if err != nil {
		t.Fatal("Open failed:", err)
		return
	}

	for i := 0; i < 20; i++ {
		var count struct {
			Count int
		}

		count.Count = i
		db.Set("count", count)
	}

	err = util.AwaitError(db.Sync())

	if err != nil {
		t.Fatal("Sync failed:", err)
		return
	}

	var db2 = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db2.Open())

	if err != nil {
		t.Fatal("db2 Open failed:", err)
		return
	}

	if value, ok := db2.Get("count"); !ok || value.(map[string]interface{})["Count"] != float64(19) {
		t.Error("count failed expected 19: got:", value.(map[string]interface{})["Count"])
	}
}

func TestBigWrite(t *testing.T) {
	dbfile := filepath.Join(os.TempDir(), "test-file-bw.db")

	file, err := reset(dbfile)

	if err != nil {
		t.Fatal("Reset failed:", err)
		return
	}

	var db = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db.Open())

	if err != nil {
		t.Fatal("Open failed:", err)
		return
	}

	var doc struct {
		Data string
	}

	doc.Data = "data"
	db.Set("test", doc)

	doc.Data = fmt.Sprintf("%x", make([]byte, 512))

	db.Set("test", doc)

	var db2 = NewJsonDb(util.NewFileStorage(file, os.ModePerm), true)

	err = util.AwaitError(db2.Open())

	if err != nil {
		t.Fatal("db2 Open failed:", err)
		return
	}

	if value, ok := db2.Get("test"); !ok || value.(map[string]interface{})["Data"] != doc.Data {
		t.Error("Failed to update doc")
	}
}
