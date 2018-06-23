package act

type Queue struct {
	in   chan interface{}
	out  chan interface{}
	done chan struct{}
}

func NewQueue() *Queue {
	q := &Queue{
		in:   make(chan interface{}),
		out:  make(chan interface{}),
		done: make(chan struct{}),
	}
	go q.run()
	return q
}

func (q *Queue) Push() chan<- interface{} {
	return q.in
}

func (q *Queue) Pop() <-chan interface{} {
	return q.out
}

func (q *Queue) Close() {
	q.done <- struct{}{}
}

func (q *Queue) run() {
	var last interface{}
	var buf []interface{}

	last = <-q.in

	for {
		select {
		case value := <-q.in:
			buf = append(buf, value)

		case q.out <- last:
			sz := len(buf)
			if sz == 0 {
				last = <-q.in
			} else {
				last, buf = buf[0], buf[1:]
			}

		case <-q.done:
			return
		}
	}
}
