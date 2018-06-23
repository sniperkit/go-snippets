package domain

import (
	"bytes"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

type idTest struct {
	schemeId    []byte
	idRoot      []byte
	incarnation *uint32
	crcLength   uint
	versionType versiontype.VersionType
	encodedId   string
}

var incarnation0 = uint32(0)
var incarnation1 = uint32(1)
var incarnation1011 = uint32(1011)
var incarnation201011 = uint32(201011)

var idTests = []idTest{
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 0, versiontype.UNVERSIONED, "H7C"},
	{base62.MustDecode("1"), base62.MustDecode("12"), nil, 0, versiontype.UNVERSIONED, "H8C"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), nil, 0, versiontype.UNVERSIONED, "4keRnR"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation0, 0, versiontype.UNVERSIONED, "5tbciO"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 0, versiontype.UNVERSIONED, "5tbciP"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1011, 0, versiontype.UNVERSIONED, "OQJ9VFj"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation201011, 0, versiontype.UNVERSIONED, "6kzTmoi1Vj"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation0, 0, versiontype.UNVERSIONED, "1d9wvIxu4"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation1, 0, versiontype.UNVERSIONED, "1d9wvIxu5"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation1011, 0, versiontype.UNVERSIONED, "6kzKBJzG9T"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation201011, 0, versiontype.UNVERSIONED, "1rvLheMVxt1pz"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 8, versiontype.UNVERSIONED, "1QGxu"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 16, versiontype.UNVERSIONED, "1QHEQ"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 32, versiontype.UNVERSIONED, "1QHUw"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 0, versiontype.NUMERIC, "1QHlS"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 0, versiontype.SEMANTIC, "1QIpW"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 8, versiontype.UNVERSIONED, "5tciuX"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 16, versiontype.UNVERSIONED, "5tdp6f"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 32, versiontype.UNVERSIONED, "5tevIn"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation0, 0, versiontype.NUMERIC, "5tg1Uu"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation0, 0, versiontype.SEMANTIC, "5tkQHQ"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 32, versiontype.NUMERIC, "5tjK5J"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 8, versiontype.SEMANTIC, "5tlWTZ"},
}

func TestDomainIdFormatting(t *testing.T) {

	for _, test := range idTests {
		id, err := ToId(test.schemeId, test.idRoot, test.incarnation, test.crcLength, test.versionType, false, false)

		if err != nil {
			t.Errorf("TestDomainIdFormatting: ToId failed with err %s:", err)
		}

		if base62.Encode(id) != test.encodedId {
			t.Errorf("TestDomainIdFormatting: ToId failed expected %s: got %s", test.encodedId, base62.Encode(id))
		}

		if !bytes.Equal(SchemeId(id), test.schemeId) {
			t.Errorf("TestDomainIdFormatting: ToId failed expected schemeId %v: got %v", test.schemeId, SchemeId(id))
		}

		if !bytes.Equal(test.idRoot, IdRoot(id)) {
			t.Errorf("TestDomainIdFormatting: ToId failed expected idRoot %v: got %v", test.idRoot, IdRoot(id))
		}

		if test.incarnation == nil {
			if IncarnationValue(id) != nil {
				t.Errorf("TestDomainIdFormatting: ToId failed expected incarnation %v: got %v", test.incarnation, *IncarnationValue(id))
			}
		} else {
			if IncarnationValue(id) == nil {
				t.Errorf("TestDomainIdFormatting: ToId failed expected incarnation %v: got %v", test.incarnation, IncarnationValue(id))
			} else {
				if *test.incarnation != *IncarnationValue(id) {
					t.Errorf("TestDomainIdFormatting: ToId failed expected incarnation %v: got %v", *test.incarnation, *IncarnationValue(id))
				}
			}
		}

		crcLengthValue, err := CrcLengthValue(id)

		if err != nil {
			t.Errorf("TestDomainIdFormatting: ToId failed with err %s:", err)
		}

		if test.crcLength != crcLengthValue {
			t.Errorf("TestDomainIdFormatting: ToId failed expected crcLength %v: got %v", test.crcLength, crcLengthValue)
		}

		versionTypeValue, err := VersionTypeValue(id)

		if err != nil {
			t.Errorf("TestDomainIdFormatting: ToId failed with err %s:", err)
		}

		if test.versionType != versionTypeValue {
			t.Errorf("TestDomainIdFormatting: ToId failed expected versionTypeValue %v: got %v", test.versionType, versionTypeValue)
		}

		domain := Wrap(id)

		if domain.String() != test.encodedId {
			t.Errorf("TestDomainIdFormatting: ToId failed expected %s: got %s", test.encodedId, domain.String())
		}

		if !bytes.Equal(domain.SchemeId(), test.schemeId) {
			t.Errorf("TestDomainIdFormatting: ToId failed expected schemeId %v: got %v", test.schemeId, domain.SchemeId())
		}

		if !bytes.Equal(test.idRoot, domain.IdRoot()) {
			t.Errorf("TestDomainIdFormatting: ToId failed expected idRoot %v: got %v", test.idRoot, domain.IdRoot())
		}

		if test.incarnation == nil {
			if domain.Incarnation() != nil {
				t.Errorf("TestDomainIdFormatting: ToId failed expected incarnation %v: got %v", test.incarnation, *domain.Incarnation())
			}
		} else {
			if domain.Incarnation() == nil {
				t.Errorf("TestDomainIdFormatting: ToId failed expected incarnation %v: got %v", test.incarnation, domain.Incarnation())
			} else {
				if *test.incarnation != *domain.Incarnation() {
					t.Errorf("TestDomainIdFormatting: ToId failed expected incarnation %v: got %v", *test.incarnation, *domain.Incarnation())
				}
			}
		}

		if domain.CrcLength() != crcLengthValue {
			t.Errorf("TestDomainIdFormatting: ToId failed expected crcLength %v: got %v", crcLengthValue, domain.CrcLength())
		}

		if domain.VersionType() != versionTypeValue {
			t.Errorf("TestDomainIdFormatting: ToId failed expected versionTypeValue %v: got %v", versionTypeValue, domain.VersionType())
		}
	}
}

func TestDomainAccessors(t *testing.T) {

	for _, test := range idTests {
		domain, err := New(test.schemeId, test.idRoot, test.incarnation, test.crcLength, test.versionType, false, false)

		if err != nil {
			t.Errorf("TestDomainAccessors: New failed with err %s:", err)
		}

		if !bytes.Equal(domain.SchemeId(), test.schemeId) {
			t.Errorf("TestDomainAccessors: domain.SchemeId failed expected %v: got %v", test.schemeId, domain.SchemeId())
		}

		if !bytes.Equal(domain.IdRoot(), test.idRoot) {
			t.Errorf("TestDomainAccessors: domain.IdRoot failed expected %v: got %v", test.idRoot, domain.IdRoot())
		}

		if domain.VersionType() != test.versionType {
			t.Errorf("TestDomainAccessors: domain.VersionType failed expected %v: got %v", test.versionType, domain.VersionType())
		}

		if test.versionType == versiontype.UNVERSIONED {
			if VersionLengthLength(domain.Id()) != 0 {
				t.Errorf("TestDomainAccessors: VersionLengthLength failed expected %v: got %v", 0, VersionLengthLength(domain.Id()))
			}

		} else {
			if VersionLengthLength(domain.Id()) != 1 {
				t.Errorf("TestDomainAccessors: VersionLengthLength failed expected %v: got %v", 1, VersionLengthLength(domain.Id()))
			}
		}
	}
}

func TestWithIncarnation(t *testing.T) {

	for _, test := range idTests {
		root, err := New(test.schemeId, test.idRoot, nil, 0, test.versionType, false, false)

		if err != nil {
			t.Errorf("TestWithIncarnation: New failed with err %s:", err)
		}

		if !bytes.Equal(root.SchemeId(), test.schemeId) {
			t.Errorf("TestWithIncarnation: domain.SchemeId failed expected %v: got %v", test.schemeId, root.SchemeId())
		}

		if !bytes.Equal(root.IdRoot(), test.idRoot) {
			t.Errorf("TestWithIncarnation: domain.IdRoot failed expected %v: got %v", test.idRoot, root.IdRoot())
		}

		if root.VersionType() != test.versionType {
			t.Errorf("TestWithIncarnation: domain.IdRoot failed expected %v: got %v", test.versionType, root.VersionType())
		}

		if !root.IsRoot() {
			t.Errorf("TestDomainWithCrc Failed: !root.IsRoot")
		}

		if test.incarnation != nil {
			domain, err := WithIncarnation(root, *test.incarnation, test.crcLength)

			if err != nil {
				t.Errorf("TestWithIncarnation: WithIncarnation failed with err %s:", err)
			}

			if !bytes.Equal(root.IdRoot(), domain.IdRoot()) {
				t.Errorf("TestWithIncarnation Failed: root.IdRoot !=  domain.IdRoot: expected: '%v' got '%v'", root.IdRoot(), domain.IdRoot())
			}

			if domain.Incarnation() == nil || *domain.Incarnation() != *test.incarnation {
				t.Errorf("TestWithIncarnation Failed: root.Incarnation: expected: '%v' got '%v'", test.incarnation, domain.Incarnation())
			}

			if domain.CrcLength() != test.crcLength {
				t.Errorf("TestWithIncarnation Failed: root.CrcLength: expected: '%v' got '%v'", test.crcLength, domain.CrcLength())
			}

			if domain.IsRoot() {
				t.Errorf("TestWithIncarnation Failed: domain IsRoot")
			}

			if !root.IsRootOf(domain) {
				t.Errorf("TestWithIncarnation Failed: !root.IsRootOf domain")
			}

			if !root.Matches(domain) {
				t.Errorf("TestWithIncarnation Failed: !root.Matches domain")
			}
		}
	}
}

func TestGoTypeIdDomainDefinition(t *testing.T) {
	id, err := DecodeId(encodertype.BASE62, "2", "0")

	if err != nil {
		t.Errorf("TestGoTypeIdDomainDefinition: ToId failed with err %s:", err)
	}

	if base62.Encode(id) != "YAC" {
		t.Errorf("TestGoTypeIdDomainDefinition: ToId failed expected %s: got %s", "YAC", base62.Encode(id))
	}

	id = MustDecodeId(encodertype.BASE62, "2", "0")

	if base62.Encode(id) != "YAC" {
		t.Errorf("TestGoTypeIdDomainDefinition: ToId failed expected %s: got %s", "YAC", base62.Encode(id))
	}

	incarnation := uint32(0)

	id, err = DecodeId(encodertype.BASE62, "T", "0", &incarnation)

	if base62.Encode(id) != "2DAEU1w" {
		t.Errorf("TestGoTypeIdDomainDefinition: ToId failed expected %s: got %s", "2DAEU1w", base62.Encode(id))
	}

}
