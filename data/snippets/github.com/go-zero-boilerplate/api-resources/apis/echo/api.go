package echo

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"

	"github.com/go-zero-boilerplate/api-resources/apis"
	"github.com/go-zero-boilerplate/api-resources/resources"
	"github.com/go-zero-boilerplate/api-resources/resources/responses"
)

//NewAPI creates a new echo API
func NewAPI(methodContextFactory MethodContextFactory, base EchoOrGroup, middlewares ...echo.MiddlewareFunc) apis.API {
	return &api{
		methodContextFactory: methodContextFactory,
		base:                 base,
		middlewares:          middlewares,
	}
}

type api struct {
	methodContextFactory MethodContextFactory
	base                 EchoOrGroup
	middlewares          []echo.MiddlewareFunc
}

func (a *api) resourceResponseToEchoResponse(echoCtx echo.Context, response responses.Response) EchoResponse {
	responseFormat := responses.FormatJSON

	switch responseFormat {
	case responses.FormatJSON:
		return EchoResponse(echoCtx.JSON(response.Status.Int(), response.Body))
	case responses.FormatXML:
		return EchoResponse(echoCtx.XML(response.Status.Int(), response.Body))
	default:
		//For now JSON is the default "fallback" representation type
		return EchoResponse(echoCtx.JSON(response.Status.Int(), response.Body))
	}
}

func (a *api) AddResource(path string, resource resources.Resource) apis.API {
	a.base.GET(
		path,
		func(echoCtx echo.Context) error {
			ctx := a.methodContextFactory.Context(echoCtx, path, resource)
			return a.resourceResponseToEchoResponse(echoCtx, resource.Get(ctx, newGetInput(echoCtx)))
		},
		a.middlewares...,
	)

	a.base.POST(
		path,
		func(echoCtx echo.Context) error {
			ctx := a.methodContextFactory.Context(echoCtx, path, resource)
			return a.resourceResponseToEchoResponse(echoCtx, resource.Post(ctx, newPostInput(echoCtx)))
		},
		a.middlewares...,
	)

	a.base.PATCH(
		path,
		func(echoCtx echo.Context) error {
			ctx := a.methodContextFactory.Context(echoCtx, path, resource)
			return a.resourceResponseToEchoResponse(echoCtx, resource.Patch(ctx, newPatchInput(echoCtx)))
		},
		a.middlewares...,
	)

	a.base.PUT(
		path,
		func(echoCtx echo.Context) error {
			ctx := a.methodContextFactory.Context(echoCtx, path, resource)
			return a.resourceResponseToEchoResponse(echoCtx, resource.Put(ctx, newPutInput(echoCtx)))
		},
		a.middlewares...,
	)

	a.base.DELETE(
		path,
		func(echoCtx echo.Context) error {
			ctx := a.methodContextFactory.Context(echoCtx, path, resource)
			return a.resourceResponseToEchoResponse(echoCtx, resource.Delete(ctx, newDeleteInput(echoCtx)))
		},
		a.middlewares...,
	)

	return a
}

func (a *api) AddResourceWithEndingParam(path, endingParamName string, resource resources.Resource) apis.API {
	if !strings.HasPrefix(endingParamName, ":") {
		panic("Should prefix endingParamName with ':' (in echo/api)")
	}
	return a.
		AddResource(strings.TrimRight(path, "/"), resource).
		AddResource(fmt.Sprintf("%s/%s", strings.TrimRight(path, "/"), strings.Trim(endingParamName, "/")), resource)
}
