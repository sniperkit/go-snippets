package localresolver

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/mappings"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/types"
	"github.com/distributed-vision/go-resources/types/gotypeid"
	"github.com/distributed-vision/go-resources/types/publictypeid"
	"github.com/distributed-vision/go-resources/version"
)

var resolverType ids.TypeIdentifier = gotypeid.IdOf(reflect.TypeOf(LocalResolver{}))
var publicTypeVersion = version.SemanticVersion{Major: 0, Minor: 0, Patch: 1}

var PublicType = types.MustNewId(publictypeid.ResolverDomain, []byte("LocalResolver"), &publicTypeVersion)

func init() {
	mappings.Map(context.Background(), resolverType, PublicType)
	resolvers.ResisterNewFactoryFunction(PublicType, NewResolverFactory)
}

func NewResolverInfo(resolvableTypes []ids.TypeIdentifier, resolvableDomains []ids.Domain,
	keyExtractor resolvers.KeyExtractor, values map[interface{}]interface{}) resolvers.ResolverInfo {
	return resolvers.NewResolverInfo(PublicType,
		resolvableTypes, resolvableDomains, keyExtractor, values)
}

type factory struct {
	resolverInfo resolvers.ResolverInfo
}

func NewResolverFactory(resolverInfo resolvers.ResolverInfo) (resolvers.ResolverFactory, error) {
	return &factory{resolverInfo}, nil
}

func (this *factory) New(resolutionContext context.Context) (resolvers.Resolver, error) {
	return New(this.resolverInfo)
}

func (this *factory) ResolverType() ids.TypeIdentifier {
	return resolverType
}

func (this *factory) ResolverInfo() resolvers.ResolverInfo {
	return this.resolverInfo
}

type LocalResolver struct {
	resolverInfo resolvers.ResolverInfo
	entityMap    map[interface{}]interface{}
	mutex        *sync.Mutex
}

func New(baseInfo resolvers.ResolverInfo) (*LocalResolver, error) {
	if baseInfo == nil {
		return nil, fmt.Errorf("base resolver info must be defined")
	}

	return &LocalResolver{
		baseInfo.DerivedCopy(),
		make(map[interface{}]interface{}),
		&sync.Mutex{}}, nil
}

func (this *LocalResolver) ResolverInfo() resolvers.ResolverInfo {
	return this.resolverInfo
}

func (this *LocalResolver) Get(resolutionContext context.Context, selector resolvers.Selector) (interface{}, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var key string

	switch selector.Key().(type) {
	case string:
		key = selector.Key().(string)
	case []byte:
		key = string(selector.Key().([]byte))
	default:
		key = fmt.Sprintf("%v", selector.Key())
	}

	entity, ok := this.entityMap[key]

	//fmt.Printf("key=%v, entity: %v\n", key, entity)
	if ok {
		if selector.Test(entity) {
			return entity, nil
		}
	}

	for _, entity := range this.entityMap {
		if selector.Test(entity) {
			return entity, nil
		}
	}

	return nil, fmt.Errorf("Can't resolve entity for %v", selector)
}

func (this *LocalResolver) Resolve(resolutionContext context.Context, selector resolvers.Selector) (chan interface{}, chan error) {
	cres, cerr := make(chan interface{}), make(chan error)

	go func() {
		entity, err := this.Get(resolutionContext, selector)

		if err != nil {
			cerr <- err
		} else {
			cres <- entity
		}

		close(cres)
		close(cerr)
	}()

	return cres, cerr
}

func (this *LocalResolver) Put(resolutionContext context.Context, entity interface{}) error {
	keyExtractor := this.resolverInfo.KeyExtractor()

	if key, ok := keyExtractor(entity); ok {
		this.entityMap[key] = entity
	} else {
		return fmt.Errorf("Cannot extract key from: %v", entity)
	}

	return nil
}

func (this *LocalResolver) Post(resolutionContext context.Context, entity interface{}) error {
	keyExtractor := this.resolverInfo.KeyExtractor()

	if key, ok := keyExtractor(entity); ok {
		if _, ok := this.entityMap[key]; ok {
			this.entityMap[key] = entity
		} else {
			return fmt.Errorf("Can't resolve entity for %v", key)
		}
	} else {
		return fmt.Errorf("Cannot extract key from: %v", entity)
	}

	return nil
}

func (this *LocalResolver) Delete(resolutionContext context.Context, selector resolvers.Selector) error {
	var key string

	switch selector.Key().(type) {
	case string:
		key = selector.Key().(string)
	case []byte:
		key = string(selector.Key().([]byte))
	default:
		key = fmt.Sprintf("%v", selector.Key())
	}

	delete(this.entityMap, key)

	return nil
}
