/*
 * MIT License
 *
 * Copyright (c) 2017 Yang Chenglong
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

type Node interface {
	/*
		return the value < i
	*/
	Comparable(i Node) bool
}

type priorityqueue struct {
	pq []Node
}

func NewPriorityQueue() *priorityqueue {
	return &priorityqueue{pq: make([]Node, 1, 10)}
}

func (q *priorityqueue) IsEmpty() bool {
	return q.len() == 0
}

func (q *priorityqueue) Size() int {
	return q.len()
}

func (q *priorityqueue) Insert(v Node) {
	q.pq = append(q.pq, v)
	q.swim(q.len())
}

func (q *priorityqueue) DelMax() Node {
	max := q.pq[1]
	q.exchange(1, q.len())
	q.pq = q.pq[:q.len()]
	q.sink(1)
	return max
}

func (q *priorityqueue) less(i, j int) bool {
	return q.pq[i].Comparable(q.pq[j])
}

func (q *priorityqueue) exchange(i, j int) {
	q.pq[i], q.pq[j] = q.pq[j], q.pq[i]
}

func (q *priorityqueue) swim(k int) {
	for k > 1 && q.less(k/2, k) {
		q.exchange(k/2, k)
		k = k / 2
	}
}

func (q *priorityqueue) sink(k int) {
	for k*2 <= q.len() {
		j := k * 2
		if j < q.len() && q.less(j, j+1) {
			j++
		}
		if !q.less(k, j) {
			break
		}
		q.exchange(k, j)
		k = j
	}
}

func (q *priorityqueue) len() int {
	return len(q.pq) - 1
}
