// 消息分发器
package dispatcher

import (
	"fmt"
	"sync"
	"time"

	"util/logs"

	"core/net"
	"core/net/dispatcher/pb"
)

var _ = logs.Debug

// 消息分发器
type Dispatcher struct {
	Name      string // 标识
	SrvId     string
	Timestamp int64

	Id         int             // global id
	Units      map[string]Unit // global id => unit
	sync.Mutex                 //
}

func New(name string, srvId string) *Dispatcher {
	return &Dispatcher{
		Name:      name,
		SrvId:     srvId,
		Timestamp: time.Now().Unix(),
		Units:     map[string]Unit{},
	}
}

//
func (this *Dispatcher) AddUnit(u Unit) {
	this.Lock()
	defer this.Unlock()

	this.Id++
	uid := fmt.Sprintf("%v:%v", this.Timestamp, this.Id)

	u.Set(this.SrvId, uid)
	this.Units[uid] = u
}

//
func (this *Dispatcher) RemoveUnit(u Unit) {
	this.Lock()
	defer this.Unlock()

	delete(this.Units, u.GetId())
}

//
func (this *Dispatcher) GetUnit(id string) Unit {
	this.Lock()
	defer this.Unlock()

	return this.Units[id]
}

//
func (this *Dispatcher) Dispatch(f *pb.PbFrame, fOffline func(dstUrl string)) {
	for _, url := range f.DstUrls {
		//
		_, id, ok := net.Url2Part(url)
		if !ok {
			logs.Warn("invalid dst url in frame! %v:%v", this.Name, url)
			continue
		}

		unit := this.GetUnit(id)
		if nil == unit {
			// 已经不在线了
			logs.Info("not found unit! maybe offline. %v:%v", this.Name, url)

			// 向发送服务器反馈
			fOffline(url)
			continue
		}

		nf := &Frame{PbFrame: f}
		unit.AddFrame(nf)
	}
}

//
var g_dispatcher *Dispatcher

func Init(name, srvId string) {
	g_dispatcher = New(name, srvId)
}

//
func AddUnit(u Unit) {
	g_dispatcher.AddUnit(u)
}

//
func RemoveUnit(u Unit) {
	g_dispatcher.RemoveUnit(u)
}

func Dispatch(f *pb.PbFrame, fOffline func(dstUrl string)) {
	g_dispatcher.Dispatch(f, fOffline)
}
