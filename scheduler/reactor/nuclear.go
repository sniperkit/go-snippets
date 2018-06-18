package reactor

import (
	"log"
	"sync"
	"time"

	tomb "gopkg.in/tomb.v1"
)

type NuclearGenerator interface {
	Next() ([]*nuclear, error)
}

type NuclearProductor struct {
	tomb  *tomb.Tomb
	mutex sync.Mutex

	reactor   *reactor
	generator NuclearGenerator

	fillUpRun bool
	notifyCh  chan *reactor
}

func NewNuclearProductor(g NuclearGenerator) *NuclearProductor {
	np := &NuclearProductor{
		tomb:      new(tomb.Tomb),
		generator: g,
		notifyCh:  make(chan *reactor, 1),
	}
	go func() {
		defer np.tomb.Done()
		np.tomb.Kill(np.loop())
	}()
	return np
}

func (n *NuclearProductor) loop() error {
	log.Printf("nuclear: nuclear productor run")

	var next chan struct{}
	closed := make(chan struct{})
	close(closed)

	for {
		select {
		case <-n.tomb.Dying():
			log.Printf("nuclear: nuclear productor shutdown")
			return nil
		case reactor := <-n.notifyCh:
			if n.reactor == nil {
				n.reactor = reactor
			}
			next = closed
		case <-next:
			// product nuclear for reactor
			n.productNuclear()
			next = nil
		}
	}
}

func (n *NuclearProductor) productNuclear() {
	n.mutex.Lock()
	if !n.fillUpRun {
		n.fillUpRun = true
		n.mutex.Unlock()
		// generate go
		go func() {
			defer func() {
				n.mutex.Lock()
				n.fillUpRun = false
				n.mutex.Unlock()
				log.Printf("nuclear: product over")
			}()
			// do it right now
			log.Printf("nuclear: product run")
			for {
				select {
				case <-n.tomb.Dying():
					return
				default:
				}
				nuclears, err := n.generator.Next()
				if err != nil {
					log.Printf("nuclear: nuclear productor generate nuclear error: %s", err)
					return
				}
				if len(nuclears) != 0 {
					for _, nuclear := range nuclears {
						err = n.reactor.AddNuclear(nuclear)
						switch err {
						case ErrorReactorCapacity:
							log.Printf("nuclear: add nuclear error")
							return
						default:
						}
					}
				}
			}
		}()

		return
	}
	n.mutex.Unlock()
	return
}

// nuclear struct
type nuclear struct {
	id       string
	priority int
}

func NewNuclear(id string, priority int) *nuclear {
	return &nuclear{
		id:       id,
		priority: priority,
	}
}

func (n *nuclear) Reaction() {
	log.Printf("nuclear: %s reaction over: %s", n.id, time.Now())
	return
}

type ByNuclearPriority []*nuclear

func (b ByNuclearPriority) Len() int           { return len(b) }
func (b ByNuclearPriority) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByNuclearPriority) Less(i, j int) bool { return b[i].priority > b[j].priority }
