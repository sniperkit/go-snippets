package scheme

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/distributed-vision/go-resources/encoding"
	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/schemeformat"
	"github.com/distributed-vision/go-resources/ids/schemevisibility"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

// DomainPathKey is the unique context path for the file paths to seach for domain
//  .json files
var DomainPathKey = string(domain.MustDecodeId(encodertype.BASE62, "2", "1")) + "-DOMAINPATH"

type scheme struct {
	ids.Domain
	visibility ids.SchemeVisibility
	format     ids.SchemeFormat
}

func KeyExtractor(entity ...interface{}) (interface{}, bool) {
	if len(entity) > 0 {
		if domain, ok := entity[0].(ids.Scheme); ok {
			return base62.Encode(domain.Id()), true
		}
	}
	return nil, false
}

func Await(cres chan ids.Scheme, cerr chan error) (result ids.Scheme, err error) {
	if cres == nil || cerr == nil {
		return nil, fmt.Errorf("Await Failed: channels are undefined")
	}

	resolved := false
	for !resolved {
		select {
		case res, ok := <-cres:
			if ok {
				result = res
				resolved = true
			}
		case error, ok := <-cerr:
			if ok {
				err = error
				resolved = true
			}
		}
	}

	return result, err
}

func unmarshalJSON(unmarshalContext context.Context, json map[string]interface{}) (interface{}, error) {
	id, err := base62.Decode(json["id"].(string))

	if err != nil {
		return nil, err
	}

	if len(id) == 0 {
		return nil, fmt.Errorf("Empty id is invalid")
	}

	visibility, err := schemevisibility.Parse(json["visibility"].(string))

	if err != nil {
		return nil, err
	}

	format, err := schemeformat.Parse(json["format"].(string))

	if err != nil {
		return nil, err
	}

	/*
		if unmarshalContext.Value("paths")!=nil {
			unmarshalContext=context.WithValue(unmarshalContext, domain.DomainPathKey, path )
		}
	*/

	return NewScheme(
		id,
		toString(json, "name"),
		toString(json, "description"),
		visibility,
		format,
		toDomainInfo(unmarshalContext, id, json["domainInfo"].(map[string]interface{})))
}

func NewScheme(id []byte, name string, description string, visibility ids.SchemeVisibility, format ids.SchemeFormat, info map[interface{}]interface{}) (ids.Scheme, error) {

	info["name"] = name
	info["description"] = description

	base, err := domain.New(nil, id, nil, 0, versiontype.UNVERSIONED, false, false, info)

	//fmt.Printf("base id=%v, scheme=%v\n", base.Id(), base.SchemeId())
	if err != nil {
		return nil, err
	}

	return &scheme{
		Domain:     base,
		visibility: visibility,
		format:     format}, nil
}

var empty = []byte{}
var extensionBit byte = (1 << 6)

func ToId(base byte, extension []byte) ([]byte, error) {

	if base > 61 {
		return nil, errors.New("base too Long: scheme id base must be < 61")
	}

	if extension == nil {
		extension = empty
	}

	var extensionLen = len(extension)

	if extensionLen > 255 {
		return nil, errors.New("Id too Long: scheme id extension binary length must be < 255")
	}

	var baseSlice []byte

	if extensionLen == 0 {
		baseSlice = []byte{base}
	} else {
		baseSlice = []byte{base | extensionBit, byte(extensionLen & 0xff)}
	}

	return bytes.Join([][]byte{baseSlice, extension}, empty), nil
}

func DecodeId(encoderType encodertype.EncoderType, base string, extension string) ([]byte, error) {
	var baseValue []byte
	var extensionValue []byte
	var err error

	baseValue, err = encoding.Decode(base, encoderType)

	if err != nil {
		return nil, fmt.Errorf("Invalid base encoding %s", err)
	}

	if len(baseValue) != 1 {
		return nil, fmt.Errorf("Invalid base length: %v", len(baseValue))
	}

	extensionValue, err = encoding.Decode(extension, encoderType)

	if err != nil {
		return nil, fmt.Errorf("extension rootId encoding %s", err)
	}

	return ToId(baseValue[0], extensionValue)
}

func MustDecodeId(encoderType encodertype.EncoderType, base string, extension string) []byte {
	id, err := DecodeId(encoderType, base, extension)

	if err != nil {
		panic(fmt.Sprintf("Failed to encode id: %s", err))
	}

	return id
}

func (this *scheme) Visibility() ids.SchemeVisibility {
	return this.visibility
}

func (this *scheme) Format() ids.SchemeFormat {
	return this.format
}

func (this *scheme) RegisterResolvers() error {
	resolverInfos := this.InfoValue("resolverInfo").([]resolvers.ResolverInfo)
	var errs = make([]error, 0)

	if resolverInfos != nil {

		for _, resolverInfo := range resolverInfos {
			factory, err := resolvers.NewResolverFactory(resolverInfo)

			if err == nil {
				domain.RegisterResolverFactory(factory)
			} else {
				//fmt.Printf("factory err=%s\n", err)
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Resolver registration fails: %v", errs)
	}

	return nil
}

func toString(json map[string]interface{}, field string) string {
	value := json[field]

	if value == nil {
		return ""
	}

	return value.(string)
}

var fileResolverDomainId = domain.MustDecodeId(encodertype.BASE62, "T", "0", uint32(0), uint(0), versiontype.SEMANTIC)

func toDomainInfo(infoContext context.Context, schemeId []byte, infoIn map[string]interface{}) map[interface{}]interface{} {

	infoOut := make(map[interface{}]interface{})

	var sourcePaths []string

	if contextResInfo, ok := infoContext.Value("resolverInfo").(resolvers.ResolverInfo); ok {
		if paths, ok := contextResInfo.Value("paths").([]string); ok {
			sourcePaths = paths
		}
	}

	for key, value := range infoIn {

		if strings.ToUpper(key) == strings.ToUpper("resolverInfo") {
			resolverInfos, ok := value.([]interface{})
			if ok {
				resolverInfosOut := make([]resolvers.ResolverInfo, len(resolverInfos))

				for index, resolverInfo := range resolverInfos {
					var infoMap = map[interface{}]interface{}{}
					var resolverType ids.TypeIdentifier

					for key, value := range resolverInfo.(map[string]interface{}) {

						switch key {
						case "resolverType":
							resolverTypeParts := strings.Split(value.(string), "-")
							idPart := []byte(resolverTypeParts[0])

							if len(resolverTypeParts) > 1 {
								versionPart, err := version.Parse(resolverTypeParts[1])

								if err == nil {
									resolverType, err = ids.NewTypeId(fileResolverDomainId, idPart, versionPart)
								}
							} else {
								var err error
								resolverType, err = ids.NewTypeId(fileResolverDomainId, idPart, nil)

								if err != nil {
								}
							}
							break

						default:
							infoMap[key] = value
						}
					}

					if resolverType != nil {

						if sourcePaths != nil && "FileResolver" == string(resolverType.Id()) {
							if /*value*/ _, ok := infoMap["paths"]; !ok {
								infoMap["paths"] = sourcePaths
							} else {
								// TODO append source paths after predefined paths
							}
						}

						infoMap["schemeId"] = schemeId
						resolverInfosOut[index] = resolvers.NewResolverInfo(resolverType,
							[]ids.TypeIdentifier{domainEntityType}, []ids.Domain{},
							domain.KeyExtractor, infoMap)
					}
				}

				infoOut["resolverInfo"] = resolverInfosOut
			}
		} else {
			infoOut[key] = value
		}
	}

	return infoOut
}
