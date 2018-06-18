package routing

import (
	"time"

	"github.com/labstack/echo"
)

func MiddlewareResponseDelay(duration time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			time.Sleep(duration)
			return next(c)
		}
	}
}
