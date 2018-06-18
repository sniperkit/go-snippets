package apis

import (
	"github.com/go-zero-boilerplate/api-resources/resources"
)

//API is the API server
type API interface {
	AddResource(path string, resource resources.Resource) API
	AddResourceWithEndingParam(path, endingParamName string, resource resources.Resource) API
}
