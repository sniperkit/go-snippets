package routing

import (
	"github.com/labstack/echo"

	"github.com/go-zero-boilerplate/facades/tasks"
)

//EasyHandlerFunc is a simple func receiving a echo.Context and a tasks.Context
type EasyHandlerFunc func(c echo.Context, taskCtx *tasks.Context)
