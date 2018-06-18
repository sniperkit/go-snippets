package routing

import (
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var router *echo.Echo

//Router will return the singleton router instance
func Router() *echo.Echo { return router }

//InitRouter will initialize the router
func InitRouter(version, gitSha1 string, debugMode bool, responseDelay time.Duration, sentryClient *raven.Client) {
	router = echo.New()
	router.Debug = debugMode

	if responseDelay > 0 {
		router.Use(MiddlewareResponseDelay(responseDelay))
	}

	router.Use(middleware.Logger())

	router.Use(MiddlewareRecoverSentry(sentryClient, nil))

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	router.Use(MiddlewareResponseVersion(version, gitSha1))
}
