package echo

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/go-zero-boilerplate/extended-apex-logger/logging"

	"github.com/go-zero-boilerplate/api-resources/apis"
	"github.com/go-zero-boilerplate/api-resources/servers"
)

//NewBuilder creates a new Server builder with the github.com/labstack/echo framework
func NewBuilder(logger logging.Logger, baseEchoInstance *echo.Echo, apiFactory APIFactory) servers.Builder {
	return &builder{
		logger:           logger,
		baseEchoInstance: baseEchoInstance,
		apiFactory:       apiFactory,
	}
}

type builder struct {
	logger           logging.Logger
	baseEchoInstance *echo.Echo
	apiFactory       APIFactory

	jwtAuthMiddleware echo.MiddlewareFunc
}

func (b *builder) SetDebugMode(debugMode bool) servers.Builder {
	if debugMode {
		b.baseEchoInstance.Debug = debugMode
	}
	return b
}

func (b *builder) PrefixStripperMiddleware(prefixToStrip string) servers.Builder {
	if prefixToStrip == "" {
		return b
	}
	b.logger.WithField("prefix", prefixToStrip).Info("Add middleware to strip request url prefix")
	b.baseEchoInstance.Pre(middlewareRemovePrefix(prefixToStrip))
	return b
}

func (b *builder) LoggerMiddleware() servers.Builder {
	b.baseEchoInstance.Use(middleware.Logger())
	return b
}

func (b *builder) RecoverMiddleware() servers.Builder {
	b.baseEchoInstance.Use(middleware.Recover())
	return b
}

func (b *builder) CORSMiddleware(allowOrigins, allowMethods, allowHeaders []string) servers.Builder {
	b.baseEchoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: allowMethods,
		AllowHeaders: allowHeaders,
	}))
	return b
}

func (b *builder) ResponseDelayMiddleware(responseDelay time.Duration) servers.Builder {
	if responseDelay <= 0 {
		return b
	}
	b.baseEchoInstance.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		time.Sleep(responseDelay)
		return h
	})
	return b
}

func (b *builder) CreateJWTMiddleware(tokenContextKey, secretKey string) servers.Builder {
	b.jwtAuthMiddleware = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(secretKey),
		Skipper:       middleware.DefaultJWTConfig.Skipper,
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    tokenContextKey,
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		Claims:        jwt.MapClaims{},
	})

	return b
}

func (b *builder) AddAPI(publicAPIOut *apis.API) servers.Builder {
	*publicAPIOut = b.apiFactory.NewAPI(b.baseEchoInstance)
	return b
}

func (b *builder) AddJWTAPI(authedAPIOut *apis.API) servers.Builder {
	*authedAPIOut = b.apiFactory.NewAPI(b.baseEchoInstance, b.jwtAuthMiddleware)
	return b
}

func (b *builder) BuildAndRun(address string) error {
	b.logger.WithField("address", address).Info("Server listening")
	return b.baseEchoInstance.Start(address)
}
