package tasks

import (
	"context"

	"github.com/go-zero-boilerplate/facades/logging"
)

//NewContext creates a new Context
func NewContext(ctx context.Context, logger logging.Logger) *Context {
	return &Context{
		ctx:    ctx,
		Logger: logger,
	}
}

//Context holds a logger and db (as inputs) and a Body and Error (as outputs)
type Context struct {
	ctx    context.Context
	Logger logging.Logger

	Response interface{}
	Error    *Error
}

//ChainLoad multiple methods
func (c *Context) ChainLoad(loadFuncs ...LoadFunc) {
	for _, l := range loadFuncs {
		if err := c.ctx.Err(); err != nil {
			c.Error = NewErrorDefaultStatus(err)
			return
		}

		l(c)
		if c.Error != nil {
			return
		}
	}
}
