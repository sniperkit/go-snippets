package mergectx

import (
	"context"
	"time"
)

// Context ...
type Context struct {
	ctx1, ctx2   context.Context
	done1, done2 <-chan struct{}
	done         <-chan struct{}
	deadline     time.Time
	deadlineOk   bool
}

var _ context.Context = (*Context)(nil)

// New ...
func New(ctx1, ctx2 context.Context) context.Context {
	ctx := &Context{
		ctx1:  ctx1,
		ctx2:  ctx2,
		done1: ctx1.Done(),
		done2: ctx2.Done(),
	}
	ctx.initDone()
	ctx.initDeadLine()

	return ctx
}

func (ctx *Context) initDone() {
	done := make(chan struct{})
	ctx.done = done
	go func() {
		select {
		case <-ctx.done1:
		case <-ctx.done2:
		}
		close(done)
	}()
}

func (ctx *Context) initDeadLine() {
	d1, ok1 := ctx.ctx1.Deadline()
	d2, ok2 := ctx.ctx2.Deadline()
	ctx.deadlineOk = true

	switch {
	case ok1 && ok2:
		if d1.Before(d2) {
			ctx.deadline = d1
		}
		ctx.deadline = d2
	case ok1:
		ctx.deadline = d1
	case ok2:
		ctx.deadline = d2
	default:
		ctx.deadline = time.Time{}
		ctx.deadlineOk = false
	}
}

// Deadline ...
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.deadline, ctx.deadlineOk
}

// Done ...
func (ctx *Context) Done() <-chan struct{} {
	return ctx.done
}

// Err ...
func (ctx *Context) Err() error {
	err1 := ctx.ctx1.Err()
	err2 := ctx.ctx2.Err()
	if err1 != nil {
		return err1
	}
	return err2
}

// Value ...
func (ctx *Context) Value(key interface{}) interface{} {
	if v := ctx.ctx1.Value(key); v != nil {
		return v
	}
	return ctx.ctx2.Value(key)
}
