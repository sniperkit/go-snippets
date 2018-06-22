package fileresolver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sync"

	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/ids/mappings"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/translators"
	"github.com/distributed-vision/go-resources/types"
	"github.com/distributed-vision/go-resources/types/gotypeid"
	"github.com/distributed-vision/go-resources/types/publictypeid"
	"github.com/distributed-vision/go-resources/util"
	"github.com/distributed-vision/go-resources/version"
)

var contentType ids.TypeIdentifier = gotypeid.IdOf(reflect.TypeOf(map[string]interface{}{}))

type fileResolver struct {
	path         string
	resolverInfo resolvers.ResolverInfo
	entityMap    map[interface{}]interface{}
}

var resolverMap map[string]*fileResolver = make(map[string]*fileResolver)
var resolverMapMutex = &sync.Mutex{}
var resolverType ids.TypeIdentifier = gotypeid.IdOf(reflect.TypeOf(fileResolver{}))
var publicTypeVersion = version.New(0, 0, 1)

var PublicType = types.MustNewId(publictypeid.ResolverDomain, []byte("FileResolver"), publicTypeVersion)

func init() {
	mappings.Map(context.Background(), resolverType, PublicType)
	resolvers.ResisterNewFactoryFunction(PublicType, NewResolverFactory)
}

func baseTypeTranslator(translationContext context.Context, resoverInfo interface{}) (chan interface{}, chan error) {
	return nil, nil
}

func ResolverType() ids.TypeIdentifier {
	return resolverType
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

	locationValue := this.resolverInfo.Value("location")

	if locationValue == nil {
		return nil, fmt.Errorf("ResolverInfo 'location' value can't be nil")
	}

	location := locationValue.(string)

	return New(
		location,
		this.resolverInfo)
}

func (this *factory) ResolverType() ids.TypeIdentifier {
	return resolverType
}

func (this *factory) ResolverInfo() resolvers.ResolverInfo {
	return this.resolverInfo
}

func New(locator string, resolverInfo resolvers.ResolverInfo) (resolvers.Resolver, error) {
	resolverMapMutex.Lock()
	defer resolverMapMutex.Unlock()
	resolver, ok := resolverMap[locator]

	if ok {
		return resolver, nil
	}

	resolver, err := newResolver(locator, resolverInfo)

	if err != nil {
		return nil, err
	}

	resolverMap[locator] = resolver

	return resolver, nil
}

func newResolver(file string, resolverInfo resolvers.ResolverInfo) (*fileResolver, error) {

	var filePath string

	pathsValue := resolverInfo.Value("paths")

	if pathsValue != nil {

		for _, respath := range pathsValue.([]string) {
			respath = path.Join(respath, file)

			_, err := os.Stat(respath)

			if err == nil {
				filePath = respath
				break
			}
		}
	} else {
		_, err := os.Stat(file)

		if err == nil {
			filePath = file
		}
	}

	if filePath == "" {
		if pathsValue == nil {
			return nil, fmt.Errorf("Can't find: %s", file)
		}

		return nil, fmt.Errorf("Can't find: %s in: %v", file, pathsValue)
	}

	return &fileResolver{path: filePath, resolverInfo: resolverInfo}, nil
}

func (this *fileResolver) ResolverInfo() resolvers.ResolverInfo {
	return this.resolverInfo
}

func (this *fileResolver) Get(resolutionContext context.Context, selector resolvers.Selector) (entity interface{}, err error) {
	return util.Await(this.Resolve(resolutionContext, selector))
}

func (this *fileResolver) Resolve(resolutionContext context.Context, selector resolvers.Selector) (chan interface{}, chan error) {
	cres, cerr := make(chan interface{}), make(chan error)

	if this.resolverInfo != nil {
		resolutionContext = context.WithValue(resolutionContext, "resolverInfo", this.resolverInfo)
	}

	go func() {
		entityMap, err := this.getMap(resolutionContext, selector.Type())

		if err != nil {
			cerr <- err
		} else if entityMap != nil {
			//fmt.Printf("sel=%+v\n", selector)
			//fmt.Printf("entities: %v\n", entityMap)
			var key string

			switch selector.Key().(type) {
			case string:
				key = selector.Key().(string)
			case []byte:
				key = string(selector.Key().([]byte))
			default:
				key = fmt.Sprintf("%v", selector.Key())
			}

			entity, ok := entityMap[key]

			//fmt.Printf("key=%v, entity: %v\n", key, entity)
			if ok {
				if selector.Test(entity) {
					cres <- entity
					close(cres)
					close(cerr)
					return
				}
			}

			for _, entity := range entityMap {
				if selector.Test(entity) {
					cres <- entity
					close(cres)
					close(cerr)
					return
				}
			}

			cerr <- fmt.Errorf("Invalid entity selector: %+v", selector)
		}

		close(cres)
		close(cerr)
	}()

	return cres, cerr
}

func (this *fileResolver) getMap(context context.Context, targetType ids.TypeIdentifier) (map[interface{}]interface{}, error) {
	if this.entityMap != nil {
		return this.entityMap, nil
	}

	entityMap, err := this.loadMap(context, targetType)

	if err == nil {
		this.entityMap = entityMap
	}

	return entityMap, err
}

var untypedLocalDomain []byte = domain.MustDecodeId(encodertype.BASE62, "3", "")

func (this *fileResolver) loadMap(context context.Context, targetType ids.TypeIdentifier) (map[interface{}]interface{}, error) {
	data, err := ioutil.ReadFile(this.path)

	if err != nil {
		return nil, err
	}

	var jsonEntities interface{}
	err = json.Unmarshal(data, &jsonEntities)

	if err != nil {
		return nil, err
	}

	jsonEntityMap := jsonEntities.(map[string]interface{})

	//fmt.Printf("Content Type=%+v, Target Type %+v\n", contentType, targetType)

	keyExtractor := this.resolverInfo.KeyExtractor()
	entityMap := make(map[interface{}]interface{})

	for id, jsonEntity := range jsonEntityMap {
		if contentType != targetType {
			//fmt.Printf("id=%s\n", id)
			entityId, err := identifier.New(untypedLocalDomain, []byte(id), nil)
			//fmt.Printf("id=%v\n", entityId.Value())

			if err != nil {
				return nil, err
			}

			if entity, err := util.Await(
				translators.Translate(context, contentType,
					entityId, jsonEntity, targetType)); err == nil {
				if key, ok := keyExtractor(entity); ok {
					//fmt.Printf("key(ext)=%v\n", key)
					entityMap[key] = entity
				} else {
					return nil, fmt.Errorf("Can't extract key from: %v", entity)
				}
			} else {
				return nil, err
			}
		} else {
			if key, ok := keyExtractor(jsonEntity); ok {
				entityMap[key] = jsonEntity
			}
		}
	}

	return entityMap, nil
}
