package echo

import (
	"github.com/labstack/echo"

	"github.com/go-zero-boilerplate/api-resources/apis"
	apis_echo "github.com/go-zero-boilerplate/api-resources/apis/echo"
)

type APIFactory interface {
	NewAPI(base apis_echo.EchoOrGroup, middlewares ...echo.MiddlewareFunc) apis.API
}

func NewDefaultAPIFactory(methodContextFactory apis_echo.MethodContextFactory) APIFactory {
	return &defaultAPIFactory{
		methodContextFactory: methodContextFactory,
	}
}

type defaultAPIFactory struct {
	methodContextFactory apis_echo.MethodContextFactory
}

func (d *defaultAPIFactory) NewAPI(base apis_echo.EchoOrGroup, middlewares ...echo.MiddlewareFunc) apis.API {
	return apis_echo.NewAPI(d.methodContextFactory, base, middlewares...)
}
