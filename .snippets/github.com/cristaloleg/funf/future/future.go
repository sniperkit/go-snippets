package future

type Future struct {
	out   chan interface{}
	value interface{}
	done  bool
}

func New(fn func() interface{}) *Future {
	f := Future{
		out: make(chan interface{}, 1),
	}
	go func() {
		f.out <- fn()
	}()
	return &f
}

func (f *Future) IsClosed() bool {
	return f.done
}

func (f *Future) Value() interface{} {
	value, ok := <-f.out
	if ok {
		f.value = value
		f.done = true
		close(f.out)
	}
	return f.value
}
