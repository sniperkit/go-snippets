package producer

import (
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestProducer(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancelFunc()

	var wg sync.WaitGroup
	wg.Add(2)
	// start producer
	go func() {
		defer wg.Done()
		produce(ctx)
	}()

	// start consumer
	go func() {
		defer wg.Done()
		consume(ctx)
	}()

	wg.Wait()
}
