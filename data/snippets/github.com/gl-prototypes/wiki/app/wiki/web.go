package wiki

import (
	"net/http"
	"os"

	"github.com/gliderlabs/stdcom/web"
	"github.com/gorilla/handlers"
)

// Middleware is used to wrap all requests with a common request logger
func (c *Component) Middleware() []web.Middleware {
	return []web.Middleware{
		func(h http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, h)
		},
	}
}
