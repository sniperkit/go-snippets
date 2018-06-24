package core

import (
	"github.com/datatogether/sql_datastore"
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestPrimerStorage(t *testing.T) {
	store := datastore.NewMapDatastore()

	p := &Primer{Title: "Test Primer", Description: "test primer description!"}
	if err := p.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	p.Description = "new description"
	if err := p.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	p2 := &Primer{Id: p.Id}
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

func TestPrimerSQLStorage(t *testing.T) {
	defer resetTestData(appDB, "primers", "sources")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Primer{}); err != nil {
		t.Error(err.Error())
		return
	}

	p := &Primer{Title: "Test Primer", Description: "test primer description!"}
	if err := p.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	p.Description = "new description"
	if err := p.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	p2 := &Primer{Id: p.Id}
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

func TestPrimerReadSubprimers(t *testing.T) {
	p := &Primer{Id: "5b1031f4-38a8-40b3-be91-c324bf686a87"}
	if err := p.ReadSubPrimers(appDB); err != nil {
		t.Error(err.Error())
	}
}

func TestPrimerReadSources(t *testing.T) {
	p := &Primer{Id: "5b1031f4-38a8-40b3-be91-c324bf686a87"}
	if err := p.ReadSources(appDB); err != nil {
		t.Error(err.Error())
	}
}

func TestPrimerCalcStats(t *testing.T) {
	p := &Primer{Id: "5b1031f4-38a8-40b3-be91-c324bf686a87"}
	if err := p.CalcStats(appDB); err != nil {
		t.Error(err.Error())
	}
}
