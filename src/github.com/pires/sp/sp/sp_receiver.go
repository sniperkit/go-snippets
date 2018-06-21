package sp

import (
	"github.com/gdamore/mangos"
	"github.com/gdamore/mangos/protocol/rep"
	"github.com/gdamore/mangos/transport/ipc"
	"github.com/gdamore/mangos/transport/tcp"
)

type SPReceiver struct{}

func (receiver *SPReceiver) Receive(address string) error {
	var sock mangos.Socket
	var err error

	if sock, err = rep.NewSocket(); err != nil {
		return err
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err := sock.Listen(address); err != nil {
		return err
	}

	for {
		if _, err := sock.Recv(); err != nil {
			return err
		}
	}

	return nil
}
