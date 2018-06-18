package requests

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/go-zero-boilerplate/facades/tasks"
)

//BindBody binds the request body to the given object
func BindBody(c echo.Context, dest interface{}) tasks.LoadFunc {
	return func(taskCtx *tasks.Context) {
		if err := c.Bind(dest); err != nil {
			taskCtx.Error = tasks.NewError(err, http.StatusBadRequest)
			return
		}
	}
}
