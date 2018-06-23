package resources

import (
	"github.com/go-zero-boilerplate/api-resources/resources/inputs"
	"github.com/go-zero-boilerplate/api-resources/resources/responses"
)

//Resource represents any resource that gets transferred from/to the API
type Resource interface {
	Get(methodCtx MethodContext, input inputs.Get) responses.Response
	Post(methodCtx MethodContext, input inputs.Post) responses.Response
	Patch(methodCtx MethodContext, input inputs.Patch) responses.Response
	Put(methodCtx MethodContext, input inputs.Put) responses.Response
	Delete(methodCtx MethodContext, input inputs.Delete) responses.Response
}
