package reactor

import (
	"errors"
	"log"
	"math/rand"
	"sort"
	"sync"

	tomb "gopkg.in/tomb.v1"
)

var (
	ErrorReactorCapacity error = errors.New("reactor: reach max capacity")
)

type Reactor interface {
}

type reactor struct {
	tomb  *tomb.Tomb
	mutex sync.Mutex

	capacity   int
	threshold  int
	concurrent int
	pending    []*nuclear
	sorted     bool

	notifyCh chan<- *reactor
	material chan struct{}
	closed   chan struct{}
}

func NewReactor(capacity, threshold, concurrent int, notifyCh chan<- *reactor) *reactor {
	r := &reactor{
		tomb:       new(tomb.Tomb),
		capacity:   capacity,
		threshold:  threshold,
		concurrent: concurrent,
		notifyCh:   notifyCh,
		material:   make(chan struct{}, 1),
	}
	// A closed channel is used to provide immediate route through a
	// select call in a loop function
	r.closed = make(chan struct{})
	close(r.closed)
	go func() {
		defer r.tomb.Done()
		r.tomb.Kill(r.loop())
	}()
	return r
}

func (r *reactor) Kill() {
	r.tomb.Kill(nil)
}

func (r *reactor) AddNuclear(n *nuclear) error {
	if n == nil {
		return nil
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.pending) >= r.capacity {
		return ErrorReactorCapacity
	}
	r.pending = append(r.pending, n)
	r.sorted = false
	if len(r.pending) == 1 {
		r.material <- struct{}{}
	}
	return nil
}

func (r *reactor) loop() error {
	log.Printf("reactor: loop run")

	var next chan struct{} = r.closed
	var out chan struct{}
	concurrentCh := make(chan struct{}, r.concurrent)

	for {
		select {
		case <-r.tomb.Dying():
			log.Printf("reactor: reactor shutdown")
			return nil
		case <-r.material:
			// nuclear put in
			out = concurrentCh
			continue
		case <-next:
			out = concurrentCh
			continue
		case out <- struct{}{}: // concurrent control. Block maybe
			out = nil
		}
		// notify nuclear product material if current nuclear less than threshold
		r.fillUpNuclear()
		// pop one nuclear by its priority
		nu, empty := r.popOne()
		if empty {
			next = nil
		} else {
			next = r.closed
		}
		if nu != nil {
			log.Printf("reactor: prepare to execute nuclear: %+v", nu)
			go func(n *nuclear) {
				defer func() {
					select {
					case <-concurrentCh:
					default:
					}
				}()
				n.Reaction()
			}(nu)
			continue
		}
		// put back into concurrent pool
		<-concurrentCh
	}
	return nil
}

func (r *reactor) fillUpNuclear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.pending) < r.threshold {
		select {
		case r.notifyCh <- r:
		default:
		}
		return
	}
	return
}

func (r *reactor) popOne() (*nuclear, bool) {
	log.Printf("reactor: pop one nuclear")
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.pending) == 0 {
		return nil, true
	}
	if !r.sorted {
		sort.Sort(ByNuclearPriority(r.pending))
		r.sorted = true
	}
	n := r.randomNuclearWithWeigth()
	return n, len(r.pending) == 0
}

func (r *reactor) randomNuclearWithWeigth() *nuclear {
	var sumweight int
	for _, n := range r.pending {
		sumweight += n.priority
	}
	randomValue := rand.Intn(sumweight)
	for index, n := range r.pending {
		if randomValue < n.priority {
			r.pending = append(r.pending[0:index], r.pending[index+1:]...)
			return n
		}
		randomValue -= n.priority
	}
	panic("reactor: should never get hear!")
}
