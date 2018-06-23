// 使用堆接口构建优先级队列
package main

import (
	"container/heap"
	"fmt"
)

// Item为优先级队列中管理的对象
type Item struct {
	index int // heap中存储对象的索引
	value    string // 对象存储的信息，这里为string，可自己进行修改
	priority int    // 对象在队列中的优先级
}


// PriorityQueue type 的实例使用堆， 接口来保存Item
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
	}

func (pq PriorityQueue) Less(i, j int) bool {
	// 根据priority判断Queue中对象的优先级
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	// 修改对象的索引
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// 更新对queue中对象权限、数值的修改
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}


func main() {

	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// 插入对象并修改其优先级
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value, 5)

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s \n", item.priority, item.value)
	}
}
