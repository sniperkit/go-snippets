package LinkedStack

// A Node is the data stored inside of a Stack.
type Node struct {
	Data interface{}
	Next *Node
}

// A Stack is a first in, last out data structure.
type Stack struct {
	top    *Node
	length int
}

// Stack constructor which creates an empty Stack.
func new() *Stack {
	stack := new()
	stack.top = nil
	stack.length = 0
	return stack
}

// Empty empties the Stack.
func (stack *Stack) Empty() {
	stack.top = nil
	stack.length = 0
}

// IsEmpty returns true if the Stack is empty.
func (stack *Stack) IsEmpty() bool {
	if stack.length == 0 {
		return true
	}
	return false
}

// Len returns the length of the Stack.
func (stack *Stack) Len() int {
	return stack.length
}

// Peek returns the top Node.
func (stack *Stack) Peek() interface{} {
	return stack.top
}

// Push adds a Node to the top of the Stack.
func (stack *Stack) Push(Data interface{}) {
	if stack.length == 0 {
		stack.top = &Node{Data: Data, Next: nil}
	} else {
		stack.top = &Node{Data: Data, Next: stack.top}
	}
	stack.length++
}

// Pop removes the top Node in the stack and returns it.
func (stack *Stack) Pop() interface{} {
	temp := stack.top
	stack.top = stack.top.Next
	stack.length--
	return temp
}
