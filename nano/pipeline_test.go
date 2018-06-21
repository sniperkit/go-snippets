package nano

import (
	"fmt"
	"sync"
	"testing"
)

func TestPULLPUSH(T *testing.T) {

	count := 100000000
	var wg sync.WaitGroup
	Func := func(Args ...interface{}) ([]interface{}, error) {
		wg.Done()
		return nil, nil
	}
	P := Create("tcp://127.0.0.1:5050", Func)
	wg.Add(count)
	for i := 0; i < count; i++ {
		P.AddTask([]byte(fmt.Sprintf("%d", i)))
	}
	wg.Wait()
}
