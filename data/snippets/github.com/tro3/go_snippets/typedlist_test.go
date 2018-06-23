package snippets

import (
	"testing"
)

func TestPush(t *testing.T) {
	h := NewTestHelper(t)
	a := TypeAList{}
	a.Push(TYPEA{3})
	a.Push(TYPEA{4})
	a.Push(TYPEA{5})
	h.Assert(len(a) == 3, "Incorrect Push length")
}

func TestPop(t *testing.T) {
	h := NewTestHelper(t)
	a := TypeAList{TYPEA{5}, TYPEA{4}, TYPEA{3}}
	x := a.Pop()
	h.Assert(len(a) == 2, "Incorrect Pop length")
	h.Assert(x.age == 3, "Incorrect Pop value")
}

func TestInsert(t *testing.T) {
	h := NewTestHelper(t)
	a := TypeAList{TYPEA{5}, TYPEA{4}, TYPEA{3}}
	a.Insert(1, TYPEA{6})
	h.Assert(len(a) == 4, "Incorrect Insert length")
	h.Assert(a[1].age == 6, "Incorrect Insert value")
}

func TestRemove(t *testing.T) {
	h := NewTestHelper(t)
	a := TypeAList{TYPEA{5}, TYPEA{4}, TYPEA{3}}
	b := a.Remove(1)
	h.Assert(len(a) == 2, "Incorrect Remove length")
	h.Assert(b.age == 4, "Incorrect Remove value")
}

func TestIndex(t *testing.T) {
	h := NewTestHelper(t)
	a := TypeAList{TYPEA{5}, TYPEA{4}, TYPEA{3}}
	h.Assert(a.Index(a[1]) == 1, "Incorrect Index value")
}

func TestFilter(t *testing.T) {
	h := NewTestHelper(t)
	a := TypeAList{TYPEA{5}, TYPEA{4}, TYPEA{6}, TYPEA{3}}
	b := a.Filter(func(val TYPEA) bool {
		return val.age > 4
	})
	h.Assert(len(b) == 2, "Incorrect Filter length")
	h.Assert(b[1].age == 6, "Incorrect Filter value")
}

func TestSort(t *testing.T) {
	h := NewTestHelper(t)
	a := TypeAList{TYPEA{5}, TYPEA{4}, TYPEA{3}}
	a.Sort(func(a, b TYPEA) bool {
		return a.age < b.age
	})
	h.Assert(a[0].age == 3, "Incorrect Pop value")
}
