package io_queue

import "time"

// A imlemented io queue inspired by seastar io queue.

// io priority class
type IOPriorityClass struct {
	val uint64
}

func NewIOPriorityClass(val uint64) *IOPriorityClass {
	return &IOPriorityClass{val}
}

func (pc *IOPriorityClass) Id() uint64 {
	return pc.val
}

const (
	MaxClasses int = 2048
)

type ShardId uint8

type Func func()

// io queue
type Queue struct {
	shardId ShardId

	capacity uint64

	// priority classes
	priorityClasses map[uint64]*priorityClassData
	fq              *FairQueue

	// registered shares
	registeredShares [MaxClasses]uint32
	registeredNames  [MaxClasses]string
}

func NewQueue(shardId ShardId, cap uint64) *Queue {
	q := new(Queue)
	q.shardId = shardId
	q.capacity = cap
	q.priorityClasses = make(map[uint64]*priorityClassData)
	q.fq = new(FairQueue)
	return q
}

func (q *Queue) registerOnePriorityClass(name string, shares uint32) IOPriorityClass {
}

func (q *Queue) findOrCreateClass(pc IOPriorityClass, ownerId ShardId) *priorityClassData {
	pclass, ok := q.priorityClasses[pc.Id()]
	if !ok {
		shares := q.registeredShares[pc.Id()]
		name := q.registeredNames[pc.Id()]

		newpc := q.fq.RegisterPriorityClass(shares)
		pclass = NewPriorityClassData(name, newpc, owner)
		q.priorityClasses[pc.Id()] = pclass
	}
	return pclass
}

func (q *Queue) Coordinator() ShardId   { return q.shardId }
func (q *Queue) Capacity() uint64       { return q.capacity }
func (q *Queue) QueuedRequests() uint64 { return q.fq.Waiters() }

func (q *Queue) RequestsCurrentlyExecuting() uint64 {
	return q.fq.RequestsCurrentlyExecuting()
}

func (q *Queue) NotifyRequestsFinished(finished uint64) {
	q.fq.NotifyRequestsFinished(finished)
}

func (q *Queue) PollQueue() {
	q.fq.DispatchRequests()
}

func (q *Queue) UpdateSharesForClass(pc IOPriorityClass, new uint32) {
	owner := 1
	pclass := q.findOrCreateClass(pc, owner)
	q.fq.UpdateShares(pclass.ptr, new)
}

func (q *Queue) Request(pc IOPriorityClass, len uint64, ops Func) Event {
	start := time.Now()
	owner := 1
	weight := 1 + len/(16<<10)

	pclass := q.findOrCreateClass(pc, owner)
	pclass.bytes += len
	pclass.ops++
	pclass.nrQueued++

	run := func() {
		//TODO: recover this function if someting goes wrong
		pclass.nrQueued--
		pclass.queueTime = time.Now().Sub(start)
		// function run
		ops()
	}

	q.fq.Queue(pclass.ptr, weight, run)
}

// priority class data
type priorityClassData struct {
	ptr       *PriorityClass
	bytes     uint64
	ops       uint64
	nrQueued  uint32
	queueTime time.Duration
}

func NewPriorityClassData(name string, ptr *PriorityClass, ownerId ShardId) *priorityClassData {
	return &priorityClassData{
		ptr:       ptr,
		bytes:     0,
		ops:       0,
		nrQueued:  0,
		queueTime: time.Second * 1,
	}
}
