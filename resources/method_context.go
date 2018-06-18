package resources

import (
	"github.com/go-zero-boilerplate/api-resources/resources/responses"
)

type MethodContext interface {
	MiddlewareContextValue(key string) interface{} //typically info (like user object) stored between middlewares

	ResponseBuilderFactory() responses.BuilderFactory
}
