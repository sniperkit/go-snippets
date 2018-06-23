package revision

import (
	"errors"
	"sort"
)

var (
	ErrRevisionNotIntersect error = errors.New("revision: merge revision not intersect")
	ErrRevisionNotConcate   error = errors.New("revision: merge revision not concate")
)

type Revision struct {
	vclock int64
	off    uint64
	data   []byte
}

func NewRevision(v int64, o uint64, data []byte) *Revision {
	return &Revision{v, o, data}
}

func (r *Revision) Len() int      { return len(r.data) }
func (r *Revision) Begin() uint64 { return r.off }
func (r *Revision) End() uint64   { return r.off + uint64(r.Len()) - 1 }

func (r *Revision) Intersect(l *Revision) bool {
	return !(l.End() < r.Begin() || l.Begin() > r.End())
}

func (r *Revision) Concate(l *Revision) error {
	if r.End()+1 != l.Begin() {
		return ErrRevisionNotConcate
	}
	// mark this revision be concated
	r.vclock = -1
	r.data = append(r.data, l.data...)
	return nil
}

func (r *Revision) Merge(l *Revision) error {
	if !r.Intersect(l) {
		return ErrRevisionNotIntersect
	}
	correctBegin := min(r.Begin(), l.Begin())
	correctEnd := max(r.End(), l.End())
	mergeData := make([]byte, correctEnd-correctBegin+1)
	if r.vclock < l.vclock {
		r.vclock = l.vclock
		// merge r
		lOffInData := r.Begin() - correctBegin
		copy(mergeData[lOffInData:], r.data)
		// merge l
		lOffInData = l.Begin() - correctBegin
		copy(mergeData[lOffInData:], l.data)
	} else {
		// merge l
		lOffInData := l.Begin() - correctBegin
		copy(mergeData[lOffInData:], l.data)
		// merge r
		lOffInData = r.Begin() - correctBegin
		copy(mergeData[lOffInData:], r.data)
	}
	r.data = mergeData
	r.off = correctBegin
	return nil
}

type SparseRevisions struct {
	revisions []*Revision
}

func NewEmptySparseRevisions() *SparseRevisions {
	return &SparseRevisions{
		revisions: []*Revision{},
	}
}

func (s *SparseRevisions) Insert(r *Revision) *SparseRevisions {
	s.merge(r)
	// sort by vclock
	sort.Sort(RevisionsSortByVClock(s.revisions))
	return s
}

func (s *SparseRevisions) merge(rc *Revision) {
	var begin uint64 = rc.Begin()
	var end uint64 = rc.End()
	var revisions []*Revision
	for _, r := range s.revisions {
		if r.End() < begin || r.Begin() > end {
			revisions = append(revisions, r)
			continue
		}
		rc.Merge(r)
	}
	s.revisions = append(revisions, rc)
}

func (s *SparseRevisions) Compaction() {
	// sort by off
	sort.Sort(RevisionsSortByOff(s.revisions))

	var revisions []*Revision = s.revisions
	index := 0
	for {
		if index+1 >= len(revisions) {
			break
		}
		curr := revisions[index]
		next := revisions[index+1]
		if curr.End()+1 == next.Begin() {
			curr.Concate(next)
			revisions = append(revisions[0:index+1], revisions[index+2:]...)
		} else {
			index++
		}
	}
	s.revisions = revisions
	//FIXME: already lost vclock control
}

// for revision sorted
type RevisionsSortByVClock []*Revision

func (r RevisionsSortByVClock) Len() int           { return len(r) }
func (r RevisionsSortByVClock) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RevisionsSortByVClock) Less(i, j int) bool { return r[i].vclock < r[j].vclock }

// sorted by offset
type RevisionsSortByOff []*Revision

func (r RevisionsSortByOff) Len() int           { return len(r) }
func (r RevisionsSortByOff) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RevisionsSortByOff) Less(i, j int) bool { return r[i].off < r[j].off }

func max(x, y uint64) uint64 {
	if x < y {
		return y
	}
	return x
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
