package sp

import (
	"github.com/pires/sp/common"

	"github.com/gdamore/mangos"
	"github.com/gdamore/mangos/protocol/req"
	"github.com/gdamore/mangos/transport/ipc"
	"github.com/gdamore/mangos/transport/tcp"
)

type SPSender struct{}

func (sender *SPSender) Send(message *common.Message, address string) error {
	var sock mangos.Socket
	var err error

	if sock, err = req.NewSocket(); err != nil {
		return err
	}
	defer sock.Close()

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial(address); err != nil {
		return err
	}

	return sock.Send(message.ToBytes())
}
