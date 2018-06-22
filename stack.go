package stack

type Stack struct {
	stack []interface{}
}

func New() *Stack {
	return &Stack{}
}

func (s *Stack) Push(v interface{}) {
	s.stack = append(s.stack, v)
}

func (s *Stack) Pop() interface{} {
	if len(s.stack) == 0 {
		return nil
	} else {
		v := s.stack[len(s.stack)-1]
		s.stack = s.stack[0 : len(s.stack)-1]
		return v
	}
}

func (s *Stack) Top() interface{} {
	if len(s.stack) == 0 {
		return nil
	} else {
		return s.stack[len(s.stack)-1]
	}
}
