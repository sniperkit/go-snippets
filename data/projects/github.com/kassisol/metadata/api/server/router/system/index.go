package system

import (
	"net/http"

	"github.com/labstack/echo"
)

func IndexHandle(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
