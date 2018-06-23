// A binary tree of capped binary heaps
package tree

import heap "github.com/go-li/libgeneration/heap/binary"
import tree "github.com/go-li/libgeneration/dict/binary"

var Cap uint8 = 9

func Insert(h *Root, add [], cmp func(*, *) int) {

	var capacity = 1 << Cap

	h.length += len(add)

	var compareNode = func(a, b *Node) int {
		if cmp(&a.full_heap[0], &b.full_heap[0]) < 0 {
			return -1
		}
		return 1
	}

	for len(add) >= capacity-1 {

		var newlyadded = add[0 : capacity-1 : capacity-1]

		heap.Heapify(newlyadded, cmp)
		add = add[capacity-1:]

		if len(newlyadded) > 0 {
			var nod Node
			nod.full_heap = newlyadded

			tree.Append(heapTree, &h.full_heaps_tree, &nod, compareNode)

		}
	}

	var workcap int

	if len(h.working_heap) > 0 {
		workcap = len(h.working_heap)
	}

	var exceed = workcap + len(add) - (capacity - 1)

	if exceed >= 0 {

		heap.Insert(&h.working_heap, add[exceed:], cmp)

		if exceed > 0 {

			var nod Node
			nod.full_heap = h.working_heap

			tree.Append(heapTree, &h.full_heaps_tree, &nod, compareNode)

			h.working_heap = add[0:exceed:exceed]

			heap.Heapify(h.working_heap, cmp)
		}

	} else if len(h.working_heap) > 0 {
		heap.Insert(&h.working_heap, add, cmp)

	} else {

		h.working_heap = add
		heap.Heapify(add, cmp)
	}

}

func Fix(h *Root, cmp func(*, *) int) {

	var compareNode = func(a, b *Node) int {
		if cmp(&a.full_heap[0], &b.full_heap[0]) < 0 {
			return -1
		}
		return 1
	}

	var minnode *Node = tree.Find(heapTree, h.full_heaps_tree, nil)

	if minnode != nil {
		heap.Fix(minnode.full_heap, 0, cmp)

		tree.Delete(heapTree, &h.full_heaps_tree, nil, func(_, b *Node) int {
			if b == minnode {
				return 0
			}
			return -1
		})
		tree.Append(heapTree, &h.full_heaps_tree, minnode, compareNode)

	}
	if len(h.working_heap) > 0 {
		heap.Fix(h.working_heap, 0, cmp)
	}

}

func Pop(h *Root, cmp func(*, *) int) (ret *) {

	var minnode *Node = tree.Find(heapTree, h.full_heaps_tree, nil)

	if (minnode == nil) && (len(h.working_heap) == 0) {
		return nil
	}

	if (minnode != nil) && (cmp(&minnode.full_heap[0], &h.working_heap[0]) < 0) {

		var tmp []
		tmp = make([], 1)
		var temp *
		var p *
		var q *
		temp = &tmp[0]

		p = &minnode.full_heap[0]
		q = &h.working_heap[0]

		*temp = *p
		*p = *q
		*q = *temp

		Fix(h, cmp)

		minnode = tree.Find(heapTree, h.full_heaps_tree, nil)
	}

	ret = heap.Pop(&h.working_heap, cmp)
	h.length--

	if len(h.working_heap) == 0 {
		if minnode == nil {
			return ret
		}

		h.working_heap = minnode.full_heap
		tree.Delete(heapTree, &h.full_heaps_tree, nil, func(_, b *Node) int {
			if b == minnode {
				return 0
			}
			return -1
		})

	}

	return ret
}

func Peek(h *Root, cmp func(*, *) int) (ret *) {

	var minnode *Node = tree.Find(heapTree, h.full_heaps_tree, nil)
	if (minnode != nil) && (len(h.working_heap) > 0) {
		var res = cmp(&minnode.full_heap[0], &h.working_heap[0])
		if res < 0 {
			ret = &minnode.full_heap[0]
			return
		}
		ret = &h.working_heap[0]
		return
	}

	if minnode != nil {
		ret = &minnode.full_heap[0]
	} else if len(h.working_heap) > 0 {
		ret = &h.working_heap[0]
	}

	return
}

func Heapify(h *Root, cmp func(*, *) int) {

	heap.Heapify(h.working_heap, cmp)
	tree.Foreach(heapTree, h.full_heaps_tree, func(node *Node) bool {
		heap.Heapify(node.full_heap, cmp)
		return true
	})

}

func Is(h *Root, cmp func(*, *) int) (o bool) {

	o = true
	o = o && heap.Is(h.working_heap, cmp)
	tree.Foreach(heapTree, h.full_heaps_tree, func(node *Node) bool {
		o = o && heap.Is(node.full_heap, cmp)
		return true
	})

	return o
}

func heapTree(n *Node) *[2]*Node {
	return &n.order_link
}

type Root struct {
	full_heaps_tree *Node
	working_heap    []
	length          int
}

type Node struct {
	order_link [2]*Node
	full_heap  []
}
