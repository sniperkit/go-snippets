package leakybuff

const (
	defaultBufSize  int = 1024
	defaultCapacity int = 1024
)

var defaultLeakyBuff *LeakyBuff = NewLeakyBuff(defaultBufSize, defaultCapacity)

// thread safe

type LeakyBuff struct {
	bufSize  int
	bucketCh chan []byte
}

func NewLeakyBuff(bufSize, cap int) *LeakyBuff {
	lb := &LeakyBuff{
		bufSize:  bufSize,
		bucketCh: make(chan []byte, cap),
	}
	return lb
}

func (lb *LeakyBuff) Get() []byte {
	var buf []byte
	select {
	case buf = <-lb.bucketCh:
	default:
		buf = make([]byte, lb.bufSize)
	}
	return buf
}

// If we put into a same buffer more than one times.
// That is to say, more than one can hold it and operate it,
// maybe same time.
func (lb *LeakyBuff) Put(buf []byte) {
	if len(buf) != lb.bufSize {
		panic("invalid buff size that's put into leaky buffer")
	}

	// FIXME: reset this buffer ?
	//buf = buf[:0]

	select {
	case lb.bucketCh <- buf:
	default:
	}
	return
}
