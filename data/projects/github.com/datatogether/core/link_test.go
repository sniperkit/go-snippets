package core

import (
	"github.com/datatogether/sql_datastore"
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestLinkStorage(t *testing.T) {
	store := datastore.NewMapDatastore()

	l := &Link{Src: &Url{Url: "http://www.epa.gov"}, Dst: &Url{Url: "http://www.epa.gov"}}
	if err := l.Insert(store); err != nil {
		t.Error(err.Error())
		return
	}

	if err := l.Update(store); err != nil {
		t.Error(err.Error())
		return
	}

	l2 := &Link{
		Src: &Url{Url: "http://www.epa.gov"},
		Dst: &Url{Url: "http://www.epa.gov"},
	}
	if err := l2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if !l2.Created.Equal(l.Created) {
		t.Errorf("created doesn't match: %s != %s", l2.Created.String(), l.Created.String())
	}

	if !l2.Updated.Equal(l.Updated) {
		t.Errorf("updated doesn't match: %s != %s", l2.Updated.String(), l.Updated.String())
	}

	if err := l.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestLinkSQLStorage(t *testing.T) {
	defer resetTestData(appDB, "links")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Link{}); err != nil {
		t.Error(err.Error())
		return
	}

	l := &Link{Src: &Url{Url: "http://www.epa.gov"}, Dst: &Url{Url: "http://www.epa.gov"}}
	if err := l.Insert(store); err != nil {
		t.Error(err.Error())
		return
	}

	if err := l.Update(store); err != nil {
		t.Error(err.Error())
		return
	}

	l2 := &Link{
		Src: &Url{Url: "http://www.epa.gov"},
		Dst: &Url{Url: "http://www.epa.gov"},
	}
	if err := l2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if !l2.Created.Equal(l.Created) {
		t.Errorf("created doesn't match: %s != %s", l2.Created.String(), l.Created.String())
	}

	if !l2.Updated.Equal(l.Updated) {
		t.Errorf("updated doesn't match: %s != %s", l2.Updated.String(), l.Updated.String())
	}

	if err := l.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}

}
