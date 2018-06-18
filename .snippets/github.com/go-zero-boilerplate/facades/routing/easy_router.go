package routing

import (
	"context"
	"fmt"
	"strings"

	"github.com/labstack/echo"

	"github.com/go-zero-boilerplate/facades/logging"
	"github.com/go-zero-boilerplate/facades/routing/responses"
	"github.com/go-zero-boilerplate/facades/tasks"
)

var _easyRouter = &EasyRouterStruct{}

//EasyRouter returns the local/private instance/singleton of _easyRouter
func EasyRouter() *EasyRouterStruct {
	return _easyRouter
}

//EasyRouterStruct is a simple struct to manage registering of routes
type EasyRouterStruct struct {
	router               *echo.Echo
	pendingRegistrations []*PendingRegistration
}

func (e *EasyRouterStruct) register(method, path string) *PendingRegistration {
	pending := &PendingRegistration{
		method: method,
		path:   path,
	}
	e.pendingRegistrations = append(e.pendingRegistrations, pending)
	return pending
}

//GET will register a GET method to the path with the given handler and middlewares
func (e *EasyRouterStruct) GET(path string) *PendingRegistration {
	return e.register("GET", path)
}

//POST will register a POST method to the path with the given handler and middlewares
func (e *EasyRouterStruct) POST(path string) *PendingRegistration {
	return e.register("POST", path)
}

//PUT will register a PUT method to the path with the given handler and middlewares
func (e *EasyRouterStruct) PUT(path string) *PendingRegistration {
	return e.register("PUT", path)
}

//DELETE will register a DELETE method to the path with the given handler and middlewares
func (e *EasyRouterStruct) DELETE(path string) *PendingRegistration {
	return e.register("DELETE", path)
}

func (e *EasyRouterStruct) finalizeRegistration() {
	for _, pending := range e.pendingRegistrations {
		if pending.handler == nil {
			panic(fmt.Sprintf("pending.handler required for route %s", pending.path))
		}

		method := pending.method
		path := pending.path
		handler := pending.handler
		preMiddlewares := pending.preMiddlewares
		postMiddlewares := pending.postMiddlewares

		var (
			middlewares []echo.MiddlewareFunc
			echoHandler echo.HandlerFunc
		)

		logging.Get().Info(fmt.Sprintf("Registering route %s, method %s", path, method))

		middlewares = append(middlewares, preMiddlewares...)
		middlewares = append(middlewares, jwtMiddlewares...)

		//Determine echoHandler based on type of handler created from handler
		echoHandler = func(c echo.Context) error {
			/*logger := logging.Get().WithFields(map[string]interface{}{
				"request-uri": c.Request().RequestURI,
			})*/

			formatter := c.JSON
			// TODO: using default JSON formatter here, might need to override this in future

			logger := logging.Get().WithFields(map[string]interface{}{
				"path": c.Path(),
				"uri":  c.Request().RequestURI,
			})
			taskCtx := tasks.NewContext(context.TODO(), logger)
			handler(c, taskCtx)
			if taskCtx.Error != nil {
				return responses.NewResponseFromTaskError(taskCtx.Error, formatter).Response()
			}

			if taskCtx.Response == nil {
				//no content
				return responses.NewResponse(responses.StatusNoContent, nil, formatter).Response()
			}
			return responses.NewResponse(responses.StatusOK, taskCtx.Response, formatter).Response()
		}

		middlewares = append(middlewares, postMiddlewares...)

		switch strings.ToUpper(method) {
		case "CONNECT":
			e.router.CONNECT(path, echoHandler, middlewares...)
		case "DELETE":
			e.router.DELETE(path, echoHandler, middlewares...)
		case "GET":
			e.router.GET(path, echoHandler, middlewares...)
		case "HEAD":
			e.router.HEAD(path, echoHandler, middlewares...)
		case "OPTIONS":
			e.router.OPTIONS(path, echoHandler, middlewares...)
		case "PATCH":
			e.router.PATCH(path, echoHandler, middlewares...)
		case "POST":
			e.router.POST(path, echoHandler, middlewares...)
		case "PUT":
			e.router.PUT(path, echoHandler, middlewares...)
		case "TRACE":
			e.router.TRACE(path, echoHandler, middlewares...)
		default:
			panic(fmt.Errorf("Http method '%s' is not yet supported", method))
		}
	}
}

var _easyRouterIsInitialized = false

//InitEasyRouter initializes the local/private instance/singleton of _easyRouter
func InitEasyRouter() {
	_easyRouter.router = router
	_easyRouter.finalizeRegistration()

	_easyRouterIsInitialized = true
}
