package sentinel

import (
	"fmt"
	"runtime"
	"sync/atomic"

	"github.com/go-mangos/mangos"
)

const (
	workPoolCap = 1024
)

type (
	// workPool is used as description worker pool
	workPool struct {
		id        int
		workCount int32
	}
)

// createWorkPool is used as
func createWorkPool(id int, leader *leader) {
	wp := &workPool{
		id:        id,
		workCount: 0,
	}
	wp.Run(leader)
}

func (wp *workPool) Run(leader *leader) {
	var message *mangos.Message
	var err error
	for {
		message, err = leader.pull.RecvMsg()
		if err != nil {

			if err == mangos.ErrClosed {
				fmt.Println("happen error , because of %s", err.Error())
				goto end
			} else if err == mangos.ErrRecvTimeout {

			} else {
				fmt.Println("happen error , because of %s", err.Error())
			}

			runtime.Gosched()
			continue
		}
		wp.incWorker(leader, message.Body)
		runtime.Gosched()
	}
end:
	for {
		if wp.isCanExit() {
			break
		}
		runtime.Gosched()
	}
	leader.waitGroup.Done()
}

// isCanExit is used as description workpool which is can exit
func (wp *workPool) isCanExit() bool {
	if atomic.LoadInt32(&wp.workCount) == 0 {
		return true
	}
	return false
}

// IncWorker is used as increase worker for working
func (wp *workPool) incWorker(leader *leader, taskData []byte) {

	for {
		if atomic.LoadInt32(&wp.workCount) >= workPoolCap {
			runtime.Gosched()
			continue
		}
		break
	}
	atomic.AddInt32(&wp.workCount, 1)
	go leader.dealRemoteMessage(wp, taskData)
}

// decWorker is used as decrease worker when work is finshed
func (wp *workPool) decWorker() {
	atomic.AddInt32(&wp.workCount, -1)

}
