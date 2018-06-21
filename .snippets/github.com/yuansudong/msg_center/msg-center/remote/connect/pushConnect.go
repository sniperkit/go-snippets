package connect

import (
	"fmt"
	"jkt/msg-center/remote/protocol"

	"github.com/go-mangos/mangos/transport/tcp"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/push"
)

type (

	// pushConnect is used as
	pushConnect struct {
		push mangos.Socket
	}
)

// createConnect is used as instance a connect which is of remote
func createConnect(url string) *pushConnect {
	var err error
	pConnect := &pushConnect{}
	if pConnect.push, err = push.NewSocket(); err != nil {
		panic(fmt.Sprintf("create push socket failed, because of %s", err.Error()))
	}
	pConnect.push.AddTransport(tcp.NewTransport())
	pConnect.push.SetOption(mangos.OptionRaw, true)
	pConnect.push.SetOption(mangos.OptionReconnectTime, connectReconnectTime)
	pConnect.push.SetOption(mangos.OptionSendDeadline, connectSendTimeout)
	if err = pConnect.push.Dial(url); err != nil {
		panic(fmt.Sprintf("dial remote address failed, because of %s", err.Error()))
	}
	return pConnect
}

// Send  is used as send message to remote computer,cmd is used as identify the data
func (pc *pushConnect) Send(cmd int, data []byte) error {
	return pc.push.Send(protocol.Packet(cmd, data))
}
