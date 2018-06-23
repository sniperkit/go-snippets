package routing

import (
	"github.com/labstack/echo"
)

//MiddlewareResponseVersion is a middleware to add X-VERSION and X-SHA1 to response headers
func MiddlewareResponseVersion(version, gitSha1 string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()
			res.Header().Add("X-VERSION", version)
			res.Header().Add("X-SHA1", gitSha1)
			return next(c)
		}
	}
}
