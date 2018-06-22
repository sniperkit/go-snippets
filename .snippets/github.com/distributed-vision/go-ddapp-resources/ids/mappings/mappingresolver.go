package mappings

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/translators"
	"github.com/distributed-vision/go-resources/translators/maptranslator"
	"github.com/distributed-vision/go-resources/util"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

var mappingResolverInfo resolvers.ResolverInfo
var mappingResolver *resolvers.CompositeResolver

var mappingsType = reflect.TypeOf(Mappings{})
var mappingsEntityType ids.TypeIdentifier
var PublicResolverType ids.TypeIdentifier

func init() {
	ids.OnLocalTypeInit(func() {
		var err error

		if mappingsEntityType == nil {
			mappingsEntityType = ids.NewLocalTypeId(mappingsType)
		}

		mapType := ids.NewLocalTypeId(reflect.TypeOf(map[string]interface{}{}))
		translators.Register(context.Background(), mapType, mappingsEntityType, mappingsMapTranslator)

		PublicResolverType, err = ids.NewTypeId(
			domain.MustDecodeId(encodertype.BASE62, "T", "0", uint32(0), uint(0), versiontype.SEMANTIC),
			[]byte("MappingsResolver"), version.New(0, 0, 1))

		mappingResolverInfo = resolvers.NewResolverInfo(PublicResolverType,
			[]ids.TypeIdentifier{mappingsEntityType}, nil, KeyExtractor, nil)
		mappingResolver, err = resolvers.NewCompositeResolver(mappingResolverInfo)

		if err != nil {
			panic(fmt.Sprint("Scheme resolver creation failed with:", err))
		}

		resolvers.RegisterResolver(mappingResolver)
	})
}

func mappingsMapTranslator(translationContext context.Context, fromId ids.Identifier, fromValue interface{}) (chan interface{}, chan error) {
	cres := make(chan interface{}, 1)
	cerr := make(chan error, 1)

	json := fromValue.(map[string]interface{})

	toValue, err := maptranslator.FromMap(json, mappingsType)

	if err != nil {
		cerr <- err
	} else {
		cres <- toValue
	}

	close(cres)
	close(cerr)

	return cres, cerr
}

type Selector struct {
	From ids.Identifier
	To   ids.IdentityDomain
	At   time.Time
}

func (this *Selector) Type() ids.TypeIdentifier {
	return mappingsEntityType
}

func (this *Selector) Test(candidate interface{}) bool {
	mapping, ok := candidate.(Mappings)

	if ok {
		for _, mappedid := range mapping.mappedIds {
			if this.At.After(mappedid.from) && this.At.Before(mappedid.to) &&
				this.From.Equals(mapping.fromId) && this.To.Equals(mapping.toDomain) {
				return true
			}
		}
	}

	return false
}

func (this *Selector) Key() interface{} {
	return this.From.String() + "->" + this.To.String()
}

type MappingResolver interface {
	resolvers.Resolver
	Map(mappingContext context.Context, from ids.Identifier, to ids.Identifier, between ...time.Time) chan error
}

func RegisterResolver(resolver resolvers.Resolver) error {
	return mappingResolver.RegisterComponent(resolver)
}

func RegisterResolverFactory(resolverFactory resolvers.ResolverFactory) error {
	return mappingResolver.RegisterComponentFactory(resolverFactory, false)
}

func Get(resolutionContext context.Context, selector Selector) (domain ids.Mapping, err error) {
	return AwaitMapping(Resolve(resolutionContext, selector))
}

func Resolve(resolutionContext context.Context, selector Selector) (chan ids.Mapping, chan error) {
	cResOut := make(chan ids.Mapping, 1)
	cErrOut := make(chan error, 1)

	go func() {
		res, err := util.Await(mappingResolver.Resolve(resolutionContext, &selector))

		if err == nil {
			if mappings, ok := res.(Mappings); ok {
				mappingIndex := -1
				for index, mappedid := range mappings.mappedIds {
					if selector.At.After(mappedid.from) && selector.At.Before(mappedid.to) {
						cResOut <- &mapping{&mappings, index}
						mappingIndex = index
					}
				}

				if mappingIndex < 0 {
					cErrOut <- fmt.Errorf("Can't find mapping for: %s at: %v", selector.Key(), selector.At)
				}
			} else {
				cErrOut <- fmt.Errorf("Resolver returned invalid type, expected: mappings.Mappings got: %s", reflect.TypeOf(res))
			}
		} else {
			cErrOut <- err
		}

		close(cResOut)
		close(cErrOut)
	}()

	return cResOut, cErrOut
}

func Map(mappingContext context.Context, from ids.Identifier, to ids.Identifier, between ...time.Time) chan error {
	cErrOut := make(chan error, 1)

	go func() {
		defer close(cErrOut)

		selector := Selector{From: from, To: to.Domain()}
		mutableResolvers := mappingResolver.GetMutableComponents(mappingContext, &selector)

		if len(mutableResolvers) == 0 {
			cErrOut <- fmt.Errorf("No mutable mapping resolvers installed for: %s", from.Domain())
			return
		}

		var mappingResolver = mutableResolvers[0]

		if len(mutableResolvers) > 1 {

		}

		switch mappingResolver.(type) {
		case MappingResolver:
			err := util.AwaitError(mappingResolver.(MappingResolver).Map(mappingContext, from, to, between...))
			if err != nil {
				cErrOut <- err
			}
		case resolvers.MutableResolver:
			result, err := mappingResolver.Get(mappingContext, &selector)

			after, before := MinTime, MaxTime

			if len(between) > 0 {
				after = between[0]
			}

			if len(between) > 1 {
				after = between[1]
			}

			if result != nil {
				mappings := result.(Mappings)
				var currentMids = mappings.mappedIds
				var updatedMids []mappedId

				for index, mid := range currentMids {
					if mid.from.After(after) {
						updatedMids = append(currentMids, mappedId{})
						copy(updatedMids[index+1:], updatedMids[index:])
						updatedMids[index] = mappedId{after, before, to}
						break
					} else if mid.from.Equal(after) {
						if mid.to.Equal(after) {
							currentMids[index] = mappedId{after, before, to}
							updatedMids = currentMids
						} else {
							insertIndex := index
							if mid.to.Before(before) {
								insertIndex++
							}
							updatedMids = append(currentMids, mappedId{})
							copy(updatedMids[insertIndex+1:], updatedMids[insertIndex:])
							updatedMids[insertIndex] = mappedId{after, before, to}
						}
						break
					}
				}
				if updatedMids != nil {
					mappings.mappedIds = updatedMids
				} else {
					mappings.mappedIds = append(currentMids, mappedId{after, before, to})
				}
				mappingResolver.(resolvers.MutableResolver).Post(mappingContext, mappings)
			} else if _, ok := err.(*resolvers.EntityNotFound); ok {
				mappingResolver.(resolvers.MutableResolver).Put(mappingContext,
					Mappings{from, to.Domain(), []mappedId{mappedId{after, before, to}}})
			} else {
				cErrOut <- err
			}
		}

		close(cErrOut)
	}()

	return cErrOut
}
