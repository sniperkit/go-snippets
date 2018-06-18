package queue

type Node struct {
	value interface{}
	next  *Node
}

type Queue struct {
	first *Node
	last  *Node
	size  int
}

//First returns the front of the queue
func (q *Queue) First() interface{} {
	if q.size == 0 {
		return nil
	}
	return q.first.value
}

//Push pushes given element to the front of the queue and returns the created Node
func (q *Queue) Push(value interface{}) *Node {
	n := new(Node)
	n.value = value
	n.next = nil
	if q.size == 0 {
		q.first = n
	} else {
		q.last.next = n
	}
	q.last = n
	q.size = q.size + 1
	return n
}

//Pop removes the front of the queue and returns his value
func (q *Queue) Pop() interface{} {
	if q.size == 0 {
		return nil
	}
	value := q.first.value
	q.first = q.first.next
	q.size = q.size - 1
	return value
}

//Size returns the size of the queue
func (q *Queue) Size() int {
	return q.size
}

//Empty returns true if the queue is empty
func (q *Queue) Empty() bool {
	return q.size == 0
}

//New creates and returns a new Queue instance
func New() *Queue {
	q := new(Queue)
	q.first = nil
	q.last = nil
	q.size = 0
	return q
}
