package act

type Stack struct {
	in   chan interface{}
	out  chan interface{}
	done chan struct{}
}

func NewStack() *Stack {
	s := &Stack{
		in:   make(chan interface{}),
		out:  make(chan interface{}),
		done: make(chan struct{}),
	}
	go s.run()
	return s
}

func (s *Stack) Push() chan<- interface{} {
	return s.in
}

func (s *Stack) Pop() <-chan interface{} {
	return s.out
}

func (s *Stack) Close() {
	s.done <- struct{}{}
}

func (s *Stack) run() {
	var last interface{}
	var buf []interface{}

	last = <-s.in

	for {
		select {
		case value := <-s.in:
			buf = append(buf, last)
			last = value

		case s.out <- last:
			sz := len(buf)
			if sz == 0 {
				last = <-s.in
			} else {
				last, buf = buf[sz-1], buf[:sz-1]
			}

		case <-s.done:
			return
		}
	}
}
