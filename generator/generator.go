package generator

// Generator represents a generator structure
type Generator struct {
	bound int64
	fn    func(int64) interface{}
	out   chan interface{}
}

//New creates new unbounded generator
func New(fn func(int64) interface{}) *Generator {
	g := &Generator{
		fn:  fn,
		out: make(chan interface{}),
	}
	go func() {
		for i := int64(0); ; i++ {
			g.out <- g.fn(i)
		}
	}()
	return g
}

//NewBounded creates new bounded generator, which will generate no more than `bound` values
func NewBounded(bound int64, fn func(int64) interface{}) *Generator {
	g := &Generator{
		bound: bound,
		fn:    fn,
		out:   make(chan interface{}, bound),
	}
	go func() {
		for i := int64(0); i < bound; i++ {
			g.out <- g.fn(i)
		}
		g.Stop()
	}()
	return g
}

// Next will return next generated value
func (g *Generator) Next() interface{} {
	return <-g.out
}

// Stop will stop generator
func (g *Generator) Stop() {
	close(g.out)
}
