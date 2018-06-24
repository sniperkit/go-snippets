/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/11/08        Yang Chenglong
 */

package main

import (
	"fmt"
	"io/ioutil"
	"sort"
)

func main() {
	dat, err := ioutil.ReadFile("./test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	a := GenerateHuffmanCodes(string(dat))
	for i, _ := range a {
		fmt.Println(i, ":", a[i])
	}
}

type Node struct {
	value  string
	weigth int
	left   *Node
	right  *Node
	index  int
}

type PriorityQueue struct {
	Queue []*Node
}

func NewPriorityQueue(size int) *PriorityQueue {
	p := new(PriorityQueue)
	return p
}

func (pq *PriorityQueue) Len() int {
	return len(pq.Queue)
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.Queue[i], pq.Queue[j] = pq.Queue[j], pq.Queue[i]
	pq.Queue[i].index = i
	pq.Queue[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := pq.Len()
	node := x.(*Node)
	node.index = n
	pq.Queue = append(pq.Queue, node)
	sort.Sort(pq)
}

func (pq *PriorityQueue) Pop() interface{} {
	min := pq.Queue[0]
	if pq.Len() == 1 {
		pq.Queue = make([]*Node, 0)
		return min
	}
	pq.Queue = pq.Queue[1:pq.Len()]
	return min
}

func (q *PriorityQueue) Less(i, j int) bool {
	return q.Queue[i].weigth < q.Queue[j].weigth
}

func Counter(s string) map[string]int {
	c := make(map[string]int)
	for i := 0; i < len(s); i++ {
		c[string(s[i])] = c[string(s[i])] + 1
	}
	return c
}

func BuildHuffmanTree(nodes map[string]int) *Node {
	var root *Node
	size := len(nodes)
	queue := NewPriorityQueue(size)
	for value, weight := range nodes {
		node := &Node{
			value:  value,
			weigth: weight,
		}
		queue.Push(node)
	}

	for queue.Len() > 2 {
		leftChild := queue.Pop().(*Node)
		rightChild := queue.Pop().(*Node)
		root = &Node{
			weigth: leftChild.weigth + rightChild.weigth,
			left:   leftChild,
			right:  rightChild,
		}

		queue.Push(root)
	}

	leftChild := queue.Pop().(*Node)
	rightChild := queue.Pop().(*Node)
	root = &Node{
		weigth: leftChild.weigth + rightChild.weigth,
		left:   leftChild,
		right:  rightChild,
	}

	queue.Push(root)

	return queue.Pop().(*Node)
}

func huffmanCodes(node *Node, prefix string, codes map[string]string) {
	if node.value != "" {
		codes[node.value] = prefix
	} else {
		leftPrefix := prefix + "0"
		huffmanCodes(node.left, leftPrefix, codes)
		rightPrefix := prefix + "1"
		huffmanCodes(node.right, rightPrefix, codes)
	}
}

func GenerateHuffmanCodes(s string) map[string]string {
	c := Counter(s)
	fmt.Println(c)
	tree := BuildHuffmanTree(c)
	codes := make(map[string]string)
	huffmanCodes(tree, "", codes)
	return codes
}
