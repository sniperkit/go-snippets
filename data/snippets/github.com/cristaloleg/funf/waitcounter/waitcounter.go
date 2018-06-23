package waitcounter

import (
	"sync"
	"sync/atomic"
)

// WaitCounter ...
type WaitCounter struct {
	count int64
	wg    sync.WaitGroup
}

// Get ...
func (w *WaitCounter) Get() int64 {
	return atomic.LoadInt64(&w.count)
}

// Inc ...
func (w *WaitCounter) Inc() int64 {
	w.wg.Add(1)
	return atomic.AddInt64(&w.count, 1)
}

// Dec ...
func (w *WaitCounter) Dec() int64 {
	count := atomic.AddInt64(&w.count, -1)
	w.wg.Done()
	return count
}

// Upd ...
func (w *WaitCounter) Upd(delta int) int64 {
	var count int64
	if delta > 0 {
		w.wg.Add(delta)
		count = atomic.AddInt64(&w.count, int64(delta))
	} else {
		count = atomic.AddInt64(&w.count, int64(-delta))
		w.wg.Add(delta)
	}
	return count
}

// IsActive ...
func (w *WaitCounter) IsActive() bool {
	return atomic.LoadInt64(&w.count) > 0
}

// Wait ...
func (w *WaitCounter) Wait() {
	w.wg.Wait()
}
