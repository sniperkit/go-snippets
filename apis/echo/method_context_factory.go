package echo

import (
	"github.com/labstack/echo"

	"github.com/go-zero-boilerplate/api-resources/resources"
	"github.com/go-zero-boilerplate/api-resources/resources/responses"
)

type MethodContextFactory interface {
	Context(echoCtx echo.Context, path string, resource resources.Resource) resources.MethodContext
}

func NewDefaultMethodContextFactory(responseBuilderFactory responses.BuilderFactory) MethodContextFactory {
	return &defaultMethodContextFactory{
		responseBuilderFactory: responseBuilderFactory,
	}
}

type defaultMethodContextFactory struct {
	responseBuilderFactory responses.BuilderFactory
}

func (d *defaultMethodContextFactory) Context(echoCtx echo.Context, path string, resource resources.Resource) resources.MethodContext {
	return &defaultMethodContext{
		responseBuilderFactory: d.responseBuilderFactory,
		echoCtx:                echoCtx,
	}
}

type defaultMethodContext struct {
	responseBuilderFactory responses.BuilderFactory
	echoCtx                echo.Context
}

func (d *defaultMethodContext) ResponseBuilderFactory() responses.BuilderFactory {
	return d.responseBuilderFactory
}

func (d *defaultMethodContext) MiddlewareContextValue(key string) interface{} {
	return d.echoCtx.Get(key)
}
