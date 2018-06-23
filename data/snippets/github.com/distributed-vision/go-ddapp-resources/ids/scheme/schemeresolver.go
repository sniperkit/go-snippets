package scheme

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/translators"
	"github.com/distributed-vision/go-resources/util"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

var schemeResolverInfo resolvers.ResolverInfo
var schemeResolver *resolver

var schemeEntityType ids.TypeIdentifier
var domainEntityType ids.TypeIdentifier

var PublicResolverType ids.TypeIdentifier

func init() {
	ids.OnLocalTypeInit(func() {
		var err error

		if schemeEntityType == nil {
			schemeEntityType = ids.NewLocalTypeId(reflect.TypeOf((*ids.Scheme)(nil)).Elem())
		}

		if domainEntityType == nil {
			domainEntityType = ids.NewLocalTypeId(reflect.TypeOf((*ids.Domain)(nil)).Elem())
		}

		mapType := ids.NewLocalTypeId(reflect.TypeOf(map[string]interface{}{}))
		translators.Register(context.Background(), mapType, schemeEntityType, schemeMapTranslator)

		PublicResolverType, err = ids.NewTypeId(
			domain.MustDecodeId(encodertype.BASE62, "T", "0", uint32(0), uint(0), versiontype.SEMANTIC),
			[]byte("SchemeResolver"), version.New(0, 0, 1))

		schemeResolverInfo = resolvers.NewResolverInfo(PublicResolverType,
			[]ids.TypeIdentifier{schemeEntityType}, nil, KeyExtractor, nil)
		baseResolver, err := resolvers.NewCompositeResolver(schemeResolverInfo)
		schemeResolver = &resolver{baseResolver}

		if err != nil {
			panic(fmt.Sprint("Scheme resolver creation failed with:", err))
		}

		resolvers.RegisterResolver(schemeResolver)
	})
}

var untypedDomain []byte = domain.MustDecodeId(encodertype.BASE62, "3", "")

func schemeMapTranslator(translationContext context.Context, fromId ids.Identifier, fromValue interface{}) (chan interface{}, chan error) {

	cres := make(chan interface{}, 1)
	cerr := make(chan error, 1)

	json := fromValue.(map[string]interface{})
	json["id"] = string(fromId.Id())

	toValue, err := unmarshalJSON(translationContext, json)
	//fmt.Printf("val: %+v err: %s\n", toValue, err)

	if err != nil {
		cerr <- err
	} else {
		cres <- toValue
	}

	close(cres)
	close(cerr)

	return cres, cerr
}

type SelectorOpts struct {
	IgnoreCase       bool
	IgnoreWhitespace bool
}

type Selector struct {
	Id   []byte
	Name string
	Opts SelectorOpts
}

func (this *Selector) Test(candidate interface{}) bool {
	scheme, ok := candidate.(ids.Scheme)

	if !ok {
		return false
	}

	if this.Id != nil && !bytes.Equal(this.Id, scheme.Id()) {
		return false
	}

	if this.Name != "" {
		if this.Opts.IgnoreCase {
			if strings.ToUpper(this.Name) != strings.ToUpper(scheme.Name()) {
				return false
			}
		} else {
			if this.Name != scheme.Name() {
				return false
			}
		}
	}

	return true
}

func (this *Selector) Key() interface{} {
	return base62.Encode(this.Id)
}

func (this *Selector) Type() ids.TypeIdentifier {
	return schemeEntityType
}

type resolver struct {
	*resolvers.CompositeResolver
}

func (this *resolver) Get(resolutionContext context.Context, selector resolvers.Selector) (interface{}, error) {
	return util.Await(this.Resolve(resolutionContext, selector))
}

func (this *resolver) Resolve(resolutionContext context.Context, selector resolvers.Selector) (chan interface{}, chan error) {
	cResOut := make(chan interface{}, 1)
	cErrOut := make(chan error, 1)

	go func() {
		res, err := util.Await(this.CompositeResolver.Resolve(resolutionContext, selector))

		if err == nil {
			if scheme, ok := res.(ids.Scheme); ok {
				scheme.RegisterResolvers()
				cResOut <- scheme
			} else {
				cErrOut <- fmt.Errorf("Resolver returned invalid type, expected: ids.Domain got: %s", reflect.TypeOf(res))
			}
		} else {
			cErrOut <- err
		}

		close(cResOut)
		close(cErrOut)
	}()

	return cResOut, cErrOut
}

func RegisterResolverFactory(resolverFactory resolvers.ResolverFactory) error {
	return schemeResolver.RegisterComponentFactory(resolverFactory, false)
}

func Get(resolutionContext context.Context, selector Selector) (ids.Scheme, error) {
	res, err := schemeResolver.Get(resolutionContext, &selector)

	if err != nil {
		return nil, err
	}

	if scheme, ok := res.(ids.Scheme); ok {
		return scheme, err
	}

	return nil, fmt.Errorf("Resolver returned invalid type, expected: ids.Domain got: %s", reflect.TypeOf(res))
}

func Resolve(resolutionContext context.Context, selector Selector) (chan ids.Scheme, chan error) {

	cResOut := make(chan ids.Scheme, 1)
	cErrOut := make(chan error, 1)

	go func() {
		res, err := util.Await(schemeResolver.Resolve(resolutionContext, &selector))

		if err == nil {
			if scheme, ok := res.(ids.Scheme); ok {
				cResOut <- scheme
			} else {
				cErrOut <- fmt.Errorf("Resolver returned invalid type, expected: ids.Domain got: %s", reflect.TypeOf(res))
			}
		} else {
			cErrOut <- err
		}

		close(cResOut)
		close(cErrOut)
	}()

	return cResOut, cErrOut
}
