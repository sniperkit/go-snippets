/*
* @Author: souravray
* @Date:   2015-01-25 03:27:42
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 03:27:54
 */

/*
* @Author: souravray
* @Date:   2015-01-25 01:21:27
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 03:26:40
 */
// taken from github.com/souravray/polybolos with permission

package stacker

import (
	"sync/atomic"
	"time"
)

type Bucket struct {
	tokens     int32
	usedTokens int32
	capacity   int32
	freq       time.Duration
	stop       chan bool
}

func NewBucket(capacity int32, rate int32) (b *Bucket, err error) {

	b = &Bucket{capacity: capacity}
	minFreq := time.Duration(1e9 / int64(capacity))
	freq := time.Duration(1e9 / int64(rate))
	if minFreq > freq {
		b.freq = minFreq
	} else {
		b.freq = freq
	}
	return
}

func (b *Bucket) setupUsedTokens(delta int32) {
	for {
		usedTokens := atomic.LoadInt32(&b.usedTokens)
		if !atomic.CompareAndSwapInt32(&b.usedTokens, usedTokens, usedTokens+delta) {
			continue
		} else {
			break
		}
	}
	return
}

func (b *Bucket) setdownUsedTokens(delta int32) {
	for {
		if usedTokens := atomic.LoadInt32(&b.usedTokens); usedTokens < delta {
			if !atomic.CompareAndSwapInt32(&b.usedTokens, usedTokens, 0) {
				continue
			} else {
				break
			}
		} else {
			if !atomic.CompareAndSwapInt32(&b.usedTokens, usedTokens, usedTokens-delta) {
				continue
			} else {
				break
			}
		}
	}
	return
}

func (b *Bucket) Put() (success bool) {
	for {
		tokens := atomic.LoadInt32(&b.tokens)
		usedTokens := atomic.LoadInt32(&b.usedTokens)
		if tokens+usedTokens < b.capacity {
			if !atomic.CompareAndSwapInt32(&b.tokens, tokens, tokens+1) {
				continue
			} else {
				break
			}
		} else {
			break
		}
	}
	return
}

func (b *Bucket) Take(n int32) <-chan int32 {
	c := make(chan int32)
	go func(c chan int32, b *Bucket, n int32) {
		var tokens int32
		defer close(c)
	TryToTake:
		for {
			if tokens = atomic.LoadInt32(&b.tokens); tokens == 0 {
				break
			} else if n <= tokens {
				if !atomic.CompareAndSwapInt32(&b.tokens, tokens, tokens-n) {
					continue
				} else {
					break
				}
			} else {
				if !atomic.CompareAndSwapInt32(&b.tokens, tokens, 0) {
					continue
				} else {
					break
				}
			}
		}

		if tokens > 0 {
			b.setupUsedTokens(tokens)
			c <- tokens
		} else {
			time.Sleep(time.Duration(n * int32(b.freq.Nanoseconds())))
			goto TryToTake
		}
	}(c, b, n)
	return c
}

func (b *Bucket) Spend() {
	b.setdownUsedTokens(1)
	return
}

func (b *Bucket) Fill() {

	b.stop = make(chan bool)
	ticker := time.NewTicker(b.freq)
	for _ = range ticker.C {
		select {
		case <-b.stop:
			ticker.Stop()
			return
		default:
			b.Put()
		}
	}
}

func (b *Bucket) Close() {

	b.stop <- true
	defer close(b.stop)
	return
}
