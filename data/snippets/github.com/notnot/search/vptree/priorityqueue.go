// priorityqueue.go, jpad 2016

package vptree

type _Item struct {
	p Point
	d float32
}

type priorityQueue []*_Item

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].d > pq[j].d // max-heap behavior
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(i interface{}) {
	item := i.(*_Item)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (pq priorityQueue) Top() interface{} {
	return pq[0]
}
