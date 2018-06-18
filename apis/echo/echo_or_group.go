package echo

import (
	"github.com/labstack/echo"
)

//EchoOrGroup is a common interface between echo.Echo and echo.Group
type EchoOrGroup interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
}
