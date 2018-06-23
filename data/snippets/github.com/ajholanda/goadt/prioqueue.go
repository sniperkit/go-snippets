// Package adt implements some fundamental abstract data types.

package adt

// Implementation min-max priority queue
type PQType bool

// MAX sets a priority queue where greater values
// are related to greater priorities.
// MIN sets lower values related to greater priorities.
const (
	MIN PQType = false
	MAX = true
)

// HeapNode is the unit of heap, index is the position
// in the heap, prio is the priority, key is the node
//identification and value is the information associated
// with the key
type HeapNode struct {
	index int
	prio  int
	key, value   interface{}
}

// Retrieve the index of heap node h
func (h *HeapNode) Index() int {
	return h.index
}

// Retrieve the priority of heap node h
func (h *HeapNode) Prio() int {
	return h.prio
}

// Retrieve the key of heap node h
func (h *HeapNode) Key() interface{} {
	return h.key
}

// Retrieve the value associated with the key
// in the heap node h
func (h *HeapNode) Value() interface{} {
	return h.value
}

// PriorityQueue is composed by heap nodes,
// the number of heap nodes stored in size, the type
// of the priority queue, maximum or minimum, and a
// map of keys to locate a heap node using its key.
type PriorityQueue struct {
	heap []*HeapNode // R[0] is ignored, so R[1..size]
	size int
	t PQType
	keyMap map[interface{}]*HeapNode
}

// Return the size of priority queue pq
func (pq *PriorityQueue) Size() int {
	return pq.size
}

// Verify if the priority queue pq is empty.
func (pq *PriorityQueue) IsEmpty() bool {
	return pq.size <= 0
}

// Return the type of priority queue, maximum or minimum.
func (pq *PriorityQueue) Type() PQType {
	return pq.t
}

// NewPriorityQueue returns an empty priority queue with its
// data initialized.
func NewPriorityQueue(t PQType) *PriorityQueue {
	return &PriorityQueue{make([]*HeapNode, 2), 0, t, make(map[interface{}]*HeapNode, 2)}
}

func maxgreateridx(pq *PriorityQueue, i, j int) bool {
	return pq.heap[i].prio > pq.heap[j].prio
}

func mingreateridx(pq *PriorityQueue, i, j int) bool {
	return pq.heap[i].prio < pq.heap[j].prio
}

// greater compares priorities of i and j and returns true if priority
// of i is greater than priority of j.
func (pq *PriorityQueue) greater(i, j int) bool {
	if pq.t==MAX {
		return maxgreateridx(pq, i, j)
	} else {//MIN
		return mingreateridx(pq, i, j)
	}
}

func maxgreaterprio(xprio, yprio int) bool {
	return xprio > yprio
}

func mingreaterprio( xprio, yprio int) bool {
	return xprio < yprio
}

func (pq *PriorityQueue) greaterprio(xprio, yprio int) bool {
	if pq.t==MAX {
		return maxgreaterprio(xprio, yprio)
	} else {//MIN
		return mingreaterprio(xprio, yprio)
	}
}

// The following functions and methods were adapted from "Algorithms"
// by Sedgewick and Wayne (4th ed) and "Estruturas de Dados e Seus
// Algoritmos" by Szwarcfiter and Markenzon

// swap exchanges heap node located at i with heap node locate at j.
func (pq *PriorityQueue) swap(i, j int) {
	pq.heap[i], pq.heap[j] = pq.heap[j], pq.heap[i]
	pq.heap[i].index = i
	pq.heap[j].index = j
}

// When the priority of a node i becomes greater than the priority of
// its parent i/2, the heap property is violated. A fix is to swap i
// and i/2 and to repeat the process, stopping when the heap property
// is restored. swim implements this procedure.
func (pq *PriorityQueue) swim(i int) {// bottom-up	
	for i > 1 && pq.greater(i, i / 2) {
		pq.swap(i, i / 2)
		i = i / 2
	}
}

// sink is used to restore the heap property violated when the
// priority of node i becomes smaller than one or both of its children
// 2*i or 2*i+1. Node i is swapped with the child that owns a greater
// priority. The procedure is repeated, stopping when the heap
// property is restored.
func (pq *PriorityQueue) sink(i int) { // top-down
	n := pq.size
	
	for 2*i <= n {
		j := 2*i
		if j < n && pq.greater(j + 1, j) {
			j++
		}

		if pq.greater(i, j) { // heap property is ok
			break;
		}
		pq.swap(i, j)
		i = j
	}
}

// Push inserts a node with priority prio at the end of heap. Heap
// property is mantained by calling swim.
func (pq *PriorityQueue) Push(prio int, key, value interface{}) { 
	pq.size++
	
	if pq.size + 1 >= cap(pq.heap) { // consider ignored 0th element
		newHeap := make([]*HeapNode, 2*cap(pq.heap))
		copy(newHeap, pq.heap)
		pq.heap = newHeap
	}

	pq.heap[pq.size] = &HeapNode{pq.size, prio, key, value}
	pq.keyMap[pq.heap[pq.size].key] = pq.heap[pq.size]
	
	pq.swim(pq.size)
}

// Pop returns the node located at the beginning of the heap. Heap
// property is mantained using sink.
func (pq *PriorityQueue) Pop() *HeapNode {
	ret := pq.heap[1]
	
	pq.heap[1] = pq.heap[pq.size]
	pq.heap[1].index = 1  // update heap index in heap node struct
	pq.heap[pq.size] = nil // help garbage collector
	pq.size--

	pq.sink(1)

	return ret
}

// Contains returns true if the key is in the heap.
func (pq *PriorityQueue) Contains(key interface{}) bool {
	_, ok := pq.keyMap[key]

	return ok
}

// When the priority of a node is changed, Update is used the mantain
// the heap property. Update verify the new priority and if it is
// greater than old priority, the new node priority would become
// greater than its parent, swim is used to maintain the heap
// property, otherwise sink is used.
func (pq *PriorityQueue) Update(prio int, key, value interface{}) {
	node, ok := pq.keyMap[key]

	if !ok {
		pq.Push(prio, key, value)
		return
	}

	oldPrio := node.prio
	node.prio = prio
	node.value = value

	if pq.greaterprio(node.prio, oldPrio) {
		pq.swim(node.index)
	} else {
		pq.sink(node.index)
	}
}
