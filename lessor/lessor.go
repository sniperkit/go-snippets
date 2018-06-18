package lessor

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cespare/xxhash"
	tomb "gopkg.in/tomb.v1"
)

var (
	ErrorLessorGrantAlreadyExists error = errors.New("lessor: lease already exists")
	ErrorLessorLeaseNotFound      error = errors.New("lessor: lease not found")
)

type LeaseID int64

type Lessor interface {
	// Grant grant a lease for a lease item
	Grant(item *LeaseItem) (*Lease, error)
	// Revoke revoke a lease with given id.
	Revoke(item *LeaseItem) error
	// Renew renews a lease with given id
	Renew(item *LeaseItem) error
	// Lookup
	Lookup(item *LeaseItem) (*Lease, error)
	// Stop stop the lessor
	Stop() error
}

// the lessor implement the Leasor interfaces
type lessor struct {
	tomb  *tomb.Tomb
	mutex sync.Mutex

	leaseMap  map[LeaseID]*Lease
	heartbeat time.Duration
	leaseTTL  time.Duration
}

func NewLeasor(heartbeat time.Duration, ttl time.Duration) Lessor {
	l := &lessor{
		tomb:      new(tomb.Tomb),
		leaseMap:  make(map[LeaseID]*Lease),
		heartbeat: heartbeat,
		leaseTTL:  ttl,
	}
	go func() {
		defer l.tomb.Done()
		l.tomb.Kill(l.runLoop())
	}()
	return l
}

func (l *lessor) Lookup(item *LeaseItem) (*Lease, error) {
	log.Printf("lessor: lookup item: %s", item)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	id := HashId(item)
	if lease, ok := l.leaseMap[id]; ok {
		return lease, nil
	}
	return nil, ErrorLessorLeaseNotFound
}

func (l *lessor) Grant(item *LeaseItem) (*Lease, error) {
	log.Printf("lessor: grant item: %s now: %s", item, time.Now())
	l.mutex.Lock()
	defer l.mutex.Unlock()
	id := HashId(item)
	if _, ok := l.leaseMap[id]; ok {
		return nil, ErrorLessorGrantAlreadyExists
	}
	lease := &Lease{
		Id:   id,
		item: item,
		ttl:  l.leaseTTL,
	}
	l.leaseMap[id] = lease
	lease.refresh(0)
	return lease, nil
}

func (l *lessor) Revoke(item *LeaseItem) error {
	log.Printf("lessor: revoke item: %s", item)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	id := HashId(item)
	lease := l.leaseMap[id]
	if lease == nil {
		log.Printf("lessor: remove item: %s not found", item)
		return ErrorLessorLeaseNotFound
	}
	//FIXME: revoke this item from meta
	delete(l.leaseMap, id)
	return nil
}

func (l *lessor) Renew(item *LeaseItem) error {
	log.Printf("lessor: renew item: %s", item)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	id := HashId(item)
	lease := l.leaseMap[id]
	if lease == nil {
		//FIXME: grant this item from meta, because cuckoo maybe down from the last refresh time.

		lease := &Lease{
			Id:   id,
			item: item,
			ttl:  l.leaseTTL,
		}
		l.leaseMap[id] = lease
		lease.refresh(0)
		return nil
	}
	lease.refresh(0)
	return nil
}

func (l *lessor) Stop() error {
	l.tomb.Kill(nil)
	return l.tomb.Wait()
}

func (l *lessor) ejectAllExpiredLeases() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for id, lease := range l.leaseMap {
		if lease.expired() {
			log.Printf("lessor: lease: %s expired now: %s", lease.item, time.Now())
			delete(l.leaseMap, id)
		}
	}
}

func (l *lessor) refreshAllItems() {
	log.Printf("lessor: refresh all items now: %s", time.Now())

	l.mutex.Lock()
	defer l.mutex.Unlock()

	for _, lease := range l.leaseMap {
		log.Printf("lessor: refresh %s", lease.item)
	}
}

func (l *lessor) runLoop() error {
	log.Printf("lessor: loop run")

	for {
		// eject all expired leases
		l.ejectAllExpiredLeases()

		//FIXME: renew leases from meta
		l.refreshAllItems()

		select {
		case <-l.tomb.Dying():
			log.Printf("lessor: shutdown")
			return nil
		case <-time.After(l.heartbeat):
		}
	}
}

// Lease
type Lease struct {
	Id   LeaseID
	item *LeaseItem
	ttl  time.Duration

	expiry time.Time
}

func (l *Lease) expired() bool {
	return l.Remaining() <= 0
}

// refresh refreshes the expiry of the lease.
func (l *Lease) refresh(extend time.Duration) {
	l.expiry = time.Now().Add(extend + l.ttl)
}

// Remaining returns the remaining time of the lease.
func (l *Lease) Remaining() time.Duration {
	return l.expiry.Sub(time.Now())
}

type LeaseItem struct {
	Client     string
	Host       string
	MountPoint string
}

func (l LeaseItem) String() string {
	return fmt.Sprintf("%s:%s:%s", l.Client, l.Host, l.MountPoint)
}

//
func HashId(item *LeaseItem) LeaseID {
	key := item.String()
	h := xxhash.New()
	h.Write([]byte(key))
	return LeaseID(h.Sum64())
}
