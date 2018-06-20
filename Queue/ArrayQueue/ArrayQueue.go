package ArrayQueue

// A Queue is a first in, first out data structure.
type Queue struct {
	array  []interface{}
	front  interface{}
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
func (queue *Queue) Peek() interface{} {
	return queue.front
}

// Enqueue adds a Node to the end of the Queue.
func (queue *Queue) Enqueue(Data interface{}) {
	if queue.length == 0 {
		queue.front = Data
	}
	queue.array = append(queue.array, Data)
	queue.length++
}

// Dequeue removes the first Node in the Queue and returns it.
func (queue *Queue) Dequeue() interface{} {
	temp := queue.front
	queue.array = queue.array[1:]
	queue.length--
	queue.front = queue.array[0]
	return temp
}
