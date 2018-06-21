// base.go
package snapclient

import (
	vxlan "l3/tunnel/vxlan/protocol"
	"sync"
	"utils/logging"

	nanomsg "github.com/op/go-nanomsg"
)

// setup local refs to server info
var serverchannels *vxlan.VxLanConfigChannels
var logger *logging.Writer

type PortEvtCb func(ifindex int32)

// Base Snaproute Interface
type VXLANSnapClient struct {
	vxlan.BaseClientIntf
	thriftmutex         *sync.Mutex
	ribdSubSocket       *nanomsg.SubSocket
	ribdSubSocketCh     chan []byte
	ribdSubSocketErrCh  chan error
	asicdSubSocket      *nanomsg.SubSocket
	asicdSubSocketCh    chan []byte
	asicdSubSocketErrCh chan error
	portUpEventList     map[int32][]vxlan.PortEvtCb
	portDownEventList   map[int32][]vxlan.PortEvtCb
}

func NewVXLANSnapClient(l *logging.Writer) *VXLANSnapClient {
	logger = l

	client := &VXLANSnapClient{
		thriftmutex:         &sync.Mutex{},
		ribdSubSocketCh:     make(chan []byte, 0),
		ribdSubSocketErrCh:  make(chan error, 0),
		asicdSubSocketCh:    make(chan []byte, 0),
		asicdSubSocketErrCh: make(chan error, 0),
		portUpEventList:     make(map[int32][]vxlan.PortEvtCb, 0),
		portDownEventList:   make(map[int32][]vxlan.PortEvtCb, 0),
	}

	go client.ClientChanListener()
	return client

}

func (intf VXLANSnapClient) ClientChanListener() {

	for {
		select {
		case rxBuf := <-intf.asicdSubSocketCh:
			intf.processAsicdNotification(rxBuf)
		case <-intf.asicdSubSocketErrCh:
			continue
		case rxBuf := <-intf.ribdSubSocketCh:
			intf.processRibdNotification(rxBuf)
		case <-intf.ribdSubSocketErrCh:
			continue
		}
	}
}

func (v VXLANSnapClient) SetServerChannels(s *vxlan.VxLanConfigChannels) {
	serverchannels = s
}

func (v VXLANSnapClient) IsClientIntfType(client vxlan.VXLANClientIntf, clientStr string) bool {
	logger.Debug("IsClientIntfType", clientStr)
	switch client.(type) {
	case VXLANSnapClient:
		if clientStr == "SnapClient" {
			return true
		}
	default:
		logger.Debug("IsClientInfType: default did not find client type")
	}
	return false
}
