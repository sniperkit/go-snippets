package otp

import (
	"context"
	"errors"
	"reflect"

	tomb "gopkg.in/tomb.v1"
)

// gen_server implement

// errors
var (
	ErrorGenServerMailboxClosed error = errors.New("gen server: mailbox closed")
	ErrorGenServerUnknowMsg     error = errors.New("gen server: receive unknow message type")
	ErrorGenServerTerminate     error = errors.New("gen server: terminate")
)

// gen_server return operation values
type Return string
type Name string

const (
	OK      Return = "ok"
	STOP    Return = "stop"
	REPLY   Return = "reply"
	NOREPLY Return = "noreply"
)

type GenServer interface {
	// .Init initialize function
	Init(args ...interface{}) (ctx context.Context, err error)

	// handle client sync call
	HandleCall(ctx context.Context, req interface{}) (context.Context, interface{}, error)

	// handle client async call
	HandleCast(ctx context.Context, req interface{}) (context.Context, error)

	// handle client info
	HandleInfo(ctx context.Context, req interface{}) (context.Context, error)

	// terminate
	Terminate(ctx context.Context, reason interface{}) error
}

type genServer struct {
	tomb tomb.Tomb

	// mailbox
	mailbox chan interface{}

	name Name
	mod  GenServer
}

func NewGenServer(name Name, mod GenServer, args ...interface{}) error {
	g := &genServer{
		mailbox: make(chan interface{}, 100), //default mailbox buffer is set to 100
		name:    name,
		mod:     mod,
	}
	// start background gen module
	return newGen(name, g, args)
}

func (g *genServer) InitIt(args ...interface{}) error {
	ctx, err := g.mod.Init(args)
	if err != nil {
		return err
	}

	// let's rock!
	if err = g.loop(ctx); err != nil {
		return err
	}

	return nil
}

func (g *genServer) loop(ctx context.Context) (err error) {
	// start gen server loop
	for {
		select {
		case msg, ok := <-g.mailbox:
			if !ok {
				//FIXME: someone close this channel ?
				return ErrorGenServerMailboxClosed
			}
			// decode message
			if ctx, err = g.decodeMsg(ctx, msg); err != nil {
				return err
			}
		case <-g.tomb.Dying():
			// go die
			return nil
		}
	}
}

var (
	MsgTypeCall reflect.Type = reflect.TypeOf((*call)(nil)).Elem()
	MsgTypeCast reflect.Type = reflect.TypeOf((*cast)(nil)).Elem()
	MsgTypeInfo reflect.Type = reflect.TypeOf((*info)(nil)).Elem()
	MsgTypeExit reflect.Type = reflect.TypeOf((*exit)(nil)).Elem()
)

func (g *genServer) decodeMsg(ctx context.Context, msg interface{}) (context.Context, error) {
	typ := GetType(msg)
	switch typ {
	case MsgTypeExit:
		return g.terminate(ctx, msg)
	default:
		return g.handleMsg(ctx, msg)
	}
}

func (g *genServer) handleMsg(ctx context.Context, msg interface{}) (context.Context, error) {
	// operate special messages
	typ := GetType(msg)
	switch typ {
	case MsgTypeCall:
		return g.handleCall(ctx, msg)
	case MsgTypeCast:
		return g.handleCast(ctx, msg)
	case MsgTypeInfo:
		return g.handleInfo(ctx, msg)
	default:
		//FIXME: unknow message type receive
	}
	return ctx, ErrorGenServerUnknowMsg
}

func (g *genServer) handleCall(ctx context.Context, msg interface{}) (context.Context, error) {
	var res interface{}
	var err error

	call_, _ := msg.(call)
	ctx, res, err = g.mod.HandleCall(ctx, call_.req)
	if err != nil {
		return ctx, err
	}

	// send back result
	asyncSend(call_.resCh, res)

	return ctx, nil
}

func (g *genServer) handleCast(ctx context.Context, msg interface{}) (context.Context, error) {
	var err error
	cast_, _ := msg.(cast)
	ctx, err = g.mod.HandleCast(ctx, cast_.req)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (g *genServer) handleInfo(ctx context.Context, msg interface{}) (context.Context, error) {
	var err error
	cast_, _ := msg.(info)
	ctx, err = g.mod.HandleInfo(ctx, cast_.req)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (g *genServer) terminate(ctx context.Context, msg interface{}) (context.Context, error) {
	exit_, _ := msg.(exit)
	if err := g.mod.Terminate(ctx, exit_.reason); err != nil {
		exit_.errCh <- err
		return ctx, err
	}
	// always done
	exit_.errCh <- ErrorGenServerTerminate
	return ctx, ErrorGenServerTerminate
}

type call struct {
	req   interface{}
	resCh chan interface{}
}

type cast struct {
	req interface{}
}

type info struct {
	req interface{}
}

type exit struct {
	reason interface{}
	errCh  chan error
}

// tool functions
func GetType(val interface{}) reflect.Type {
	typ := reflect.TypeOf(val)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func asyncSend(returnCh chan<- interface{}, val interface{}) {
	select {
	case returnCh <- val:
	default:
	}
}
