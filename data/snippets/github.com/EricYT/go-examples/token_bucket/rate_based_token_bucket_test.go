package token_bucket

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var index int

func adjustFunc() bool {
	defer func() { index++ }()
	if index%5 != 0 {
		return true
	}
	return false
}

func TestTokenBucket(t *testing.T) {
	tb := NewRateBasedTokenBucket(50, 20, 10, 100, 10, time.Millisecond*10, adjustFunc)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			//			select {
			//			case <-time.After(time.Millisecond * 500):
			now := time.Now()
			//log.Printf("Time now: %s tokens: %d\n", time.Now(), 3)
			escape := tb.Wait(1)
			fmt.Printf("Time escape: %s escape: %d\n", time.Now().Sub(now), escape)
			//			}
		}
	}()
	wg.Wait()
}
