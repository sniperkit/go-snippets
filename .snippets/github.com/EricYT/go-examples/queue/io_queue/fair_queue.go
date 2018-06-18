package io_queue

import (
	"math"
	"time"

	"github.com/facebookgo/pqueue"
)

// priority class
type PriorityClass struct {
	shares      uint32
	accumulated float64

	queue  *CircularBuffer
	queued bool
}

func NewPriorityClass(shares uint32) *PriorityClass {
	pc := new(PriorityClass)
	pc.shares = shares
	pc.accumulated = 0
	pc.queued = false

	pc.queue = NewCircularBuffer(10)

	return pc
}

func (pc *PriorityClass) UpdateShares(shares uint32) {
	pc.shares = shares
}

func (pc *PriorityClass) Shares() uint32 {
	return pc.shares
}

type FairQueue struct {
	requestsExecuting uint32
	requestsQueued    uint32
	capacity          uint64
	base              int
	tau               time.Duration

	// priority queue
	handles    pqueue.PriorityQueue
	allClasses []*PriorityClass
}

func NewFairQueue(cap uint64) *FairQueue {
	fq := new(FairQueue)
	fq.requestsExecuting = 0
	fq.requestsQueued = 0
	fq.capacity = cap
	fq.base = time.Now().Nanosecond()
	fq.tau = time.Microsecond(time.Millisecond * 100)

	fq.handles = pqueue.New(cap)
	fq.allClasses = make([]*PriorityClass)

	return fq
}

func (fq *FairQueue) PushPriorityClass(pc *PriorityClass) {
	if !pc.queued {
		fq.handles.Push(&pqueue.Item{Value: pc, Priority: pc.accumulated})
		pc.queued = true
	}
}

func (fq *FairQueue) PopPriorityClass() *PriorityClass {
	if fq.handles.Len() == 0 {
		return nil
	}

	item := fq.handles.Pop().(*pqueue.Item)
	hnd := item.Value.(*PriorityClass)
	if !hnd.queued {
		panic("handle priority class queued is false")
	}
	hnd.queued = false
	return hnd
}

func (fq *FairQueue) normalizeFactor() float64 {
	return math.SmallestNonzeroFloat64
}

func (fq *FairQueue) normalizeStats() {
	timeDelta := math.Log(fq.normalizeFactor()) * float64(fq.tau)
	fq.base -= int(timeDelta)
	for _, pc := range fq.allClasses {
		pc.accumulated *= fq.normalizeFactor()
	}
}

func (fq *FairQueue) RegisterPriorityClass(shares uint32) *PriorityClass {
	pc := NewPriorityClass(shares)
	fq.allClasses = append(fq.allClasses, pc)
	return pc
}

func (fq *FairQueue) UnregisterPriorityClass(pc *PriorityClass) {
	var withoutThePc []*PriorityClass
	for _, p := range fq.allClasses {
		if p != pc {
			withoutThePc = append(withoutThePc, p)
		}
	}
	fq.allClasses = withoutThePc
}

func (fq *FairQueue) Waiters() {
	return fq.requestsQueued
}

func (fq *FairQueue) RequestsCurrentlyExecuting() {
	return fq.requestsExecuting
}

func (fq *FairQueue) Queue(pc *PriorityClass, weight uint32, fun Func) {
	fq.PushPriorityClass(pc)
	pc.queue.PushBack(NewRequest(fun, weight))
	fq.requestsQueued++
}

func (fq *FairQueue) NotifyRequestsFinished(finished uint32) {
	fq.requestsExecuting -= finished
}

func (fq *FairQueue) DispatchRequests() {
	for (fq.requestsQueued != 0) && (fq.requestsExecuting < fq.capacity) {
		var h *PriorityClass
		for {
			h = fq.PopPriorityClass()
			if !h.queue.Empty() {
				break
			}
		}

		req := h.queue.PopFront()
		fq.requestsExecuting++
		fq.requestsQueued--

		delta := (time.Now().Nanosecond() - fq.base) / 1000 // convert to microsecond
		reqCost := float64(req.weight) / float64(h.shares)
		cost := math.Exp(float64(1)/float64(fq.tau)*delta) * reqCost
		nextAccumulated := h.accumulated + cost
		for math.IsInf(nextAccumulated, 0) {
			fq.normalizeStats()
			// If we renormalized, out time base will have changed. This should happen very infrequently
			delta = (time.Now().Nanosecond() - fq.base) / 1000 // convert to microsecond
			cost = math.Exp(float64(1)/float64(fq.tau)*delta) * reqCost
			nextAccumulated = h.accumulated + cost
		}
		h.accumulated = nextAccumulated

		if !h.queue.Empty() {
			fq.PushPriorityClass(h)
		}

		req.fun()
	}
}

func (fq *FairQueue) UpdateShares(pc *PriorityClass, shares uint32) {
	pc.UpdateShares(shares)
}
