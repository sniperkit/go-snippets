package pagerank

import "sync/atomic"

type Base struct {
	traversals uint64
	nonstarter bool
}

func (b *Base) IsStarter() bool {
	return !b.nonstarter
}

func (b *Base) Traverse() {
	atomic.AddUint64(&(b.traversals), 1)
}

func (b *Base) Traversals() uint64 {
	return atomic.LoadUint64(&(b.traversals))
}
