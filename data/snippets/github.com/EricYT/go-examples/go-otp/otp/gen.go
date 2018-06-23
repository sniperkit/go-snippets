package otp

import (
	"errors"
	"log"
	"time"
)

// gen.erl implement

var (
	ErrorGenCallTimeout error = errors.New("gen: gen call timeout")
	ErrorGenCallResult  error = errors.New("gen: gen call wrong result")
	ErrorGenCast        error = errors.New("gen: gen cast error")
	ErrorGenInfo        error = errors.New("gen: gen info error")
	ErrorGenTerminate   error = errors.New("gen: gen terminate error")
)

type gen struct {
	genMod *genServer
}

func newGen(name Name, genMod *genServer, args ...interface{}) error {
	g := &gen{
		genMod: genMod,
	}

	errCh := make(chan error)
	// start gen loop goroutine
	go func() {
		// FIXME: recover
		defer g.genMod.tomb.Done()
		g.genMod.tomb.Kill(g.initIt(errCh, name, args...))
	}()

	// wait this gen started
	return <-errCh
}

func (g *gen) initIt(errCh chan<- error, name Name, args ...interface{}) error {
	// register gen server by name
	if err := registerName(name, g); err != nil {
		errCh <- err
		return err
	}
	defer unregisterName(name)

	// gen started
	errCh <- nil

	// callback initialize gen module
	if err := g.genMod.InitIt(args...); err != nil {
		return err
	}
	return nil
}

// functions for client call server
func Call(server Name, req interface{}) (interface{}, error) {
	gen, err := getGenByName(server)
	if err != nil {
		log.Printf("[gen:Call] server: %s req: %#v get by name error: %s\n", server, req, err)
		return nil, err
	}
	resCh := make(chan interface{}, 1)
	call := call{req, resCh}
	gen.genMod.mailbox <- call
	return <-resCh, nil
}

func CallWithTimeout(server Name, req interface{}, timeout time.Duration) (interface{}, error) {
	gen, err := getGenByName(server)
	if err != nil {
		log.Printf("[gen:CallWithTimeout] server: %s req: %#v get by name error: %s\n", server, req, err)
		return nil, err
	}

	// timeout
	timer := time.NewTimer(timeout)
	resCh := make(chan interface{}, 1)
	call := call{req, resCh}

	// send call
	select {
	case gen.genMod.mailbox <- call:
	case <-timer.C:
		return nil, ErrorGenCallTimeout
	}

	// receive
	select {
	case result, ok := <-resCh:
		if ok {
			return result, nil
		}
		return nil, ErrorGenCallResult
	case <-timer.C:
		return nil, ErrorGenCallTimeout
	}
}

func Cast(server Name, req interface{}) error {
	gen, err := getGenByName(server)
	if err != nil {
		log.Printf("[gen:Cast] server: %s req: %#v get by name error: %s\n", server, req, err)
		return err
	}
	select {
	case gen.genMod.mailbox <- cast{req}:
		return nil
	default:
		return ErrorGenCast
	}
}

func Info(server Name, req interface{}) error {
	gen, err := getGenByName(server)
	if err != nil {
		log.Printf("[gen:Info] server: %s req: %#v get by name error: %s\n", server, req, err)
		return err
	}
	select {
	case gen.genMod.mailbox <- info{req}:
		return nil
	default:
		return ErrorGenInfo
	}
}

func Terminate(server Name, reason interface{}) error {
	gen, err := getGenByName(server)
	if err != nil {
		log.Printf("[gen:Terminate] server: %s req: %#v get by name error: %s\n", server, reason, err)
		return err
	}
	errCh := make(chan error, 1)
	gen.genMod.mailbox <- exit{reason, errCh}
	return <-errCh
}
