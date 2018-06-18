package routing

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	//JWTContextKey is used to set the ContextKey of the JWT middleware
	JWTContextKey = "user_token"
)

var (
	jwtMiddlewares = []echo.MiddlewareFunc{}
)

//AddJWTMiddleware will add a JWT middleware
func AddJWTMiddleware(jwtSecretKey string, skipper func(c echo.Context) bool) {
	if _easyRouterIsInitialized {
		panic("AddJWTMiddleware should be called (at startup) before InitEasyRouter")
	}

	jwtMiddlewares = append(jwtMiddlewares, middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(jwtSecretKey),
		Skipper:       skipper,
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    JWTContextKey,
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		Claims:        jwt.MapClaims{},
	}))
}
