//Package raw provides implementation of a raw mangos server.
package raw

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
)

type Server struct {
	url  string
	sock mangos.Socket
}

type ResponseHandler func(*Server, *mangos.Message)

//Starts a Raw Socket on the specified url.  A set of workers can run to handle more traffic.
func (self *Server) Listen(url string, workers int, handler ResponseHandler) error {

	self.url = url

	var err error

	if self.sock, err = rep.NewSocket(); err != nil {
		return err
	}

	err = self.sock.SetOption(mangos.OptionRaw, true)
	if err != nil {
		return err
	}

	self.sock.AddTransport(ipc.NewTransport())
	self.sock.AddTransport(tcp.NewTransport())

	if err = self.sock.Listen(url); err != nil {
		return err
	}

	for id := 0; id < workers; id++ {
		go self.processData(handler)
	}

	return nil

}

//Reply to the Raw Request Message.
func (self *Server) Reply(m *mangos.Message) error {

	var err error

	if err = self.sock.SendMsg(m); err != nil {
		return err
	}

	return nil
}

//Handles the raw Request Messages.
func (self *Server) processData(handler ResponseHandler) {

	for {

		msg, err := self.sock.RecvMsg()
		if err != nil {
			continue
		}
		go handler(self, msg)
	}
}
