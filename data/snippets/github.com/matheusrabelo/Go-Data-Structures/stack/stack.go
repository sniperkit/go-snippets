package stack

type Node struct {
	value    interface{}
	previous *Node
}

type Stack struct {
	top  *Node
	size int
}

//Top returns the top of stack
func (s *Stack) Top() interface{} {
	if s.size == 0 {
		return nil
	}
	return s.top.value
}

//Push pushes given element to the top of the stack and returns the created Node
func (s *Stack) Push(value interface{}) *Node {
	n := new(Node)
	n.value = value
	n.previous = s.top
	s.top = n
	s.size = s.size + 1
	return n
}

//Pop removes the top of the stack and returns his value
func (s *Stack) Pop() interface{} {
	if s.size == 0 {
		return nil
	}
	value := s.top.value
	s.top = s.top.previous
	s.size = s.size - 1
	return value
}

//Size returns the size of the Stack
func (s *Stack) Size() int {
	return s.size
}

//Empty returns true if the stack is empty
func (s *Stack) Empty() bool {
	return s.size == 0
}

//New creates and returns a new Stack instance
func New() *Stack {
	s := new(Stack)
	s.size = 0
	s.top = nil
	return s
}
