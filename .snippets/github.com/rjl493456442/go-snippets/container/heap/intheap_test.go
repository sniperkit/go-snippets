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

type IntHeap []int

func (h IntHeap) Len() int {return len(h)}
func (h IntHeap) Less(i, j int) bool {return h[i] < h[j]}
func (h IntHeap) Swap(i, j int) {h[i], h[j] = h[j], h[i]}
func (h *IntHeap) Push(x interface{}) {
	// Because we will modify the slice length here, so a pointer will be passed
	// as the receiver instead of the slice itself.
	*h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
	// TODO The reason to create a copy old here is?
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func ExampleSimpleHeap() {
	// Note only the IntHeap pointer implements the interface.
	h := &IntHeap{2, 1, 5}
	stdheap.Init(h)
	// Print the first element
	fmt.Println((*h)[0])

	// Insert a minimum element
	stdheap.Push(h, 0)
	fmt.Println((*h)[0])

	// Modify the underlying storage directly
	(*h)[0] = 100
	stdheap.Fix(h, 0)
	fmt.Println((*h)[0])

	// Pop the minimum element from the heap
	x := stdheap.Pop(h)
	fmt.Println(x.(int))

	// Output:
	// 1
	// 0
	// 1
	// 1
}


