package pool

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
)

var (
	InitialCap = 5
	MaximumCap = 30
	network    = "tcp"
	address    = "127.0.0.1:7777"
	factory    = func() (Object, error) { return net.Dial(network, address) }
)

func init() {
	// used for factory function
	go simpleTCPServer()
	time.Sleep(time.Millisecond * 300) // wait until tcp server has been settled

	rand.Seed(time.Now().UTC().UnixNano())
}

func TestNew(t *testing.T) {
	_, err := newChannelPool()
	if err != nil {
		t.Errorf("New error: %s", err)
	}
}
func TestPool_Get_Impl(t *testing.T) {
	p, _ := newChannelPool()
	defer p.Close()

	_, err := p.Borrow()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
}

func TestPool_Get(t *testing.T) {
	p, _ := newChannelPool()
	defer p.Close()

	_, err := p.Borrow()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}

	// after one get, current capacity should be lowered by one.
	if p.Len() != (InitialCap - 1) {
		t.Errorf("Get error. Expecting %d, got %d",
			(InitialCap - 1), p.Len())
	}

	// get them all
	var wg sync.WaitGroup
	for i := 0; i < (InitialCap - 1); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.Borrow()
			if err != nil {
				t.Errorf("Get error: %s", err)
			}
		}()
	}
	wg.Wait()

	if p.Len() != 0 {
		t.Errorf("Get error. Expecting %d, got %d",
			(InitialCap - 1), p.Len())
	}

	_, err = p.Borrow()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
}

func TestPool_Put(t *testing.T) {
	p, err := NewChannelPool(0, 30, factory)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	// get/create from the pool
	objs := make([]PoolObjecter, MaximumCap)
	for i := 0; i < MaximumCap; i++ {
		obj, _ := p.Borrow()
		objs[i] = obj
	}

	// now put them all back
	for _, obj := range objs {
		obj.Return()
	}

	if p.Len() != MaximumCap {
		t.Errorf("Put error len. Expecting %d, got %d",
			1, p.Len())
	}

	obj, _ := p.Borrow()
	p.Close() // close pool

	obj.Return() // try to put into a full pool
	if p.Len() != 0 {
		t.Errorf("Put error. Closed pool shouldn't allow to put connections.")
	}
}

func TestPool_PutUnusableConn(t *testing.T) {
	p, _ := newChannelPool()
	defer p.Close()

	// ensure pool is not empty
	obj, _ := p.Borrow()
	obj.Return()

	poolSize := p.Len()
	obj, _ = p.Borrow()
	obj.Return()
	if p.Len() != poolSize {
		t.Errorf("Pool size is expected to be equal to initial size")
	}

	obj, _ = p.Borrow()
	if po, ok := obj.(*PoolObject); !ok {
		t.Errorf("impossible")
	} else {
		po.MarkUnusable()
	}
	obj.Return()
	if p.Len() != poolSize-1 {
		t.Errorf("Pool size is expected to be initial_size - 1", p.Len(), poolSize-1)
	}
}

func TestPool_UsedCapacity(t *testing.T) {
	p, _ := newChannelPool()
	defer p.Close()

	if p.Len() != InitialCap {
		t.Errorf("InitialCap error. Expecting %d, got %d",
			InitialCap, p.Len())
	}
}

func TestPool_Close(t *testing.T) {
	p, _ := newChannelPool()

	// now close it and test all cases we are expecting.
	p.Close()

	c := p.(*channelPool)

	if c.objects != nil {
		t.Errorf("Close error, objects channel should be nil")
	}

	if c.factory != nil {
		t.Errorf("Close error, factory should be nil")
	}

	_, err := p.Borrow()
	if err == nil {
		t.Errorf("Close error, get object should return an error")
	}

	if p.Len() != 0 {
		t.Errorf("Close error used capacity. Expecting 0, got %d", p.Len())
	}
}

func TestPoolConcurrent(t *testing.T) {
	p, _ := newChannelPool()
	pipe := make(chan PoolObjecter, 0)

	go func() {
		p.Close()
	}()

	for i := 0; i < MaximumCap; i++ {
		go func() {
			obj, _ := p.Borrow()

			pipe <- obj
		}()

		go func() {
			obj := <-pipe
			if obj == nil {
				return
			}
			obj.Return()
		}()
	}
}

func TestPoolWriteRead(t *testing.T) {
	p, _ := NewChannelPool(0, 30, factory)

	obj, _ := p.Borrow()

	obj_ := obj.(*PoolObject)

	conn := obj_.Object.(net.Conn)
	msg := "hello"
	_, err := conn.Write([]byte(msg))
	if err != nil {
		t.Error(err)
	}
}

func TestPoolConcurrent2(t *testing.T) {
	p, _ := NewChannelPool(0, 30, factory)

	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				obj, _ := p.Borrow()
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				obj.Return()
				wg.Done()
			}(i)
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			obj, _ := p.Borrow()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			obj.Return()
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestPoolTransaction(t *testing.T) {
	p, _ := newChannelPool()
	defer p.Close()

	// invalied func
	f := func() {}
	_, err := p.Transaction(f)
	if err == nil {
		t.Errorf("pool transaction must have function with one argument")
	}

	// argument type is wrong
	f1 := func(i int) {}
	_, err = p.Transaction(f1)
	if err == nil {
		t.Errorf("pool transaction function first argument is Must be Object")
	}

	// right one
	f2 := func(obj Object) {}
	_, err = p.Transaction(f2)
	if err != nil {
		t.Error(err)
	}

	if p.Len() != InitialCap {
		t.Errorf("InitialCap error. Expecting %d, got %d",
			InitialCap, p.Len())
	}

}

func newChannelPool() (Pool, error) {
	return NewChannelPool(InitialCap, MaximumCap, factory)
}

func simpleTCPServer() {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			buffer := make([]byte, 256)
			conn.Read(buffer)
		}()
	}
}
