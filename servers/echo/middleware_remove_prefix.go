package echo

import (
	"strings"

	"github.com/labstack/echo"
)

func middlewareRemovePrefix(prefix string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			url := req.URL
			path := url.Path
			qs := c.QueryString()

			if len(path) > 0 && len(prefix) > 0 && strings.HasPrefix(path, prefix) {
				path = path[len(prefix):]
				uri := path
				if qs != "" {
					uri += "?" + qs
				}

				req.RequestURI = uri
				url.Path = path
			}

			return next(c)
		}
	}
}
