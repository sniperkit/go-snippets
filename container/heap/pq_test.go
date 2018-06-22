// Copyright 2016-2018 Hyperchain Corp.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package heap

import (
	stdheap "container/heap"
	"fmt"
)

type Item struct {
	Value string
	Index int
	Priority float64
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {return len(pq)}
func (pq PriorityQueue) Less(i, j int) bool {return pq[i].Priority > pq[j].Priority}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index, pq[j].Index = i, j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}
func (pq *PriorityQueue) update(item *Item, value string, priority float64) {
	(*pq)[item.Index].Value = value
	(*pq)[item.Index].Priority = priority
	stdheap.Fix(pq, item.Index)
}

func ExamplePriorityQueue() {
	items := map[string]float64{
		"hello": 1.1,
		"world": 2.2,
	}
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			Value:    value,
			Priority: priority,
			Index:    i,
		}
		i++
	}
	stdheap.Init(&pq)
	// Insert a new item and then modify its priority.
	item := &Item{
		Value:    "orange",
		Priority: 1.0,
	}
	stdheap.Push(&pq, item)
	// Move the new inserted item to the head
	pq.update(item, item.Value, 5.0)
	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := stdheap.Pop(&pq).(*Item)
		fmt.Printf("%.1f:%s ", item.Priority, item.Value)
	}
	// Output:
	// 5.0:orange 2.2:world 1.1:hello
}


