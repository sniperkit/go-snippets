package routing

import (
	"fmt"

	raven "github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"

	"github.com/go-zero-boilerplate/facades/debugging"
	"github.com/go-zero-boilerplate/facades/logging"
)

//MiddlewareRecoverSentry returns a middleware which recovers from panics anywhere in the chain
//and sends the event to Sentry using the client
//and handles the control to the centralized HTTPErrorHandler
func MiddlewareRecoverSentry(client *raven.Client, config *middleware.RecoverConfig) echo.MiddlewareFunc {
	if config == nil {
		config = &middleware.DefaultRecoverConfig
	}
	// Defaults
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = middleware.DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}
					if !config.DisablePrintStack {
						c.Logger().Printf("[%s] error: %s\n", color.Red("PANIC RECOVER"), err)
						logging.Get().WithError(err).Error(fmt.Sprintf("PANIC RECOVER, stack: %s", debugging.GetNormalStack(false)))
						if client != nil {
							client.CaptureError(err, nil)
						}
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
