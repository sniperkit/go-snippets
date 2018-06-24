package core

import (
	"github.com/datatogether/sql_datastore"
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestUncrawlableStorage(t *testing.T) {

	store := datastore.NewMapDatastore()

	p := &Uncrawlable{Url: "http://www.epa.gov"}
	if err := p.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	p.Comments = "new comment"
	if err := p.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	p2 := &Uncrawlable{Id: p.Id}
	if err := p2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if !p2.Created.Equal(p.Created) {
		t.Errorf("created doesn't match: %s != %s", p2.Created.String(), p.Created.String())
	}

	if !p2.Updated.Equal(p.Updated) {
		t.Errorf("updated doesn't match: %s != %s", p2.Updated.String(), p.Updated.String())
	}

	if err := p.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestUncrawlableSQLStorage(t *testing.T) {
	defer resetTestData(appDB, "uncrawlables")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Uncrawlable{}); err != nil {
		t.Error(err)
		return
	}

	p := &Uncrawlable{Url: "http://www.epa.gov"}
	if err := p.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	p.Comments = "new comment"
	if err := p.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	p2 := &Uncrawlable{Url: "http://www.epa.gov"}
	if err := p2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if !p2.Created.Equal(p.Created) {
		t.Errorf("created doesn't match: %s != %s", p2.Created.String(), p.Created.String())
	}

	if !p2.Updated.Equal(p.Updated) {
		t.Errorf("updated doesn't match: %s != %s", p2.Updated.String(), p.Updated.String())
	}

	if err := p.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}
