package resolvers

import (
	"context"
	"fmt"
	"sync"

	"github.com/distributed-vision/go-resources/ids"
)

type Selector interface {
	Type() ids.TypeIdentifier
	Key() interface{}
	Test(entity interface{}) bool
}

type Resolver interface {
	Get(resolutionContext context.Context, selector Selector) (interface{}, error)
	Resolve(resolutionContext context.Context, selector Selector) (chan interface{}, chan error)
	ResolverInfo() ResolverInfo
}

type MutableResolver interface {
	Resolver
	Put(resolutionContext context.Context, entity interface{}) (interface{}, error)
	Post(resolutionContext context.Context, entity interface{}) (interface{}, error)
	Delete(resolutionContext context.Context, selector Selector) error
}

type ResolverInfo interface {
	ResolverType() ids.TypeIdentifier
	IsMutable() bool
	ResolvableTypes() []ids.TypeIdentifier
	ResolvableDomains() []ids.Domain
	KeyExtractor() KeyExtractor
	Matches(selector Selector) bool
	Value(key interface{}) interface{}
	WithValue(key, value interface{}) ResolverInfo
	WithValues(values map[interface{}]interface{}) ResolverInfo
	WithExtractor(keyExtractor KeyExtractor) ResolverInfo
	DerivedCopy() ResolverInfo
}

type ResolverFactory interface {
	ResolverType() ids.TypeIdentifier
	ResolverInfo() ResolverInfo
	New(resolutionContext context.Context) (Resolver, error)
}

type KeyExtractor func(entity ...interface{}) (interface{}, bool)

var RootInfo = NewResolverInfo(nil, nil, nil, nil, nil)
var rootResolver, _ = NewCompositeResolver(RootInfo)

func Get(resolutionContext context.Context, selector Selector) (entity interface{}, err error) {
	return rootResolver.Get(resolutionContext, selector)
}

func Resolve(resolutionContext context.Context, selector Selector) (chan interface{}, chan error) {
	return rootResolver.Resolve(resolutionContext, selector)
}

func RegisterResolver(resolver Resolver) {
	rootResolver.RegisterComponent(resolver)
}

func RegisterFactory(resolverFactory ResolverFactory) {
	rootResolver.RegisterComponentFactory(resolverFactory, false)
}

type NewFactoryFunction func(resolverInfo ResolverInfo) (ResolverFactory, error)

var newFactoryFunctionRegistry map[string]NewFactoryFunction = make(map[string]NewFactoryFunction)
var newFactoryFunctionRegistryMutex = &sync.Mutex{}

func NewResolverFactory(resolverInfo ResolverInfo) (ResolverFactory, error) {
	newFactoryFunctionRegistryMutex.Lock()
	defer newFactoryFunctionRegistryMutex.Unlock()
	//fmt.Printf("Getting func for: %v\n", resolverInfo.ResolverType())
	if newFactoryFunction, ok := newFactoryFunctionRegistry[string(resolverInfo.ResolverType().Value())]; ok {
		return newFactoryFunction(resolverInfo)
	}

	return nil, fmt.Errorf("No factory function availible for: %s", string(resolverInfo.ResolverType().Id()))
}

func ResisterNewFactoryFunction(resolverType ids.TypeIdentifier, newFactoryFunction NewFactoryFunction) {
	newFactoryFunctionRegistryMutex.Lock()
	//fmt.Printf("Regestering func for: %v\n", resolverType)
	newFactoryFunctionRegistry[string(resolverType.Value())] = newFactoryFunction
	newFactoryFunctionRegistryMutex.Unlock()
}

type resolverInfo struct {
	resolverType      ids.TypeIdentifier
	resolvableTypes   []ids.TypeIdentifier
	resolvableDomains []ids.Domain
	keyExtractor      KeyExtractor
	values            map[interface{}]interface{}
	parent            ResolverInfo
}

func NewResolverInfo(resolverType ids.TypeIdentifier, resolvableTypes []ids.TypeIdentifier, resolvableDomains []ids.Domain, keyExtractor KeyExtractor, values map[interface{}]interface{}) ResolverInfo {
	return &resolverInfo{resolverType, resolvableTypes, resolvableDomains, keyExtractor, values, nil}
}

func (this *resolverInfo) ResolverType() ids.TypeIdentifier {
	return this.resolverType
}

func (this *resolverInfo) ResolvableTypes() []ids.TypeIdentifier {
	return this.resolvableTypes
}

func (this *resolverInfo) IsMutable() bool {
	return false
}

func (this *resolverInfo) ResolvableDomains() []ids.Domain {
	return this.resolvableDomains
}

func (this *resolverInfo) Value(key interface{}) (res interface{}) {
	if this.values != nil {
		res = this.values[key]
	}

	if res == nil && this.parent != nil {
		res = this.parent.Value(key)
	}
	//fmt.Printf("value[%s]=%v\n", key, res)
	return res
}

func (this *resolverInfo) Matches(selector Selector) bool {
	if this.resolvableTypes != nil && selector.Type() != nil {
		for _, resolvableType := range this.resolvableTypes {
			if resolvableType.Equals(selector.Type()) {
				return true
			}
		}
	} else {
		return true
	}

	return false
}

func (this *resolverInfo) KeyExtractor() KeyExtractor {
	return this.keyExtractor
}

func (this *resolverInfo) DerivedCopy() ResolverInfo {
	var resolvableTypes []ids.TypeIdentifier
	var resolvableDomains []ids.Domain

	if this.resolvableTypes != nil {
		resolvableTypes = make([]ids.TypeIdentifier, len(this.ResolvableTypes()))
		copy(resolvableTypes, this.ResolvableTypes())
	}

	if resolvableDomains != nil {
		resolvableDomains = make([]ids.Domain, len(this.ResolvableDomains()))
		copy(resolvableDomains, this.ResolvableDomains())
	}

	return &resolverInfo{this.resolverType,
		resolvableTypes,
		resolvableDomains,
		this.keyExtractor,
		nil, this}
}

func (this *resolverInfo) WithExtractor(keyExtractor KeyExtractor) ResolverInfo {
	return &resolverInfo{this.resolverType,
		this.resolvableTypes,
		this.resolvableDomains,
		keyExtractor, nil, this}
}

func (this *resolverInfo) WithValue(key, value interface{}) ResolverInfo {
	values := make(map[interface{}]interface{})
	values[key] = value
	return &resolverInfo{this.resolverType,
		this.resolvableTypes,
		this.resolvableDomains,
		this.keyExtractor, values, this}
}

func (this *resolverInfo) WithValues(values map[interface{}]interface{}) ResolverInfo {
	return &resolverInfo{this.resolverType,
		this.resolvableTypes,
		this.resolvableDomains,
		this.keyExtractor,
		values, this}
}
