package identifier

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/distributed-vision/go-resources/encoding"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/mappings"
	"github.com/distributed-vision/go-resources/ids/scheme"
	"github.com/distributed-vision/go-resources/util/hton"
	"github.com/distributed-vision/go-resources/util/ntoh"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versiontype"
	"github.com/howeyc/crc16"
	"github.com/sigurn/crc8"
)

func Init() {
}

func Await(cres chan ids.Identifier, cerr chan error) (result ids.Identifier, err error) {
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

type identifier struct {
	value []byte
}

func New(domainValue interface{}, id []byte, optionalValues ...interface{}) (ids.Identifier, error) {

	var crcLength uint
	var domainId []byte
	var versionType versiontype.VersionType
	var allowPaths bool
	var allowFragments bool
	var err error

	if domainValue == nil {
		return nil, errors.New("Invalid domain: undefined")
	}

	switch t := domainValue.(type) {
	case ids.Domain:
		dom := domainValue.(ids.Domain)
		crcLength = dom.CrcLength()
		domainId = dom.Id()
		versionType = dom.VersionType()
		allowPaths = dom.HasPaths()
		allowFragments = dom.HasFragments()
	case []byte:
		domainId = domainValue.([]byte)
		crcLength, err = identifierCrcLength(domainId)
		versionType, err = identifierVersionType(domainId)
		allowPaths = domain.HasPaths(domainId)
		allowFragments = domain.HasFragments(domainId)

		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Invalid domain type: %s", t)
	}

	if id == nil {
		return nil, errors.New("Invalid id: undefined")
	}

	if domainId == nil {
		return nil, errors.New("Invalid domain: id undefined")
	}

	var value []byte
	var idVersion version.Version
	var pathValue []byte
	var fragmentValue []byte

	if len(optionalValues) > 0 && optionalValues[0] != nil {
		switch optType := optionalValues[0].(type) {
		case version.Version:
			idVersion = optionalValues[0].(version.Version)
		case []byte:
			if pathValue == nil && allowPaths {
				pathValue = optionalValues[0].([]byte)
			} else if allowFragments {
				fragmentValue = optionalValues[0].([]byte)
			} else {
				return nil, fmt.Errorf("Domain can't accept paths or fragments")
			}
		case string:
			pathValue = []byte(optionalValues[0].(string))
		default:
			return nil, fmt.Errorf("Invalid id optional value 0: unexpected type: %v", optType)
		}
	}

	if len(optionalValues) > 1 && optionalValues[1] != nil {
		switch optType := optionalValues[1].(type) {
		case []byte:
			if pathValue == nil && allowPaths {
				pathValue = optionalValues[1].([]byte)
			} else if allowFragments {
				fragmentValue = optionalValues[1].([]byte)
			} else {
				if allowPaths {
					return nil, fmt.Errorf("Domain can't accept fragments")
				}
				return nil, fmt.Errorf("Domain can't accept paths or fragments")
			}
		case string:
			if pathValue == nil && allowPaths {
				pathValue = []byte(optionalValues[1].(string))
			} else if allowFragments {
				fragmentValue = []byte(optionalValues[1].(string))
			} else {
				if allowPaths {
					return nil, fmt.Errorf("Domain can't accept fragments")
				}
				return nil, fmt.Errorf("Domain can't accept paths or fragments")
			}
		default:
			return nil, fmt.Errorf("Invalid id optional value 1: unexpected type: %v", optType)
		}
	}

	if len(optionalValues) > 2 && optionalValues[2] != nil {
		if fragmentValue == nil && allowFragments {
			switch optType := optionalValues[2].(type) {
			case []byte:
				fragmentValue = optionalValues[2].([]byte)
			case string:
				fragmentValue = []byte(optionalValues[2].(string))
			default:
				return nil, fmt.Errorf("Invalid id optional value 2: unexpected type: %v", optType)
			}
		} else {
			return nil, fmt.Errorf("Domain can't accept fragments")
		}
	}

	pathLength := []byte{}

	if len(pathValue) > 0 || allowPaths {
		if !allowPaths {
			return nil, fmt.Errorf("Domain can't accept path value: %v", pathValue)
		}

		if len(pathValue) > 255 {
			return nil, errors.New("Path too Long: (path+fragment) must be < 256")
		}

		pathLength = []byte{byte(len(pathValue))}
	}

	fragmentLength := []byte{}

	if len(fragmentValue) > 0 || allowFragments {
		if !allowFragments {
			return nil, fmt.Errorf("Domain can't accept fragment value")
		}

		if len(fragmentValue) > 255 {
			return nil, errors.New("Fragment too Long: must be < 256")
		}

		fragmentLength = []byte{byte(len(fragmentValue))}
	}

	id = bytes.Join([][]byte{pathLength, fragmentLength, id, pathValue, fragmentValue}, []byte{})
	//fmt.Printf("fl=&%v,pl=%v, id=%v\n", fragmentLength, pathLength, id)
	switch versionType {
	case versiontype.UNVERSIONED:
		value = bytes.Join([][]byte{domainId, id}, []byte{})
		break
	case versiontype.NUMERIC:
		if nver, ok := idVersion.(version.NumericVersion); ok {
			value = bytes.Join([][]byte{domainId, []byte{nver.ByteLength()}, id, nver.Bytes()}, []byte{})
		} else {
			return nil, errors.New("Expected numeric version")
		}
		break
	case versiontype.SEMANTIC:
		if sver, ok := idVersion.(*version.SemanticVersion); ok {
			verbytes := sver.Bytes()
			blen := len(verbytes)

			if blen > 255 {
				return nil, errors.New("Version length > 255")
			}
			value = bytes.Join([][]byte{domainId, []byte{byte(blen & 0xff)}, id, verbytes}, []byte{})
		} else {
			return nil, errors.New("Expected semantic version")
		}
		break
	default:
		return nil, errors.New("Unknown version type")
	}

	if crcLength > 0 {
		crc, err := crcCalc(value, crcLength)

		if err != nil {
			return nil, err
		}

		value = bytes.Join([][]byte{value, crc}, []byte{})
	}

	//fmt.Printf("val=%v\n", value)
	return Wrap(value), nil
}

func Wrap(id []byte) ids.Identifier {
	return &identifier{id}
}

func Unwrap(id ids.Identifier) []byte {
	return id.Value()
}

func Parse(id string) (ids.Identifier, error) {
	return nil, nil
}

func AsLocator(id ids.Identifier) ids.Locator {
	locator, ok := id.(*locator)

	if ok {
		return locator
	}

	return getLocator(id)
}

func (id *identifier) Id() []byte {
	domainOffset := domain.DomainOffset(id.value)
	domainLength := domain.DomainLength(id.value) + domain.VersionLengthLength(id.value)
	versionLength := versionLength(id.value)
	identifierLength, _ := identifierLength(id.value)
	pathLengthLength := domain.PathLengthLength(id.value)
	fragmentLengthLength := domain.FragmentLengthLength(id.value)
	//fmt.Printf("do=%v, dl=%v, vl=%v, il=%v\n", domainOffset, domainLength, versionLength, identifierLength)
	return id.value[domainOffset+domainLength+pathLengthLength+fragmentLengthLength : domainOffset+domainLength+identifierLength-versionLength]
}

func (id *identifier) IdRoot() []byte {
	domainLength := domain.DomainLength(id.value) + domain.VersionLengthLength(id.value)
	pathLengthLength := domain.PathLengthLength(id.value)
	fragmentLengthLength := domain.FragmentLengthLength(id.value)
	//fmt.Printf("pll=%d, fll=%d\n", pathLengthLength, fragmentLengthLength)
	startOffset := domain.DomainOffset(id.value) + domainLength + pathLengthLength + fragmentLengthLength
	endOffset := startOffset + rootIdLength(id.value)

	return id.value[startOffset:endOffset]
}

func (id *identifier) Path() []byte {
	domainLength := domain.DomainLength(id.value) + domain.VersionLengthLength(id.value)
	pathLengthLength := domain.PathLengthLength(id.value)
	fragmentLengthLength := domain.FragmentLengthLength(id.value)

	//fmt.Printf("domainOffset=%d, domainLength=%d, pathLengthLength=%d\n", domain.DomainOffset(id.value), domainLength, pathLengthLength)
	idOffset := domain.DomainOffset(id.value) + domainLength + pathLengthLength + fragmentLengthLength
	//fmt.Printf("idOffset=%d, fragmentLengthLength=%d, rootIdLength=%d\n", idOffset, fragmentLengthLength, rootIdLength(id.value))
	startOffset := idOffset + rootIdLength(id.value)
	endOffset := startOffset + pathLength(id.value)

	//fmt.Printf("startOffset=%d, endOffset=%d\n", startOffset, endOffset)
	return id.value[startOffset:endOffset]
}

func (id *identifier) Fragment() []byte {
	domainLength := domain.DomainLength(id.value) + domain.VersionLengthLength(id.value)
	pathLengthLength := domain.PathLengthLength(id.value)
	fragmentLengthLength := domain.FragmentLengthLength(id.value)

	idOffset := domain.DomainOffset(id.value) + domainLength + pathLengthLength + fragmentLengthLength
	startOffset := idOffset + rootIdLength(id.value) + pathLength(id.value)
	endOffset := startOffset + fragmentLength(id.value)

	return id.value[startOffset:endOffset]
}

func (id *identifier) SchemeId() []byte {
	return domain.SchemeId(id.value)
}

func (id *identifier) DomainId() []byte {
	return id.value[:domain.DomainOffset(id.value)+domain.DomainLength(id.value)]
}

func (id *identifier) HasVersion() bool {
	return domain.VersionLengthLength(id.value) > 0
}

func (id *identifier) VersionId() []byte {
	versionLength := versionLength(id.value)
	if versionLength == 0 {
		return nil
	}
	domainOffset := domain.DomainOffset(id.value)
	domainLength := domain.DomainLength(id.value) + domain.VersionLengthLength(id.value)
	identifierLength, _ := identifierLength(id.value)
	return id.value[domainOffset+domainLength+identifierLength-versionLength:]
}

func (id *identifier) Version() version.Version {
	versionId := id.VersionId()

	if versionId == nil {
		return nil
	}

	vtype, err := domain.VersionTypeValue(id.value)

	if err == nil {
		switch vtype {
		case versiontype.NUMERIC:
			return version.NumericVersion(numericVersionValue(versionId))
		case versiontype.SEMANTIC:
			result, err := version.Parse(string(versionId))
			if err == nil {
				return result
			}
		}
	}

	return nil
}

func (id *identifier) DomainIdRoot() []byte {
	return domain.IdRoot(id.value)
}

func (id *identifier) DomainIncarnation() *uint32 {
	return domain.IncarnationValue(id.value)
}

func (id *identifier) Checksum() []byte {
	crcLength, _ := identifierCrcLength(id.value)
	if crcLength == 0 {
		return nil
	}
	return id.value[uint(len(id.value))-crcLength:]
}

func (id *identifier) sign(signatureDomain ids.SignatureDomain) (ids.Signature, error) {
	crcLength, _ := identifierCrcLength(id.value)
	signatureBytes := id.value[:uint(len(id.value))-crcLength]
	signatureBytes[0] = signatureBytes[0] & 0x3f
	return signatureDomain.CreateSignature(struct {
		domainId       []byte
		id             []byte
		signatureBytes []byte
	}{domainId: id.DomainId(),
		id:             id.Id(),
		signatureBytes: signatureBytes})
}

func (id *identifier) Scheme() ids.Scheme {
	result, _ := scheme.Get(context.Background(), scheme.Selector{Id: id.SchemeId()})
	return result
}

func (id *identifier) Domain() ids.IdentityDomain {
	domain, _ := domain.Get(context.Background(), domain.Selector{Id: id.DomainId()})
	return domain
}

func (id *identifier) IsUndefined() bool {
	return id.DomainId() == nil
}

func (id *identifier) IsNull() bool {
	return id.value == nil
}

func (id *identifier) IsValid() bool {
	return id.Validate() == nil
}

func (id *identifier) Validate() error {
	return validate(id.value)
}

func (id *identifier) Matches(other ids.Identifier) bool {

	if id.Equals(other) {
		return true
	}

	if bytes.Equal(id.DomainId(), other.DomainId()) {
		return false
	}

	as, err := Await(other.As(id.Domain()))
	if err != nil {
		return false
	}

	return id.Equals(as)
}

func (id *identifier) Equals(other ids.Identifier) bool {
	if other == nil {
		return false
	}

	return bytes.Equal(id.value, Unwrap(other))
}

func (id *identifier) CompareTo(o ids.Identifier) int {
	result := bytes.Compare(id.DomainId(), o.DomainId())

	if result != 0 {
		return result
	}

	return bytes.Compare(id.Id(), o.Id())
}

func (id *identifier) Encode(seperator string, encoders ...encodertype.EncoderType) string {
	domainEncoder := encodertype.BASE62
	idEncoder := encodertype.BASE62
	versionEncoder := encodertype.BASE62

	if len(encoders) > 0 {
		domainEncoder = encoders[0]

		if len(encoders) > 1 {
			idEncoder = encoders[1]
		} else {
			idEncoder = domainEncoder
		}

		if len(encoders) > 2 {
			versionEncoder = encoders[2]
		} else {
			versionEncoder = domainEncoder
		}
	}

	if seperator != "" {
		// id effectivelt removes the CRC from the Identifier
		// so format the domain as a non crc domain
		did := id.DomainId()
		did = bytes.Join([][]byte{[]byte{did[0] & 0x3f}, did[1:]}, []byte{})
		dom, _ := encoding.Encode(did, domainEncoder)
		eid, _ := encoding.Encode(id.Id(), idEncoder)

		result := dom + seperator + eid

		if id.HasVersion() {
			version, _ := encoding.Encode(id.VersionId(), versionEncoder)
			result = result + seperator + version
		}

		return result
	}

	result, _ := encoding.Encode(id.value, domainEncoder)
	return result
}

func (id *identifier) String() string {
	return id.Encode("", encodertype.BASE62)
}

func (id *identifier) Value() []byte {
	return id.value
}

func (id *identifier) Sign(signatureDomain ids.SignatureDomain) (ids.Signature, error) {
	return nil, nil
}

func (id *identifier) IsFor(typeId ids.TypeIdentifier) bool {
	return id.Domain().IsFor(typeId)
}

func (id *identifier) As(domain ids.IdentityDomain, at ...time.Time) (chan ids.Identifier, chan error) {
	cid := make(chan ids.Identifier, 1)
	cerr := make(chan error)

	go func() {
		mapping, err := mappings.Get(context.Background(), mappings.Selector{From: id, To: domain, At: time.Now()})

		if err != nil {
			cerr <- err
		} else {
			cid <- mapping.ToId()
		}

		close(cid)
		close(cerr)
	}()

	return cid, cerr
}

func (id *identifier) MapFrom(from ids.Identifier, between ...time.Time) chan error {
	return mappings.Map(context.Background(), from, id, between...)
}

func (id *identifier) MapTo(to ids.Identifier, between ...time.Time) chan error {
	return mappings.Map(context.Background(), id, to, between...)
}

func (id *identifier) MapBetween(other ids.Identifier, between ...time.Time) chan error {
	err := id.MapFrom(other, between...)
	if err != nil {
		return err
	}
	return id.MapTo(other, between...)
}

func (id *identifier) Get() (interface{}, error) {
	return AsLocator(id).Get()
}

func (id *identifier) GetAs(typeId ids.TypeIdentifier) (interface{}, error) {
	return AsLocator(id).GetAs(typeId)
}

func (id *identifier) Resolve() (chan interface{}, chan error) {
	return AsLocator(id).Resolve()
}

func (id *identifier) ResolveAs(typeId ids.TypeIdentifier) (chan interface{}, chan error) {
	return AsLocator(id).ResolveAs(typeId)
}

func (id *identifier) TypeId() ids.TypeIdentifier {
	return id.Domain().TypeId()
}

func identifierVersionType(value []byte) (versiontype.VersionType, error) {
	return domain.VersionTypeValue(value)
}

func rootIdLength(value []byte) uint {
	identifierLength, _ := identifierLength(value)
	//fmt.Printf("identifierLength=%d, pathLengthLength=%d, pathLength=%d, versionLength=%d\n",
	//	identifierLength, domain.PathLengthLength(value), pathLength(value), versionLength(value))
	return identifierLength -
		domain.PathLengthLength(value) -
		pathLength(value) -
		domain.FragmentLengthLength(value) -
		fragmentLength(value) -
		versionLength(value)
}

func versionLength(value []byte) uint {
	lengthLength := domain.VersionLengthLength(value)

	if lengthLength > 0 {
		domainOffset := domain.DomainOffset(value)
		domainLength := domain.DomainLength(value)
		return uint(value[domainOffset+domainLength])
	}

	return 0
}

func pathLength(value []byte) uint {
	lengthLength := domain.PathLengthLength(value)
	//fmt.Printf("value=%v\n", value)
	if lengthLength > 0 {
		domainOffset := domain.DomainOffset(value)
		domainLength := domain.DomainLength(value)
		versionLengthLength := domain.VersionLengthLength(value)
		return uint(value[domainOffset+domainLength+versionLengthLength])
	}

	return 0
}

func fragmentLength(value []byte) uint {
	lengthLength := domain.FragmentLengthLength(value)

	if lengthLength > 0 {
		domainOffset := domain.DomainOffset(value)
		domainLength := domain.DomainLength(value)
		versionLengthLength := domain.VersionLengthLength(value)
		pathLengthLength := domain.PathLengthLength(value)
		return uint(value[domainOffset+domainLength+versionLengthLength+pathLengthLength])
	}

	return 0
}

func numericVersionValue(versionValue []byte) uint32 {
	vlen := len(versionValue)
	if vlen > 0 {
		if vlen == 1 {
			return uint32(versionValue[0])
		}
		if vlen == 2 {
			return uint32(ntoh.U16(versionValue, 0))
		}
		return ntoh.U32(versionValue, 0)
	}

	return 0
}

func identifierCrcLength(value []byte) (uint, error) {
	crcLength, err := domain.CrcLengthValue(value)
	if err != nil {
		return 0, err
	}
	return uint(crcLength / 8), nil
}

func identifierLength(value []byte) (uint, error) {
	identifierCrcLength, err := identifierCrcLength(value)

	if err != nil {
		return 0, err
	}

	length := uint(len(value)) - domain.DomainOffset(value) -
		domain.DomainLength(value) - domain.VersionLengthLength(value) -
		identifierCrcLength

	if length < 0 {
		return 0, nil
	}
	return length, nil
}

func validate(value []byte) error {
	schemeLength := domain.SchemeLength(value)
	domainLength := domain.DomainLength(value)
	incarnationLength := domain.IncarnationLength(value)
	crcLength, err := identifierCrcLength(value)

	if err != nil {
		return err
	}

	crc, err := crcCalc(value[:uint(len(value))-crcLength], crcLength*8)

	if err != nil {
		return err
	}

	if !(crcLength <= 4 &&
		incarnationLength <= 4 &&
		schemeLength <= uint(len(value))-domainLength-crcLength &&
		domainLength > incarnationLength &&
		domainLength <= uint(len(value))-domainLength-crcLength &&
		bytes.Equal(value[uint(len(value))-crcLength:], crc)) {
		return fmt.Errorf("Invalid value: %v", value)
	}

	return nil
}

var crc8Table *crc8.Table = crc8.MakeTable(crc8.CRC8_MAXIM)
var crc16Table *crc16.Table = crc16.MakeTable(crc16.IBM)
var crc32Table *crc32.Table = crc32.MakeTable(crc32.IEEE)

func crcCalc(value []byte, crcLength uint) ([]byte, error) {
	switch crcLength {
	case 0:
		return make([]byte, 0), nil
	case 8:
		buf := [1]byte{crc8.Checksum(value, crc8Table)}
		return buf[:], nil
	case 16:
		buf := make([]byte, 2)
		return hton.U16(buf, 0, crc16.Checksum(value, crc16Table)), nil
	case 32:
		buf := make([]byte, 4)
		return hton.U32(buf, 0, crc32.Checksum(value, crc32Table)), nil
	default:
		return nil, errors.New("Invalid crc length")
	}
}

func ignore(err error) {

}
