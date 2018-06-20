package ArrayStack

// A Stack is a first in, last out data structure.
type Stack struct {
	array  []interface{}
	top    interface{}
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
	stack.top = Data
	stack.array = append(stack.array, Data)
	stack.length++
}

// Pop removes the top Node in the stack and returns it.
func (stack *Stack) Pop() interface{} {
	temp := stack.top
	stack.array = stack.array[:stack.length-1]
	stack.length--
	stack.top = stack.array[stack.length-1]
	return temp
}
