package reactor

import (
	"reflect"
	"sort"
	"testing"
)

func TestByNuclearPrioritySort(t *testing.T) {
	var ns []*nuclear
	for i := 1; i < 10; i++ {
		ns = append(ns, NewNuclear("", i))
	}

	var reverse []*nuclear
	for _, nul := range ns {
		tmp := *nul
		reverse = append([]*nuclear{&tmp}, reverse...)
	}

	sort.Sort(ByNuclearPriority(ns))
	if !reflect.DeepEqual(ns, reverse) {
		t.Fatalf("nuclear sort by priority error")
	}
}
