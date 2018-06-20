package lib

import (
	"fmt"
	"sync"
	"time"
)

const (
	idleThreshold     = 300 //seconds
	statusServing     = "serving"
	statusIdle        = "Idle"
	statusDestroyed   = "Destroyed"
	maxLenOfDestroyed = 100
)

//Runtime ...
type Runtime struct {
	ID         string      `json:"id"`
	Target     ProxyTarget `json:"target"`
	ActiveTime int64       `json:"active_time"`
	Status     string      `json:"status"`
	Image      string      `json:"container_image"`
}

//RuntimePool ...
type RuntimePool struct {
	pool          map[string]*Runtime
	destroyedOnes []*Runtime
	destroyedPtr  uint16
	lock          *sync.RWMutex
}

//NewRuntimePool ...
func NewRuntimePool() *RuntimePool {
	return &RuntimePool{
		pool:          make(map[string]*Runtime),
		destroyedOnes: make([]*Runtime, 0, maxLenOfDestroyed),
		destroyedPtr:  0,
		lock:          new(sync.RWMutex),
	}
}

//Index ...
func (rp *RuntimePool) Index(key string) (*Runtime, bool) {
	rp.lock.RLock()
	defer rp.lock.RUnlock()

	r, ok := rp.pool[key]

	return r, ok
}

//Put ...
func (rp *RuntimePool) Put(key string, r *Runtime) error {
	rp.lock.Lock()
	defer rp.lock.Unlock()

	if _, ok := rp.pool[key]; ok {
		return fmt.Errorf("%s existing", key)
	}

	//Set status
	r.Status = statusServing
	rp.pool[key] = r

	return nil
}

//Use ...
func (rp *RuntimePool) Use(key string) (*Runtime, error) {
	rp.lock.Lock()
	defer rp.lock.Unlock()

	r, ok := rp.pool[key]
	if ok {
		//Update active time
		r.ActiveTime = time.Now().Unix()
		r.Status = statusServing
		return r, nil
	}

	return r, fmt.Errorf("%s not existing", key)
}

//Remove ...
//NOT used yet
func (rp *RuntimePool) Remove(prefix, ID string) error {
	rp.lock.Lock()
	defer rp.lock.Unlock()

	if _, ok := rp.pool[giveMeKey(prefix, ID)]; ok {
		delete(rp.pool, giveMeKey(prefix, ID))
	}

	return nil
}

//Garbages ...
func (rp *RuntimePool) Garbages() []*Runtime {
	rp.lock.Lock()
	defer rp.lock.Unlock()

	var garbages []*Runtime
	now := time.Now().Unix()
	for k, v := range rp.pool {
		if v.Status == statusIdle && now >= v.ActiveTime+idleThreshold {
			//Garbage
			garbages = append(garbages, v)
			delete(rp.pool, k) //removed from pool

			v.Status = statusDestroyed
			//Keep in the destroyed list for a while
			if len(rp.destroyedOnes) < maxLenOfDestroyed {
				rp.destroyedOnes = append(rp.destroyedOnes, v)
			} else {
				//override
				if rp.destroyedPtr < maxLenOfDestroyed {
					rp.destroyedOnes[rp.destroyedPtr] = v
				} else {
					//back to the start
					rp.destroyedPtr = 0
					rp.destroyedOnes[rp.destroyedPtr] = v
				}

				rp.destroyedPtr++
			}
		}
	}

	return garbages
}

//GetAll runtimes as a list
func (rp *RuntimePool) GetAll() []*Runtime {
	rp.lock.RLock()
	defer rp.lock.RUnlock()

	list := make([]*Runtime, 0)
	for _, v := range rp.pool {
		list = append(list, v)
	}

	for _, v := range rp.destroyedOnes {
		list = append(list, v)
	}

	return list
}

//SetIdle ...
func (rp *RuntimePool) SetIdle(key string) {
	rp.lock.Lock()
	defer rp.lock.Unlock()

	if runtime, ok := rp.pool[key]; ok {
		runtime.Status = statusIdle
	}
}

func giveMeKey(prefix, ID string) string {
	return fmt.Sprintf("%s:%s", prefix, ID)
}
