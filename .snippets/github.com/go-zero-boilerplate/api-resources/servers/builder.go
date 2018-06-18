package servers

import (
	"time"

	"github.com/go-zero-boilerplate/api-resources/apis"
)

//Builder is used to simplify the creation of a basic API server
type Builder interface {
	SetDebugMode(debugMode bool) Builder
	PrefixStripperMiddleware(prefixToStrip string) Builder
	LoggerMiddleware() Builder
	RecoverMiddleware() Builder
	CORSMiddleware(allowOrigins, allowMethods, allowHeaders []string) Builder
	ResponseDelayMiddleware(responseDelay time.Duration) Builder
	CreateJWTMiddleware(tokenContextKey, secretKey string) Builder
	AddAPI(publicAPIOut *apis.API) Builder
	AddJWTAPI(authedAPIOut *apis.API) Builder

	BuildAndRun(address string) error
}
