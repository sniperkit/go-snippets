package inputs

import (
	"strings"
)

//ValueProviderFactory is a factory to create ValueProviders
type ValueProviderFactory interface {
	Query(queryKey string) QueryValueProvider
	Param(paramKey string) ParamValueProvider
}

//NewValueProviderFactory creates a new ValueProviderFactory
func NewValueProviderFactory() ValueProviderFactory {
	return &valueProviderFactory{}
}

type valueProviderFactory struct {
}

func (v *valueProviderFactory) Query(queryKey string) QueryValueProvider {
	return &queryValueProvider{
		key: queryKey,
	}
}

func (v *valueProviderFactory) Param(paramKey string) ParamValueProvider {
	if strings.HasPrefix(paramKey, ":") {
		panic("Should not prefix paramKey with ':' (in ValueProviderFactory)")
	}
	return &paramValueProvider{
		key: paramKey,
	}
}
