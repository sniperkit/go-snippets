package routing

import (
	"github.com/labstack/echo"
)

//PendingRegistration holds the route that will get registered
type PendingRegistration struct {
	method string
	path   string

	handler EasyHandlerFunc

	preMiddlewares  []echo.MiddlewareFunc
	postMiddlewares []echo.MiddlewareFunc
}

func (p *PendingRegistration) Handler(h EasyHandlerFunc) *PendingRegistration {
	p.handler = h
	return p
}
