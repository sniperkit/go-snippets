package signaturedomain

import (
	"context"
	"fmt"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domaintype"
	"github.com/distributed-vision/go-resources/ids/identitydomain"
	"github.com/distributed-vision/go-resources/ids/scheme"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

func Init() {
}

func init() {
	domain.RegisterJSONUnmarshaller(domaintype.SIGNATURE, unmarshalJSON)
}

type signatureDomain struct {
	ids.IdentityDomain
}

func unmarshalJSON(unmarshalContext context.Context, json map[string]interface{}) (ids.Domain, error) {
	rootId, err := base62.Decode(json["id"].(string))
	resolverInfo, hasResolverInfo := unmarshalContext.Value("resolverInfo").(resolvers.ResolverInfo)

	if err != nil {
		return nil, err
	}

	var incarnation *uint32
	var crcLength uint
	var versionType = versiontype.UNVERSIONED

	info := make(map[interface{}]interface{})

	for key, value := range json {
		info[key] = value
	}

	var idScheme ids.Scheme

	if hasResolverInfo && resolverInfo.Value("schemeId") != nil {
		//fmt.Printf("schemeid=%v\n", resolverInfo.Value("schemeId"))
		idScheme, err = scheme.Get(unmarshalContext, scheme.Selector{Id: resolverInfo.Value("schemeId").([]byte)})
		//fmt.Printf("scheme=%v\n", scheme)
		if err != nil {
			return nil, err
		}
	}

	if idScheme == nil {
		return nil, fmt.Errorf("Can't unmarshal domain for unknow scheme")
	}

	base, err := identitydomain.New(idScheme, rootId, incarnation, crcLength, versionType, false, false, info)
	//fmt.Printf("base=%v\n", base)

	if err != nil {
		return nil, err
	}

	return &signatureDomain{base}, nil
}

func NewElements() (ids.SignatureElements, error) {
	return nil, nil
}
