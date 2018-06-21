package dispatcher

import (
	//"fmt"
	//"strings"

	"util/logs"

	"core/net"
	"core/net/dispatcher/pb"
)

var _ = logs.Debug

//
type Frame struct {
	*pb.PbFrame
}

// 消息处理单元
type Unit interface {
	Set(srvId, uid string)
	GetId() string
	AddFrame(f *Frame) bool
}

////
//func Url2Part(url string) (srvId, uid string, ok bool) {
//	ss := strings.Split(url, "#")
//	if len(ss) != 2 {
//		return
//	}
//	return ss[0], ss[1], true
//}

////
//func GenUrl(srvId, uid string) string {
//	return fmt.Sprintf("%v#%v", srvId, uid)
//}

//
type BaseUnit struct {
	Id  string // global id
	Url string

	//
	Frames chan *Frame // servers' msg(frame)
	CurF   *Frame      // 当前正在处理的服务端发送过来的消息
}

func NewBaseUnit(frameNum int) *BaseUnit {
	return &BaseUnit{Frames: make(chan *Frame, frameNum)}
}

func (this *BaseUnit) Set(srvId, uid string) {
	this.Id = uid
	this.Url = net.GenUrl(srvId, uid)
}

func (this *BaseUnit) GetId() string {
	return this.Id
}

func (this *BaseUnit) AddFrame(f *Frame) bool {
	select {
	case this.Frames <- f:
		return true
	default:
	}

	return false
}
