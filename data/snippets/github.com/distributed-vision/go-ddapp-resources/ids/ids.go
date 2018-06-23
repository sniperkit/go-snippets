package ids

import (
	"reflect"
	"sync"
	"time"

	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

type Domain interface {
	String() string
	Encode(encoder encodertype.EncoderType) string
	ToJSON() string
	IsFor(typeId TypeIdentifier) bool
	Matches(other Domain) bool
	Id() []byte
	SchemeId() []byte
	Scheme() Scheme
	IdRoot() []byte
	IsRoot() bool
	IsRootOf(domain Domain) bool
	Incarnation() *uint32
	CrcLength() uint
	HasPaths() bool
	HasFragments() bool
	VersionType() versiontype.VersionType
	Name() string
	Description() string
	TypeId() TypeIdentifier
	Source() string
	InfoValue(interface{}) interface{}
	Equals(other Domain) bool
}

type SchemeVisibility int

type SchemeFormat int
type DomainType int

type Scheme interface {
	Domain
	Visibility() SchemeVisibility
	Format() SchemeFormat
	RegisterResolvers() error
}

type IdentityDomain interface {
	Domain
}

type SignatureDomain interface {
	IdentityDomain
	CreateSignature(elements interface{}) (Signature, error)
}

type SequenceDomain interface {
	IdentityDomain
}

type Identifier interface {
	String() string
	Id() []byte
	IdRoot() []byte
	Path() []byte
	Fragment() []byte
	SchemeId() []byte
	DomainId() []byte
	VersionId() []byte
	DomainIdRoot() []byte
	DomainIncarnation() *uint32
	Checksum() []byte
	Scheme() Scheme
	Domain() IdentityDomain
	Version() version.Version
	HasVersion() bool
	IsUndefined() bool
	IsNull() bool
	IsValid() bool
	Validate() error
	Value() []byte

	Sign(signatureDomain SignatureDomain) (Signature, error)

	Matches(other Identifier) bool
	Equals(other Identifier) bool

	As(domain IdentityDomain, at ...time.Time) (chan Identifier, chan error)
	MapFrom(from Identifier, between ...time.Time) chan error
	MapTo(from Identifier, between ...time.Time) chan error
	MapBetween(from Identifier, between ...time.Time) chan error

	Get() (interface{}, error)
	GetAs(typeId TypeIdentifier) (interface{}, error)
	Resolve() (chan interface{}, chan error)
	ResolveAs(typeId TypeIdentifier) (chan interface{}, chan error)
}

type Locator interface {
	Identifier
}

type Mapping interface {
	FromId() Identifier
	ToDomain() IdentityDomain
	From() time.Time
	To() time.Time
	ToId() Identifier
}

var typeInitFunctions = []func(){}
var typeInitMutex = sync.Mutex{}

func OnLocalTypeInit(initFunction func()) {
	typeInitMutex.Lock()
	if NewLocalTypeId == nil {
		typeInitFunctions = append(typeInitFunctions, initFunction)
	} else {
		initFunction()
	}
	typeInitMutex.Unlock()
}

func LocalTypeInit(localIdInitialiser func(gotype reflect.Type) TypeIdentifier,
	publicIdInitialiser func(domainValue interface{}, id []byte, idVersion version.Version) (TypeIdentifier, error)) {
	typeInitMutex.Lock()
	if NewLocalTypeId == nil {
		NewLocalTypeId = localIdInitialiser
		NewTypeId = publicIdInitialiser
		for _, initFunction := range typeInitFunctions {
			initFunction()
		}
		typeInitFunctions = []func(){}
	}
	typeInitMutex.Unlock()
}

var NewLocalTypeId func(gotype reflect.Type) TypeIdentifier
var NewTypeId func(domainValue interface{}, id []byte, idVersion version.Version) (TypeIdentifier, error)

type SignatureElements interface {
	SignatureBytes() []byte
	Signature() Signature
}

type Signature interface {
	Identifier
	Elements() (SignatureElements, error)
}

type TypeIdentifier interface {
	Signature
	IsAssignableFrom(typeId TypeIdentifier) bool
}

type IdGenerator interface {
	GenerateId() ([]byte, error)
}
