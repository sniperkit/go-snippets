package busybody

import (
	"fmt"

	"github.com/gdamore/mangos"
	"github.com/gdamore/mangos/protocol/bus"
	"github.com/gdamore/mangos/protocol/pull"
	"github.com/gdamore/mangos/protocol/push"
	"github.com/gdamore/mangos/protocol/rep"
	"github.com/gdamore/mangos/protocol/req"
	"github.com/gdamore/mangos/protocol/surveyor"
	"github.com/gdamore/mangos/transport/ipc"
	"github.com/gdamore/mangos/transport/tcp"
)

func newSurveySocket(options map[string]interface{}) (mangos.Socket, error) {
	var sock mangos.Socket
	var err error

	if sock, err = surveyor.NewSocket(); err != nil {
		return nil, err
	}

	tp := tcp.NewTransport()
	ipctp := ipc.NewTransport()

	sock.AddTransport(tp)
	sock.AddTransport(ipctp)

	for k, v := range options {
		if err := sock.SetOption(k, v); err != nil {
			return nil, fmt.Errorf("error setting option %s: %v", k, err)
		}
	}

	return sock, nil
}

func newBusSocket(options map[string]interface{}) (mangos.Socket, error) {
	var sock mangos.Socket
	var err error

	if sock, err = bus.NewSocket(); err != nil {
		return sock, err
	}

	tp := tcp.NewTransport()
	ipctp := ipc.NewTransport()

	sock.AddTransport(tp)
	sock.AddTransport(ipctp)

	for k, v := range options {
		if err := sock.SetOption(k, v); err != nil {
			return nil, fmt.Errorf("error setting option %s: %v", k, err)
		}
	}

	return sock, nil
}

func newPushSocket(options map[string]interface{}) (mangos.Socket, error) {
	var sock mangos.Socket
	var err error

	if sock, err = push.NewSocket(); err != nil {
		return sock, err
	}

	tp := tcp.NewTransport()
	ipctp := ipc.NewTransport()

	sock.AddTransport(tp)
	sock.AddTransport(ipctp)

	for k, v := range options {
		if err := sock.SetOption(k, v); err != nil {
			return nil, fmt.Errorf("error setting option %s: %v", k, err)
		}
	}

	return sock, nil
}

func newPullSocket(options map[string]interface{}) (mangos.Socket, error) {
	var sock mangos.Socket
	var err error

	if sock, err = pull.NewSocket(); err != nil {
		return sock, err
	}

	tp := tcp.NewTransport()
	ipctp := ipc.NewTransport()

	sock.AddTransport(tp)
	sock.AddTransport(ipctp)

	for k, v := range options {
		if err := sock.SetOption(k, v); err != nil {
			return nil, fmt.Errorf("error setting option %s: %v", k, err)
		}
	}

	return sock, nil
}

func newRequestSocket(options map[string]interface{}) (mangos.Socket, error) {
	var sock mangos.Socket
	var err error

	if sock, err = req.NewSocket(); err != nil {
		return sock, err
	}

	tp := tcp.NewTransport()
	ipctp := ipc.NewTransport()

	sock.AddTransport(tp)
	sock.AddTransport(ipctp)

	for k, v := range options {
		if err := sock.SetOption(k, v); err != nil {
			return nil, fmt.Errorf("error setting option %s: %v", k, err)
		}
	}

	return sock, nil
}

func newReplySocket(options map[string]interface{}) (mangos.Socket, error) {
	var sock mangos.Socket
	var err error

	if sock, err = rep.NewSocket(); err != nil {
		return sock, err
	}

	tp := tcp.NewTransport()
	ipctp := ipc.NewTransport()

	sock.AddTransport(tp)
	sock.AddTransport(ipctp)

	for k, v := range options {
		if err := sock.SetOption(k, v); err != nil {
			return nil, fmt.Errorf("error setting option %s: %v", k, err)
		}
	}

	return sock, nil
}
