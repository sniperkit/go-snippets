package domain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domaintype"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/translators"
	"github.com/distributed-vision/go-resources/util"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

var domainResolverInfo resolvers.ResolverInfo
var domainResolver *resolvers.CompositeResolver

var domainEntityType ids.TypeIdentifier
var schemeEntityType ids.TypeIdentifier

var PublicResolverType ids.TypeIdentifier

func init() {
	ids.OnLocalTypeInit(func() {
		var err error

		if domainEntityType == nil {
			domainEntityType = ids.NewLocalTypeId(reflect.TypeOf((*ids.Domain)(nil)).Elem())
		}

		if schemeEntityType == nil {
			schemeEntityType = ids.NewLocalTypeId(reflect.TypeOf((*ids.Scheme)(nil)).Elem())
		}

		mapType := ids.NewLocalTypeId(reflect.TypeOf(map[string]interface{}{}))
		translators.Register(context.Background(), mapType, domainEntityType, domainMapTranslator)

		PublicResolverType, err = ids.NewTypeId(
			MustDecodeId(encodertype.BASE62, "T", "0", uint32(0), uint(0), versiontype.SEMANTIC),
			[]byte("DomainResolver"), version.New(0, 0, 1))

		domainResolverInfo = resolvers.NewResolverInfo(PublicResolverType,
			[]ids.TypeIdentifier{domainEntityType}, nil, KeyExtractor, nil)
		domainResolver, err = resolvers.NewCompositeResolver(domainResolverInfo)

		if err != nil {
			panic(fmt.Sprint("Domain resolver creation failed with:", err))
		}

		resolvers.RegisterResolver(domainResolver)
	})
}

func domainMapTranslator(translationContext context.Context, fromId ids.Identifier, fromValue interface{}) (chan interface{}, chan error) {

	cres := make(chan interface{}, 1)
	cerr := make(chan error, 1)

	json := fromValue.(map[string]interface{})
	json["id"] = string(fromId.Id())

	toValue, err := unmarshalJSON(translationContext, json)
	//fmt.Printf("id: %+v val: %+v err: %s\n", fromId.Id(), toValue, err)

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
	SchemeId []byte
	IdRoot   []byte
	Id       []byte
	Name     string
	Opts     SelectorOpts
}

func (this *Selector) Test(candidate interface{}) bool {

	domain, ok := candidate.(ids.Domain)

	if !ok {
		return false
	}

	if this.IdRoot != nil && !bytes.Equal(this.IdRoot, domain.IdRoot()) {
		return false
	}

	if this.SchemeId != nil && !bytes.Equal(this.SchemeId, domain.SchemeId()) {
		return false
	}

	if this.Id != nil && !bytes.Equal(this.Id, domain.Id()) {
		return false
	}

	if this.Name != "" {
		if this.Opts.IgnoreCase {
			if strings.ToUpper(this.Name) != strings.ToUpper(domain.Name()) {
				return false
			}
		} else {
			if this.Name != domain.Name() {
				return false
			}
		}
	}

	return true
}

func (this *Selector) Key() interface{} {
	if this.Id != nil {
		return base62.Encode(this.Id)
	}

	id, err := ToId(this.SchemeId, this.IdRoot, nil, 0, versiontype.UNVERSIONED, false, false)

	if err != nil {
		return nil
	}

	return base62.Encode(id)
}

var entityType ids.TypeIdentifier

func (this *Selector) Type() ids.TypeIdentifier {
	if entityType == nil {
		entityType = ids.NewLocalTypeId(reflect.TypeOf((*ids.Domain)(nil)).Elem())
	}

	return entityType
}

type schemeSelector struct {
	id []byte
}

func (this *schemeSelector) Test(candidate interface{}) bool {
	scheme, ok := candidate.(ids.Scheme)

	if !ok {
		return false
	}

	if this.id != nil && !bytes.Equal(this.id, scheme.Id()) {
		return false
	}

	return true
}

func (this *schemeSelector) Key() interface{} {
	return base62.Encode(this.id)
}

func (this *schemeSelector) Type() ids.TypeIdentifier {
	return schemeEntityType
}

func RegisterResolverFactory(resolverFactory resolvers.ResolverFactory) error {
	return domainResolver.RegisterComponentFactory(resolverFactory, false)
}

func Get(resolutionContext context.Context, selector Selector) (domain ids.Domain, err error) {
	return Await(Resolve(resolutionContext, selector))
}

var resolveMutex = &sync.Mutex{}

func Resolve(resolutionContext context.Context, selector Selector) (chan ids.Domain, chan error) {

	cResOut := make(chan ids.Domain, 1)
	cErrOut := make(chan error, 1)

	schemeId := selector.SchemeId

	if schemeId == nil && selector.Id != nil {
		schemeId = Wrap(selector.Id).SchemeId()
	}

	go func() {
		// this forces the load of any demain definition resolvers associated with this scheme
		if schemeId != nil {
			_, err := util.Await(resolvers.Resolve(resolutionContext, &schemeSelector{id: schemeId}))

			if err != nil {
				cErrOut <- err
				close(cResOut)
				close(cErrOut)
				return
			}
		}

		res, err := util.Await(domainResolver.Resolve(resolutionContext, &selector))

		if err == nil {
			if domain, ok := res.(ids.Domain); ok {
				cResOut <- domain
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

type unmarshaller func(unmarshalContext context.Context, json map[string]interface{}) (ids.Domain, error)

var unmarshalers map[ids.DomainType]unmarshaller = make(map[ids.DomainType]unmarshaller)

func RegisterJSONUnmarshaller(domainType ids.DomainType, unmarshaller unmarshaller) {
	unmarshalers[domainType] = unmarshaller
}

func unmarshalJSON(unmarshalContext context.Context, json map[string]interface{}) (ids.Domain, error) {
	dt, err := domaintype.Parse(json["domainType"].(string))
	//fmt.Printf("dt=%v\n", dt)
	if err != nil {
		return nil, err
	}

	unmarshaler, ok := unmarshalers[dt]

	if !ok {
		return nil, errors.New("Unknown domain type: " + json["domainType"].(string))
	}

	return unmarshaler(unmarshalContext, json)
}
