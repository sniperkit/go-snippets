package main

import (
	"context"
	"fmt"
	"time"
)

//
// Plugin that Counts numbers
type NumberPlugin struct{}

func (n *NumberPlugin) DoWork(ctx context.Context, e Event, id int, resultChan chan interface{}) {

	var num int

	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("\t%s (NumberPlugin worker %d) only made it to %d when event cancel() was called\n", e.Name, id, num)
			return
		default:
			// Do some work
			num++
			time.Sleep(time.Millisecond * time.Duration(randInt(50, 150)))
		}
	}

	// Finished
	resultChan <- fmt.Sprintf("NumberPlugin.worker-%d = %d", id, num)

}

// MathPlugin that does math
type MathPlugin struct{}

// DoWork well
func (n *MathPlugin) DoWork(ctx context.Context, e Event, id int, resultChan chan interface{}) {

	mathChan := make(chan []int64)

	go func() {
		time.Sleep(time.Millisecond * 500)
		mathChan <- []int64{9, 11}
	}()

	// Which ever comes first...
	select {
	case <-ctx.Done():
		fmt.Printf("\t%s (MathPlugin worker %d) failed to finish maths\n", e.Name, id)
	case v := <-mathChan:
		resultChan <- fmt.Sprintf("MathPlugin.worker-%d = %v", id, v)
	}

}
