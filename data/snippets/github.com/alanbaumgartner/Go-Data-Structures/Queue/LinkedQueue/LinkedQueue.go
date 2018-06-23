package LinkedQueue

// A Node is the data stored inside of a Queue.
type Node struct {
	Data interface{}
	Next *Node
}

// A Queue is a first in, first out data structure.
type Queue struct {
	front  *Node
	length int
}

// Queue constructor which creates an empty Queue.
func new() *Queue {
	queue := new()
	queue.front = nil
	queue.length = 0
	return queue
}

// Empty empties the Queue.
func (queue *Queue) Empty() {
	queue.front = nil
	queue.length = 0
}

// IsEmpty returns true if the Queue is empty and false if not.
func (queue *Queue) IsEmpty() bool {
	if queue.length == 0 {
		return true
	}
	return false
}

// Len returns the length of the Queue.
func (queue *Queue) Len() int {
	return queue.length
}

// Peek returns next Node to be dequeued.
func (queue *Queue) Peek() *Node {
	return queue.front
}

// Enqueue adds a Node to the end of the Queue.
func (queue *Queue) Enqueue(Data interface{}) {
	if queue.length == 0 {
		queue.front = &Node{Data: Data, Next: nil}
	} else {
		temp := queue.front
		for temp.Next != nil {
			temp = temp.Next
		}
		temp.Next = &Node{Data: Data, Next: nil}
	}
	queue.length++
}

// Dequeue removes the first Node in the Queue and returns it.
func (queue *Queue) Dequeue() *Node {
	temp := queue.front
	queue.front = queue.front.Next
	queue.length--
	return temp
}
