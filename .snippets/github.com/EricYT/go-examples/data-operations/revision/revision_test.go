package revision

import (
	"bytes"
	"log"
	"sort"
	"testing"
)

func revisionEqual(r *Revision, e *Revision) bool {
	return r.vclock == e.vclock && r.off == e.off && bytes.Equal(r.data, e.data)
}

func revisionsCompare(t *testing.T, r []*Revision, e []*Revision) {
	if len(r) != len(e) {
		t.Errorf("revision compare left length(%d) not equal right length(%d)", len(r), len(e))
		return
	}
	index := 0
	for index < len(r) {
		if !revisionEqual(r[index], e[index]) {
			t.Errorf("revision compare left(%v) not equal right(%v)", r[index], e[index])
			break
		}
		index++
	}
}

func TestRevisionsSort(t *testing.T) {
	revisions := []*Revision{
		{9, 1, []byte{}},
		{1, 23, []byte{}},
		{3, 123, []byte{}},
		{2, 3, []byte{}},
	}
	sortedRevisionsByVClock := []*Revision{
		{1, 23, []byte{}},
		{2, 3, []byte{}},
		{3, 123, []byte{}},
		{9, 1, []byte{}},
	}
	sortedRevisionsByOff := []*Revision{
		{9, 1, []byte{}},
		{2, 3, []byte{}},
		{1, 23, []byte{}},
		{3, 123, []byte{}},
	}
	// sort revisions
	sort.Sort(RevisionsSortByVClock(revisions))
	// compare
	revisionsCompare(t, revisions, sortedRevisionsByVClock)

	// sort revisions
	sort.Sort(RevisionsSortByOff(revisions))
	// compare
	revisionsCompare(t, revisions, sortedRevisionsByOff)
}

/*

vc 0 1 2 3 4 5 6 7 8 9
0  - - - - 1 2 3 - - - -

1  4 5 6 7 - - - - - - -
3  - 1 3 5 7 - - - - - -
5  - - 2 2 6 3 - - - - -
2  - - - 0 5 1 1 - - - -
9  - - - - 3 1 6 3 - - -
7  - - - - - 8 2 9 3 - -
8  - - - - - - 9 3 6 3 -
10 - - - - - - - 4 9 9 9

   4 1 2 2 3 1 6 4 9 9 9
*/

var (
	r *Revision = &Revision{0, 4, []byte{1, 2, 3}}

	r1 *Revision = &Revision{1, 0, []byte{4, 5, 6, 7}}
	r2 *Revision = &Revision{3, 1, []byte{1, 3, 5, 7}}
	r3 *Revision = &Revision{5, 2, []byte{2, 2, 6, 3}}
	r4 *Revision = &Revision{2, 3, []byte{0, 5, 1, 1}}
	r5 *Revision = &Revision{9, 4, []byte{3, 1, 6, 3}}
	r6 *Revision = &Revision{7, 5, []byte{8, 2, 9, 3}}
	r7 *Revision = &Revision{8, 6, []byte{9, 3, 6, 3}}
	r8 *Revision = &Revision{10, 7, []byte{4, 9, 9, 9}}
)

func TestRevisionIntersect(t *testing.T) {
	if r.Intersect(r1) {
		t.Errorf("revision r(%v) should not intersect with r1(%v)", r, r1)
	}
	if !r.Intersect(r2) {
		t.Errorf("revision r(%v) should intersect with r2(%v)", r, r2)
	}
	if !r.Intersect(r3) {
		t.Errorf("revision r(%v) should intersect with r3(%v)", r, r3)
	}
	if !r.Intersect(r4) {
		t.Errorf("revision r(%v) should intersect with r4(%v)", r, r4)
	}
	if !r.Intersect(r5) {
		t.Errorf("revision r(%v) should intersect with r5(%v)", r, r5)
	}
	if !r.Intersect(r6) {
		t.Errorf("revision r(%v) should intersect with r6(%v)", r, r6)
	}
	if !r.Intersect(r7) {
		t.Errorf("revision r(%v) should intersect with r7(%v)", r, r7)
	}
	if r.Intersect(r8) {
		t.Errorf("revision r(%v) should not intersect with r8(%v)", r, r8)
	}
}

func TestReversionMerge(t *testing.T) {
	if err := r.Merge(r1); err == nil {
		t.Errorf("revision r(%v) is not intersect with r1(%v)", r, r1)
	}
	// merge
	r1 := NewRevision(1, 3, []byte{7, 8, 9, 10})
	r2 := NewRevision(3, 1, []byte{1, 3, 5, 7})
	if err := r1.Merge(r2); err != nil {
		t.Errorf("revision r1(%v) merge with r2(%v) error: %s", r1, r2, err)
	}
	r := &Revision{3, 1, []byte{1, 3, 5, 7, 9, 10}}
	if !revisionEqual(r1, r) {
		t.Errorf("revision r1([7 8 9 10]) merge(%v) not equal (%v)", r2, r)
	}

	// merge
	r1 = NewRevision(3, 3, []byte{7, 8, 9, 10})
	r2 = NewRevision(1, 1, []byte{1, 3, 5, 7})
	if err := r1.Merge(r2); err != nil {
		t.Errorf("revision r1(%v) merge with r2(%v) error: %s", r1, r2, err)
	}
	r = &Revision{3, 1, []byte{1, 3, 7, 8, 9, 10}}
	if !revisionEqual(r1, r) {
		t.Errorf("revision r1([7 8 9 10]) merge(%v) not equal (%v)", r2, r)
	}

	// merge
	r1 = NewRevision(3, 3, []byte{7})
	r2 = NewRevision(1, 1, []byte{1, 3, 5, 7})
	if err := r1.Merge(r2); err != nil {
		t.Errorf("revision r1(%v) merge with r2(%v) error: %s", r1, r2, err)
	}
	r = &Revision{3, 1, []byte{1, 3, 7, 7}}
	if !revisionEqual(r1, r) {
		t.Errorf("revision r1([7]) merge(%v) not equal (%v)", r2, r)
	}
}

func TestSparseRevisionsInsert1(t *testing.T) {
	sr := NewEmptySparseRevisions()
	sr.Insert(r)
	sr.Insert(r1)
	sr.Insert(r2)
	sr.Insert(r3)
	sr.Insert(r4)
	sr.Insert(r5)
	sr.Insert(r6)
	sr.Insert(r7)
	sr.Insert(r8)
	displayRevisions(t, sr.revisions)
	r = &Revision{10, 0, []byte{4, 1, 2, 2, 3, 1, 6, 4, 9, 9, 9}}
	if !revisionEqual(sr.revisions[0], r) {
		t.Errorf("revision sparse insert result: %v not equal right %v", sr.revisions[0], r)
	}
}

/*
vc 0 1 2 3 4 5 6 7 8 9

1  1 2 3 - - - - - - - -
3  - - - 0 6 3 - - - - -
2  - - - - - - - 9 1 9 -
9  - - - - 3 2 - - - - -

   1 2 3 0 3 2 - 9 1 9 -

*/

func TestSparseRevisionsInsert2(t *testing.T) {
	sr := NewEmptySparseRevisions()
	sr.Insert(&Revision{1, 0, []byte{1, 2, 3}})
	sr.Insert(&Revision{3, 3, []byte{0, 6, 3}})
	sr.Insert(&Revision{2, 7, []byte{9, 1, 9}})
	sr.Insert(&Revision{9, 4, []byte{3, 2}})
	sr.Compaction()
	displayRevisions(t, sr.revisions)
	r1 = &Revision{-1, 0, []byte{1, 2, 3, 0, 3, 2}}
	r2 = &Revision{2, 7, []byte{9, 1, 9}}
	if !revisionEqual(sr.revisions[0], r1) {
		t.Errorf("revision sparse revision 0(%v) not equal r1(%v)", sr.revisions[0], r1)
	}
	if !revisionEqual(sr.revisions[1], r2) {
		t.Errorf("revision sparse revision 1(%v) not equal r2(%v)", sr.revisions[1], r2)
	}
}

func displayRevisions(t *testing.T, rs []*Revision) {
	for index, r := range rs {
		log.Printf("revision index: %d revision: %v\n", index, r)
	}
	log.Println()
}
