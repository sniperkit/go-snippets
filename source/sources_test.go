package source

import (
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestListPrimers(t *testing.T) {
	egs := []string{"a", "b", "c"}

	store := datastore.NewMapDatastore()
	for _, val := range egs {
		s := &Source{Url: val, Title: val}
		if err := s.Save(store); err != nil {
			t.Error(err)
			return
		}
	}

	sources, err := ListSources(store, "", 20, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(sources) != len(egs) {
		t.Errorf("sources length mismatch: %d != %d", len(sources), len(egs))
	}

	// sources, err = ListSources(store, "", 20, len(egs))
	// if err != nil {
	// 	t.Errorf(err.Error())
	// }
	// if len(sources) != 0 {
	// 	t.Errorf("sources length mismatch")
	// }
}
