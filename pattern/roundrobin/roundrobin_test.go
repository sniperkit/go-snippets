package roundrobin

import (
	"log"
	"testing"
)

func BenchmarkRoundRobin(b *testing.B) {
	// <setup code>
	// ...

	b.Run("TestSequentialWithMutex", doBenchmarkRoundRobin)
	b.Run("TestParallelWithMutext", doBenchmarkRoundRobinParallel)
	b.Run("TestSequentialWithChannel", doBenchmarkRobinHood)
	b.Run("TestParallelWithChannel", doBenchmarkRobinHoodParallel)

	// <tear-down code>
	// ...
}

func doBenchmarkRoundRobin(b *testing.B) {
	rr := NewRoundRobin([]int{1, 2, 3, 4, 5, 6})
	for k := 0; k < b.N; k++ {
		rr.Get()
	}
}

func doBenchmarkRoundRobinParallel(b *testing.B) {
	rr := NewRoundRobin([]int{1, 2, 3, 4, 5, 6})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rr.Get()
		}
	})
}

func doBenchmarkRobinHood(b *testing.B) {
	rr := NewRobinHood([]int{1, 2, 3, 4, 5, 6})
	for k := 0; k < b.N; k++ {
		rr.Get()
	}
}

func doBenchmarkRobinHoodParallel(b *testing.B) {
	rr := NewRobinHood([]int{1, 2, 3, 4, 5, 6})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rr.Get()
		}
	})
}

func TestRobinHood(t *testing.T) {
	resources := []int{1, 2, 3, 4, 5}
	rr := NewRobinHood(resources)

	for range resources {
		log.Println(rr.Get())
	}
}
