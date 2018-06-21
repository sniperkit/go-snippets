package sentinel

import (
	"fmt"
	"node/protocol"
	"node/router"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/transport/tcp"
)

const (
	sentinelDefaultTimeout = time.Millisecond
	sentinelAddress        = "tcp://*:5050"
)

type (
	// leader is used as description leader
	leader struct {
		// signal is used as receive system signal, example,kill,interrupt
		signalChannel chan os.Signal
		// childrenNumber is used as record which is a count of corotinue
		workPoolCount int32
		// waitGroup is used as wait  all workpool quit
		waitGroup sync.WaitGroup
		//  pull is used as pull message from remote computer
		pull mangos.Socket
		//  router is used as call function from map
		router *router.Router
	}
)

// createLeader is used as create a instance of workpool
func createLeader(url string) *leader {
	l := &leader{
		signalChannel: make(chan os.Signal, 1),
		workPoolCount: 0,
		router:        router.CreateRouter(),
	}
	signal.Notify(l.signalChannel, os.Interrupt, os.Kill)
	l.initSocket(url)
	return l
}

func (l *leader) initSocket(url string) {
	var err error
	if l.pull, err = pull.NewSocket(); err != nil {
		panic(fmt.Sprintf("create pull socket failed,because of %s", err.Error()))
	}
	l.pull.AddTransport(tcp.NewTransport())
	if err = l.pull.SetOption(mangos.OptionRaw, true); err != nil {
		panic(fmt.Sprintf("pull socket set raw mode failed,because of %s ", err.Error()))
	}
	if err = l.pull.SetOption(mangos.OptionRecvDeadline, sentinelDefaultTimeout); err != nil {
		panic(fmt.Sprintf("pull socket set timeout mode failed,because of %s", err.Error()))
	}
	if err = l.pull.Listen(url); err != nil {
		panic(fmt.Sprintf("pull socket listen address failed,because of %s", err.Error()))
	}
}

// Run is used as Run leader to receive message which is from remote computer
func (l *leader) Run() {
	var signal os.Signal
	cpuCount := runtime.NumCPU()
	for i := 0; i < cpuCount; i++ {
		l.waitGroup.Add(1)
		go createWorkPool(i, l)
	}
	fmt.Println("我执行到了leader run")
	for {
		select {
		case signal = <-l.signalChannel:
			fmt.Println("当前接受到的信号是:", signal.String())
			l.pull.Close()
			goto end
		default:
		}
		runtime.Gosched()
	}
end:
	fmt.Println("我在等待着你们退出了")
	l.waitGroup.Wait()
}

func (l *leader) dealRemoteMessage(wp *workPool, requestData []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	l.router.Execute(protocol.UnPacket(requestData))
	wp.decWorker()
}
