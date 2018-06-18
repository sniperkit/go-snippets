package io

import (
	"errors"
	"sync"
)

// io.Pipe copy

var ErrorPipeClosed error = errors.New("io pipe: read/write on closed pipe")

type pipe struct {
	rl    sync.Mutex // gates readers one at a time
	wl    sync.Mutex // gates writers one at a time
	l     sync.Mutex // protects remaining fields
	data  []byte     // data remaining in pending write
	rwait sync.Cond  // waiting reader
	wwait sync.Cond  // waiting writer
	rerr  error      // if reader closed, error to give writer
	werr  error      // if writer closed, error to give reader
}

func (p *pipe) read(b []byte) (n int, err error) {
	// One reader at a time
	p.rl.Lock()
	defer p.rl.Unlock()

	p.l.Lock()
	defer p.l.Unlock()
	for {
		if p.rerr != nil {
			return 0, ErrorPipeClosed
		}
		if p.data != nil {
			break
		}
		if p.werr != nil {
			return 0, p.werr
		}
		p.rwait.Wait()
	}
	n = copy(b, p.data)
	p.data = p.data[n:]
	if len(p.data) == 0 {
		p.data = nil
		p.wwait.Signal()
	}
	return
}

var zero [0]byte

func (p *pipe) write(b []byte) (n int, err error) {
	if b == nil {
		b = zero[:]
	}
	// One writer at a time
	p.wl.Lock()
	defer p.wl.Unlock()

	p.l.Lock()
	defer p.l.Unlock()
	p.data = b
	p.rwait.Signal()
	for {
		if p.data == nil {
			break
		}
		if p.rerr != nil {
			err = p.rerr
			break
		}
		if p.werr != nil {
			err = ErrorPipeClosed
			break
		}
		p.wwait.Wait()
	}
	n = len(b) - len(p.data)
	p.data = nil // in case of rerr or werr
	return
}

func (p *pipe) rClose(err error) {
	if err == nil {
		err = ErrorPipeClosed
	}
	p.l.Lock()
	defer p.l.Unlock()
	p.rerr = err
	p.rwait.Signal()
	p.wwait.Signal()
}

func (p *pipe) wClose(err error) {
	if err == nil {
		err = EOF
	}
	p.l.Lock()
	defer p.l.Unlock()
	p.werr = err
	p.rwait.Signal()
	p.wwait.Signal()
}

// reader
type PipeReader struct {
	p *pipe
}

func (pr *PipeReader) Read(b []byte) (n int, err error) {
	return pr.p.read(b)
}

func (pr *PipeReader) Close() error {
	return pr.CloseWithError(nil)
}

func (pr *PipeReader) CloseWithError(err error) error {
	pr.p.rClose(err)
	return nil
}

// writer
type PipeWriter struct {
	p *pipe
}

func (pw *PipeWriter) Write(b []byte) (n int, err error) {
	return pw.p.write(b)
}

func (pw *PipeWriter) Close() error {
	return pw.CloseWithError(nil)
}

func (pw *PipeWriter) CloseWithError(err error) error {
	pw.p.wClose(err)
	return nil
}

func Pipe() (*PipeReader, *PipeWriter) {
	p := new(pipe)
	p.rwait.L = &p.l
	p.wwait.L = &p.l
	r := &PipeReader{p}
	w := &PipeWriter{p}
	return r, w
}
