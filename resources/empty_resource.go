package resources

import (
	"github.com/go-zero-boilerplate/api-resources/resources/inputs"
	"github.com/go-zero-boilerplate/api-resources/resources/responses"
)

//EmptyResource just holds default "StatusMethodNotAllowed" methods for Resource methods
type EmptyResource struct{}

//Get returns Status=StatusMethodNotAllowed
func (e *EmptyResource) Get(methodCtx MethodContext, input inputs.Get) responses.Response {
	return methodCtx.ResponseBuilderFactory().Builder().MethodNotAllowed()
}

//Post returns Status=StatusMethodNotAllowed
func (e *EmptyResource) Post(methodCtx MethodContext, input inputs.Post) responses.Response {
	return methodCtx.ResponseBuilderFactory().Builder().MethodNotAllowed()
}

//Patch returns Status=StatusMethodNotAllowed
func (e *EmptyResource) Patch(methodCtx MethodContext, input inputs.Patch) responses.Response {
	return methodCtx.ResponseBuilderFactory().Builder().MethodNotAllowed()
}

//Put returns Status=StatusMethodNotAllowed
func (e *EmptyResource) Put(methodCtx MethodContext, input inputs.Put) responses.Response {
	return methodCtx.ResponseBuilderFactory().Builder().MethodNotAllowed()
}

//Delete returns Status=StatusMethodNotAllowed
func (e *EmptyResource) Delete(methodCtx MethodContext, input inputs.Delete) responses.Response {
	return methodCtx.ResponseBuilderFactory().Builder().MethodNotAllowed()
}
