package token_bucket

import (
	"errors"
	"math"
	"sync/atomic"
	"time"
)

// inspire from paper: Rate-Based Active Queue Management with Token Buckets
// Section: Rate-Based AQM with Token Buckets

var (
	ErrRateBasedTokenBucketW1OrW2Empty error = errors.New("RateBasedTokenBucket: w1 or w2 is empty")
	ErrRateBasedTokenBucketAlreadyDead error = errors.New("RateBasedTokenBucket: already dead")
)

type adjustFuncType func() bool

const (
	defaultTBW1 float64 = 32
	defaultTBW2 float64 = 64

	defaultTBMinFillupSize int64 = 50
	defaultTBMaxFillupSize int64 = 100
	defaultTBPeakBustSize  int64 = 1000

	defaultTBFrequent time.Duration = 2 * time.Second
)

type TokenBucket interface {
	Close() error
	Take(n int64) int64
	Wait(n int64) time.Duration
}

type rateBasedTokenBucket struct {
	w1 float64
	w2 float64

	minFillupSize int64
	maxFillupSize int64
	peakBustSize  int64
	freq          time.Duration

	adjustFunc adjustFuncType

	cir      float64
	tokens   int64
	closeing chan struct{}
}

func NewRateBasedTokenBucket(w1, w2 float64, minfs, maxfs, pbs int64, freq time.Duration, adjustFunc adjustFuncType) *rateBasedTokenBucket {
	if w1 == 0 || w2 == 0 {
		panic(ErrRateBasedTokenBucketW1OrW2Empty)
	}

	r := &rateBasedTokenBucket{
		w1:            w1,
		w2:            w1,
		minFillupSize: minfs,
		maxFillupSize: maxfs,
		peakBustSize:  pbs,
		freq:          freq,
		adjustFunc:    adjustFunc,
		tokens:        0, // default tokens is zero
		closeing:      make(chan struct{}),
	}

	// go fillup goroutine
	go r.fillup()

	return r
}

func NewDefaultRateBasedTokenBucket(adjustFunc adjustFuncType) *rateBasedTokenBucket {
	return NewRateBasedTokenBucket(
		defaultTBW1,
		defaultTBW2,
		defaultTBMinFillupSize,
		defaultTBMaxFillupSize,
		defaultTBPeakBustSize,
		defaultTBFrequent,
		adjustFunc,
	)
}

func (r *rateBasedTokenBucket) Close() error {
	close(r.closeing)
	return nil
}

func (r *rateBasedTokenBucket) Take(n int64) int64 {
	for {
		if tokens := atomic.LoadInt64(&r.tokens); tokens == 0 {
			return 0
		} else if n <= tokens {
			if !atomic.CompareAndSwapInt64(&r.tokens, tokens, tokens-n) {
				continue
			}
			return n
		} else if atomic.CompareAndSwapInt64(&r.tokens, tokens, 0) {
			return tokens
		}
	}
}

func (r *rateBasedTokenBucket) Put(n int64) int64 {
	for {
		if tokens := atomic.LoadInt64(&r.tokens); tokens == r.peakBustSize {
			return 0
		} else if tokens+n >= r.peakBustSize {
			if !atomic.CompareAndSwapInt64(&r.tokens, tokens, r.peakBustSize) {
				continue
			}
			return r.peakBustSize - tokens
		} else if atomic.CompareAndSwapInt64(&r.tokens, tokens, tokens+n) {
			return n
		}
	}
}

func (r *rateBasedTokenBucket) Wait(n int64) time.Duration {
	var rem int64
	if rem = n - r.Take(n); rem == 0 {
		return 0
	}

	var wait time.Duration
	for rem > 0 {
		sleep := r.wait(rem)
		wait += sleep
		time.Sleep(sleep)
		rem -= r.Take(rem)
	}
	return wait
}

func (r *rateBasedTokenBucket) wait(n int64) time.Duration {
	// FIXME: this is just a presume value due to last cir(commited information rate)
	return time.Duration(int64(math.Ceil(math.Min(float64(n), float64(r.peakBustSize)) / (r.cir / float64(r.freq)))))
}

func (r *rateBasedTokenBucket) fillup() {
	ticker := time.NewTicker(r.freq)
	defer ticker.Stop()

	// initialize cir is min fillup size
	r.cir = float64(r.minFillupSize)
	for {
		select {
		case <-ticker.C:
			// fill up tokens
			increase := r.adjustFunc()
			var cir float64
			if increase {
				cir = r.cir + float64(r.maxFillupSize)/r.w1
			} else {
				cir = r.cir * (1.0 - 1.0/r.w2)
			}
			r.cir = math.Max(float64(r.minFillupSize), math.Min(float64(r.maxFillupSize), cir))
			cir = math.Floor(.5 + r.cir)
			r.Put(int64(cir))
		case <-r.closeing:
			return
		}
	}
}
