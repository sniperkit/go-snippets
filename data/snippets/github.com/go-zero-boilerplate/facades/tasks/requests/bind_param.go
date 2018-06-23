package requests

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"

	"github.com/go-zero-boilerplate/facades/tasks"
)

//BindInt64Param gets an int64 from the echo request
func BindInt64Param(c echo.Context, paramName string, val *int64) tasks.LoadFunc {
	return func(taskCtx *tasks.Context) {
		valStr := strings.TrimSpace(c.Param(paramName))
		intVal, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			taskCtx.Error = tasks.NewError(fmt.Errorf("Cannot parse param value '%s' (name %s) as int, error: %s", valStr, paramName, err.Error()), http.StatusBadRequest)
			return
		}

		*val = intVal
	}
}

//BindStringParam gets a string from the echo request
func BindStringParam(c echo.Context, paramName string, val *string) tasks.LoadFunc {
	return func(taskCtx *tasks.Context) {
		*val = strings.TrimSpace(c.Param(paramName))
		if len(*val) == 0 {
			taskCtx.Error = tasks.NewError(fmt.Errorf("Param '%s' is required", paramName), http.StatusBadRequest)
			return
		}
	}
}
