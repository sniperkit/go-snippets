package identitydomain

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/domaintype"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

func init() {
	domain.RegisterJSONUnmarshaller(domaintype.IDENTITY, unmarshalJSON)
}

type identityDomain struct {
	ids.Domain
	idGenerator             ids.IdGenerator
	sequenceLeaseHolderId   ids.Identifier
	nextSequenceLeaseExpiry *time.Time
	sequenceRoot            []byte
	sequenceIncarnation     *uint32
}

func unmarshalJSON(unmarshalContext context.Context, json map[string]interface{}) (ids.Domain, error) {
	rootId, err := base62.Decode(json["id"].(string))

	if err != nil {
		return nil, err
	}

	var incarnation *uint32
	var crcLength uint
	var versionType = versiontype.UNVERSIONED
	var hasPaths = false
	var hasFragments = false
	info := make(map[interface{}]interface{})

	if value, ok := json["hasPaths"]; ok {
		if has, ok := value.(bool); ok {
			hasPaths = has
		}
	}

	if value, ok := json["hasFragments"]; ok {
		if has, ok := value.(bool); ok {
			hasFragments = has
		}
	}

	for key, value := range json {
		info[key] = value
	}

	return New(unmarshalContext.Value("scheme").(ids.Scheme), rootId, incarnation, crcLength, versionType, hasPaths, hasFragments, info)
}

func New(scheme ids.Scheme, rootId []byte, incarnation *uint32, crcLength uint, versionType versiontype.VersionType, hasPaths bool, hasFragments bool, infos ...map[interface{}]interface{}) (ids.IdentityDomain, error) {
	base, err := domain.New(scheme.Id(), rootId, incarnation, crcLength, versionType, hasPaths, hasFragments, infos...)

	if err != nil {
		return nil, err
	}

	return &identityDomain{
		Domain: base}, nil
}

func WithIncarnation(root ids.Domain, incarnation uint32, crcLength uint, infos ...map[interface{}]interface{}) (ids.IdentityDomain, error) {
	base, err := domain.WithIncarnation(root, incarnation, crcLength, infos...)

	if err != nil {
		return nil, err
	}

	return &identityDomain{
		Domain: base}, nil
}

func WithCrc(root ids.Domain, crcLength uint, infos ...map[interface{}]interface{}) (ids.IdentityDomain, error) {
	base, err := domain.WithCrc(root, crcLength, infos...)

	if err != nil {
		return nil, err
	}

	return &identityDomain{
		Domain: base}, nil
}

func (identityDomain *identityDomain) NewIdentifier() (ids.Identifier, error) {
	if identityDomain.idGenerator != nil {
		id, err := identityDomain.idGenerator.GenerateId()

		if err != nil {
			return nil, err
		}

		var nilVersion version.Version
		return identityDomain.ToId(id, nilVersion)
	}

	return nil, errors.New("Id creation not initialized")
}

func (identityDomain *identityDomain) ToId(id []byte, ver version.Version) (ids.Identifier, error) {
	return identifier.New(identityDomain, id, ver)
}

func (identityDomain *identityDomain) IsRootOf(domain ids.Domain) bool {
	return identityDomain.Incarnation() == nil && identityDomain.CrcLength() == 0 && bytes.Equal(domain.IdRoot(), identityDomain.IdRoot())
}

func (identityDomain *identityDomain) IsRoot() bool {
	return identityDomain.Incarnation() == nil && identityDomain.CrcLength() == 0
}

func (identityDomain *identityDomain) AsRoot() ids.IdentityDomain {
	root, _ := New(identityDomain.Scheme(), identityDomain.IdRoot(), nil, 0, identityDomain.VersionType(), identityDomain.HasPaths(), identityDomain.HasFragments())
	return root
}

func (identityDomain *identityDomain) SetIfChanged(root []byte, incarnation *int32, name string, description string, source string, typeId ids.TypeIdentifier) bool {

	//changed = super.setIfChanged(root, incarnation)

	/*
		if (type!=null) {
			if (!type.equals(identityDomain.type)) {
				identityDomain.type=type;
				changed=true;
			}
		}

		if (!(name==null||"".equals(name))) {
			if (!name.equals(identityDomain.name)) {
				identityDomain.name = name;
				changed=true;
			}
		}

		if (!(description==null||"".equals(description))) {
			if (!description.equals(identityDomain.description)) {
				identityDomain.description = description;
				changed=true;
			}
		}

		if (!(source==null||"".equals(source))) {
			if (!source.equals(identityDomain.source)) {
				identityDomain.source = source;
				changed=true;
			}
		}
	*/

	return false //changed
}

func (identityDomain *identityDomain) SequenceLeaseHolderId() ids.Identifier {
	return identityDomain.sequenceLeaseHolderId
}

func (identityDomain *identityDomain) HasSequenceLeaseHolderId() bool {
	return identityDomain.sequenceLeaseHolderId != nil
}

func (identityDomain *identityDomain) NextSequenceLeaseExpiry() *time.Time {
	return identityDomain.nextSequenceLeaseExpiry
}

func (identityDomain *identityDomain) SequenceDomain() ids.SequenceDomain {
	//return DomainStore.getSequenceDomain((IdentityDomainInfo)domainInfo);
	return nil
}

func (identityDomain *identityDomain) SequenceRoot() []byte {
	if identityDomain.sequenceRoot == nil {
		return nil
	}
	return identityDomain.sequenceRoot
}

func (identityDomain *identityDomain) SequenceDomainId() []byte {
	id, _ := domain.ToId(identityDomain.SchemeId(), identityDomain.sequenceRoot, identityDomain.sequenceIncarnation, 0, versiontype.UNVERSIONED, false, false)
	return id
}

func (identityDomain *identityDomain) SequenceIncarnation() uint32 {
	if identityDomain.sequenceIncarnation == nil {
		return 0
	}
	return *identityDomain.sequenceIncarnation
}

/*
	aquireSequenceDomain( int incarnation, Period leasePeriod, SequenceNumber<Type> lastSequence ) {

		// TODO - identityDomain really needs a lock - so should be in a transaction

		IdentityDomainInfo identityDomainInfo=(IdentityDomainInfo)domainInfo;
		Date nextExpiry=identityDomainInfo.nextSequenceLeaseExpiry;

		if (nextExpiry==null||LeaseHolderIdentifier.VM_ID.equals(getSequenceLeaseHolderId())||
				nextExpiry.getTime()<System.currentTimeMillis()) {

			if (identityDomainInfo.sequenceIncarnation==null||incarnation>identityDomainInfo.sequenceIncarnation)
				identityDomainInfo.sequenceIncarnation=incarnation;
			identityDomainInfo.nextSequenceLeaseExpiry=DateTime.now().plus(leasePeriod).toDate();
			identityDomainInfo.sequenceLeaseHolderId=LeaseHolderIdentifier.VM_ID;
			DomainStore.post(domainInfo);

		} else {
			throw new LeaseAquisitionException(identityDomainInfo.nextSequenceLeaseExpiry);
		}

		return DomainStore.getSequenceDomain(identityDomainInfo,identityDomainInfo.getSequenceRoot(),identityDomainInfo.sequenceIncarnation);
	}

	void releaseSequenceDomain() {
		((IdentityDomainInfo)domainInfo).nextSequenceLeaseExpiry=null;
		((IdentityDomainInfo)domainInfo).sequenceLeaseHolderId=null;
		DomainStore.post(domainInfo);
	}

	getNextSequenceNumber() {
		return getSequenceDomain().nextSequenceNumber();
	}
*/

func (identityDomain *identityDomain) RegisterMappingProvider(to ids.IdentityDomain, provider func() (ids.IdentityDomain, error)) (ids.IdentityDomain, error) {
	//MAPPING_PROVIDERS.set(to, provider)
	return nil, nil
}

func (identityDomain *identityDomain) UnregisterMappingProvider(to ids.IdentityDomain) {
	//MAPPING_PROVIDERS.delete(to)
}
