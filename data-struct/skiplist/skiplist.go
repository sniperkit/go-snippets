package skiplist

import (
	"fmt"
	"math/rand"
	"time"
)

// This implemention is not thread safe.

const (
	MAX_LEVEL int = 10
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type Node struct {
	key     int
	value   int
	forward []*Node // The forward node in every level it comes. This is the key for understanding the algorithm.
}

func NewNode(level, key, value int) *Node {
	n := new(Node)
	n.key = key
	n.value = value
	n.forward = make([]*Node, level)
	return n
}

func (n *Node) Key() int         { return n.key }
func (n *Node) Value() int       { return n.value }
func (n *Node) Forward() []*Node { return n.forward }

type SkipList struct {
	level  int
	header *Node
}

func NewSkipList() *SkipList {
	sl := new(SkipList)
	sl.level = 0
	sl.header = NewNode(MAX_LEVEL, 0, 0)
	return sl
}

func (sl *SkipList) randomLevel() int {
	var k = 1
	for random.Int()%2 != 0 {
		k++
	}
	if k < MAX_LEVEL {
		return k
	}
	return MAX_LEVEL
}

func (sl *SkipList) Insert(key, value int) bool {
	return sl.insert(key, value) != nil
}

func (sl *SkipList) insert(key, value int) *Node {
	var updates = make([]*Node, MAX_LEVEL)
	var p *Node
	var q *Node

	p = sl.header
	k := sl.level

	// find update nodes
	for i := k - 1; i >= 0; i-- {
		for q = p.forward[i]; q != nil && q.key < key; {
			p = q
			q = p.forward[i]
		}
		updates[i] = p
	}

	// can not insert same key
	if q != nil && q.key == key {
		return nil
	}

	// random a level to insert
	k = sl.randomLevel()
	if k > sl.level {
		// generate a new level
		for i := sl.level; i < k; i++ {
			updates[i] = sl.header
		}
		sl.level = k
	}

	// initialize forward nodes of every level. And
	// link it behind prevous node.
	node := NewNode(k, key, value)
	for i := 0; i < k; i++ {
		node.forward[i] = updates[i].forward[i]
		updates[i].forward[i] = node
	}

	return node
}

func (sl *SkipList) Delete(key int) bool {
	return sl.delete(key) != nil
}

func (sl *SkipList) delete(key int) *Node {
	var p *Node
	var q *Node

	p = sl.header
	k := sl.level

	// updates
	var updates = make([]*Node, MAX_LEVEL)
	for i := k - 1; i >= 0; i-- {
		for q = p.forward[i]; q != nil && q.key < key; {
			p = q
			q = p.forward[i]
		}
		updates[i] = p
	}

	if q != nil && q.key == key {
		for i := 0; i < k; i++ {
			if updates[i].forward[i] == q {
				updates[i].forward[i] = q.forward[i]
			}
			for i := k - 1; i >= 0; i-- {
				if sl.header.forward[i] == nil {
					sl.level--
				}
			}
		}
		return q
	}

	return nil
}

func (sl *SkipList) Search(key int) (int, bool) {
	if node, ok := sl.search(key); ok {
		return node.value, true
	}
	return -1, false
}

func (sl *SkipList) search(key int) (*Node, bool) {
	var p *Node

	p = sl.header
	k := sl.level
	for i := k - 1; i >= 0; i-- {
		// if we skip to below level, we just begin from the node of above node 'p'
		for q := p.forward[i]; q != nil && q.key <= key; {
			if q.key == key {
				return q, true
			}
			p = q
			q = p.forward[i]
		}
	}

	return nil, false
}

func (sl *SkipList) Display() {
	k := sl.level
	for i := k - 1; i >= 0; i-- {
		fmt.Printf("Level %-2d : ", i+1)
		p := sl.header
		for p = p.forward[i]; p != nil; {
			fmt.Printf(" %-2d ->", p.value)
			p = p.forward[i]
		}
		fmt.Printf(" NULL \n")
	}
	fmt.Printf("\n")
}
