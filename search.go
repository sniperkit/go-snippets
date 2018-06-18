package main

import (
	"github.com/pmylund/go-bitset"
	"sort"
)

type docScore struct {
	doc   uint32
	score uint64
}

type docScores struct {
	store []*docScore
}

type queryAnswer struct {
	query  *bitset.Bitset64
	answer chan []*docScore
}

func (self docScores) Len() int { return len(self.store) }

func (self docScores) Swap(i, j int) {
	self.store[i], self.store[j] = self.store[j], self.store[i]
}

func (self *docScores) Add(ds *docScore) {
	if self.store == nil {
		self.store = make([]*docScore, 0, 4)
	}
	n := len(self.store)
	if n+1 > cap(self.store) {
		s := make([]*docScore, n, 2*n+1)
		copy(s, self.store)
		self.store = s
	}
	self.store = self.store[0 : n+1]
	self.store[n] = ds
}

type byScore struct{ docScores }

func (self byScore) Less(i, j int) bool {
	return self.docScores.store[i].score < self.docScores.store[j].score
}

func (self *documents) Score(master *bitset.Bitset64, thresold_ float32) []*docScore {
	thresold := uint64(float64(thresold_) * float64(master.Count()))
	var results docScores = docScores{make([]*docScore, 0, 4)}
	for i, document := range self.tags {
		inter := master.Intersection(document)
		common := inter.Count()
		if common > thresold {
			for k, v := range docs.bonus {
				common += inter.Intersection(v).Count() * uint64(k)
			}
			results.Add(&docScore{uint32(i), common})
		}
	}
	sort.Sort(byScore{results})
	return results.store
}

//func (self *documents) Cluster(ids []uint32) {

/* A worker for handling search without fulling memory */
func StartSearch() {
	for {
		qa := <-searchQueue
		qa.answer <- docs.Score(qa.query, 0.2)
	}
}
