package identifier_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/util/random"
	"github.com/distributed-vision/go-resources/version"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

type testDomain struct {
	schemeId             []byte
	idRoot               []byte
	incarnation          *uint32
	crcLength            uint
	versionType          versiontype.VersionType
	encodedId            string
	encodedPathOnlyId    string
	encodedFragOnlyId    string
	encodedPathAndFragId string
}

var incarnation0 = uint32(0)
var incarnation1 = uint32(1)
var incarnation1011 = uint32(1011)
var incarnation201011 = uint32(201011)

var testDomains = []testDomain{
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 0, versiontype.UNVERSIONED, "", "1QPDu", "1QKxe", "1QTUA"},
	{base62.MustDecode("1"), base62.MustDecode("12"), nil, 0, versiontype.UNVERSIONED, "", "1QPEu", "1QKye", "1QTVA"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), nil, 0, versiontype.UNVERSIONED, "", "OScmnUT", "ORS7UBx", "OTnS6mz"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation0, 0, versiontype.UNVERSIONED, "", "5uAoyW", "5ttDqS", "5uSQ6a"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 0, versiontype.UNVERSIONED, "", "5uAoyX", "5ttDqT", "5uSQ6b"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1011, 0, versiontype.UNVERSIONED, "", "OSeU7ql", "ORTooYF", "OTp9R9H"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation201011, 0, versiontype.UNVERSIONED, "", "6ldRY5Obnn", "6lJSfS3Jel", "6lxQQijtwp"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation0, 0, versiontype.UNVERSIONED, "", "1dJd0SKoK", "1dEmxseMC", "1dOT321GS"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation1, 0, versiontype.UNVERSIONED, "", "1dJd0SKoL", "1dEmxseMD", "1dOT321GT"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation1011, 0, versiontype.UNVERSIONED, "", "6ldHwafqRX", "6lJJ3xKYIV", "6lxGpE18aZ"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation201011, 0, versiontype.UNVERSIONED, "", "1s6L2n9WVyAc7", "1s0qNDl1EvbE3", "1sBpiMY1n0k0B"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 8, versiontype.UNVERSIONED, "", "1QPUQ", "1QLEA", "1QTkg"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 16, versiontype.UNVERSIONED, "", "1QPkw", "1QLUg", "1QU1C"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 32, versiontype.UNVERSIONED, "", "1QQ1S", "1QLlC", "1QUHi"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 0, versiontype.NUMERIC, "", "1QQHy", "1QM1i", "1QUYE"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 0, versiontype.SEMANTIC, "", "1QRM2", "1QN5m", "1QVcI"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 8, versiontype.UNVERSIONED, "", "5uBvAf", "5tuK2b", "5uTWIj"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 16, versiontype.UNVERSIONED, "", "5uD1Mn", "5tvQEj", "5uUcUr"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 32, versiontype.UNVERSIONED, "", "5uE7Yv", "5twWQr", "5uVigz"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation0, 0, versiontype.NUMERIC, "", "5uFDl2", "5txccy", "5uWot6"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation0, 0, versiontype.SEMANTIC, "", "5uJcXY", "5u21PU", "5ubDfc"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 32, versiontype.NUMERIC, "", "5uIWLR", "5u0vDN", "5ua7TV"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 8, versiontype.SEMANTIC, "", "5uKijh", "5u37bd", "5ucJrl"},
}

func TestCreateIdsWithPaths(t *testing.T) {
	for _, test := range testDomains {
		domainId, err := domain.ToId(test.schemeId, test.idRoot, test.incarnation, test.crcLength, test.versionType, true, false)

		if err != nil {
			t.Errorf("TestCreateIdsWithPaths: ToId failed with err %s:", err)
			continue
		}

		if base62.Encode(domainId) != test.encodedPathOnlyId {
			t.Errorf("TestCreateIdsWithPaths: ToId failed expected %s: got %s", test.encodedPathOnlyId, base62.Encode(domainId))
		}

		domain := domain.Wrap(domainId)

		var idversion version.Version

		switch domain.VersionType() {
		case versiontype.NUMERIC:
			idversion = version.NumericVersion(rand.Int63())
		case versiontype.SEMANTIC:
			idversion = &version.SemanticVersion{
				Major: uint64(rand.Int63n(50)),
				Minor: uint64(rand.Int63n(50)),
				Patch: uint64(rand.Int63n(100))}
		}

		idlen := int(rand.Int63n(24))
		idbytes := random.RandomBytes(idlen)

		id, err := identifier.New(domain, idbytes, idversion)

		if err != nil {
			t.Errorf("TestCreateIdsWithPaths: identifier.New failed with err %s:", err)
			continue
		}

		if !id.IsValid() {
			t.Errorf("TestCreateIdsWithPaths: identifier.New failed id not valid:")
		}

		if !bytes.Equal(idbytes, id.Id()) {
			t.Errorf("TestCreateIdsWithPaths: id.Id %v: got %v", idbytes, id.Id())
		}

		if !bytes.Equal(id.IdRoot(), id.Id()) {
			t.Errorf("TestCreateIdsWithPaths: id.Id/id.IdRoot %v: got %v", id.IdRoot(), id.Id())
		}

		if len(id.Path()) > 0 {
			t.Errorf("TestCreateIdsWithPaths: id.Path len failed expected %d: got %d", 0, len(id.Path()))
		}

		if len(id.Fragment()) > 0 {
			t.Errorf("TestCreateIdsWithPaths: id.Fragment len failed expected %d: got %d", 0, len(id.Fragment()))
		}

		pathlen := int(rand.Int63n(24))
		pathbytes := random.RandomBytes(pathlen)

		id, err = identifier.New(domain, idbytes, idversion, pathbytes)

		if err != nil {
			t.Errorf("TestCreateIdsWithPaths: identifier.New failed with err %s:", err)
			continue
		}

		if !id.IsValid() {
			t.Errorf("TestCreateIdsWithPaths: identifier.New failed id not valid:")
		}

		if !bytes.Equal(idbytes, id.IdRoot()) {
			t.Errorf("TestCreateIdsWithPaths: id.IdRoot %v: got %v", idbytes, id.IdRoot())
		}

		if !bytes.Equal(pathbytes, id.Path()) {
			t.Errorf("TestCreateIdsWithPaths: id.Path %v: got %v", pathbytes, id.Path())
		}

		if len(id.Fragment()) > 0 {
			t.Errorf("TestCreateIdsWithPaths: id.Fragment len failed expected %d: got %d", 0, len(id.Fragment()))
		}

		if idversion == nil {
			id, err = identifier.New(domain, idbytes, pathbytes)

			if err != nil {
				t.Errorf("TestCreateIdsWithPaths: identifier.New failed with err %s:", err)
				continue
			}

			if !id.IsValid() {
				t.Errorf("TestCreateIdsWithPaths: identifier.New failed id not valid:")
			}

			if !bytes.Equal(idbytes, id.IdRoot()) {
				t.Errorf("TestCreateIdsWithPaths: id.IdRoot %v: got %v", idbytes, id.IdRoot())
			}

			if !bytes.Equal(pathbytes, id.Path()) {
				t.Errorf("TestCreateIdsWithPaths: id.Path %v: got %v", pathbytes, id.Path())
			}

			if len(id.Fragment()) > 0 {
				t.Errorf("TestCreateIdsWithPaths: id.Fragment len failed expected %d: got %d", 0, len(id.Fragment()))
			}
		}
	}
}

func TestCreateIdsWithFragments(t *testing.T) {
	for _, test := range testDomains {
		domainId, err := domain.ToId(test.schemeId, test.idRoot, test.incarnation, test.crcLength, test.versionType, false, true)

		if err != nil {
			t.Errorf("TestCreateIdsWithFragments: ToId failed with err %s:", err)
			continue
		}

		if base62.Encode(domainId) != test.encodedFragOnlyId {
			t.Errorf("TestCreateIdsWithFragments: ToId failed expected %s: got %s", test.encodedFragOnlyId, base62.Encode(domainId))
		}

		domain := domain.Wrap(domainId)

		var idversion version.Version

		switch domain.VersionType() {
		case versiontype.NUMERIC:
			idversion = version.NumericVersion(rand.Int63())
		case versiontype.SEMANTIC:
			idversion = &version.SemanticVersion{
				Major: uint64(rand.Int63n(50)),
				Minor: uint64(rand.Int63n(50)),
				Patch: uint64(rand.Int63n(100))}
		}

		idlen := int(rand.Int63n(24))
		idbytes := random.RandomBytes(idlen)

		id, err := identifier.New(domain, idbytes, idversion)

		if err != nil {
			t.Errorf("TestCreateIdsWithFragments: identifier.New failed with err %s:", err)
			continue
		}

		if !id.IsValid() {
			t.Errorf("TestCreateIdsWithFragments: identifier.New failed id not valid:")
		}

		if !bytes.Equal(idbytes, id.Id()) {
			t.Errorf("TestCreateIdsWithFragments: id.Id %v: got %v", idbytes, id.Id())
		}

		if !bytes.Equal(id.IdRoot(), id.Id()) {
			t.Errorf("TestCreateIdsWithFragments: id.Id/id.IdRoot %v: got %v", id.IdRoot(), id.Id())
		}

		if len(id.Path()) > 0 {
			t.Errorf("TestCreateIdsWithFragments: id.Path len failed expected %d: got %d", 0, len(id.Path()))
		}

		if len(id.Fragment()) > 0 {
			t.Errorf("TestCreateIdsWithFragments: id.Fragment len failed expected %d: got %d", 0, len(id.Fragment()))
		}

		fraglen := int(rand.Int63n(24))
		fragbytes := random.RandomBytes(fraglen)

		id, err = identifier.New(domain, idbytes, idversion, nil, fragbytes)

		if err != nil {
			t.Errorf("TestCreateIdsWithFragments: identifier.New failed with err %s:", err)
			continue
		}

		if !id.IsValid() {
			t.Errorf("TestCreateIdsWithFragments: identifier.New failed id not valid:")
		}

		if !bytes.Equal(idbytes, id.IdRoot()) {
			t.Errorf("TestCreateIdsWithFragments: id.IdRoot expected: %v: got %v", idbytes, id.IdRoot())
		}

		if !bytes.Equal(fragbytes, id.Fragment()) {
			t.Errorf("TestCreateIdsWithFragments: id.Fragment expected: %v: got %v", fragbytes, id.Fragment())
		}

		if len(id.Path()) > 0 {
			t.Errorf("TestCreateIdsWithFragments: id.Path len failed expected %d: got %d", 0, len(id.Path()))
		}

		if idversion == nil {
			id, err = identifier.New(domain, idbytes, nil, fragbytes)

			if err != nil {
				t.Errorf("TestCreateIdsWithFragments: identifier.New failed with err %s:", err)
				continue
			}

			if !id.IsValid() {
				t.Errorf("TestCreateIdsWithFragments: identifier.New failed id not valid:")
			}

			if !bytes.Equal(idbytes, id.IdRoot()) {
				t.Errorf("TestCreateIdsWithFragments: id.IdRoot expected: %v: got %v", idbytes, id.IdRoot())
				fmt.Println(id.Value())
				fmt.Println(fragbytes)
			}

			if !bytes.Equal(fragbytes, id.Fragment()) {
				t.Errorf("TestCreateIdsWithFragments: id.Fragment expected: %v: got %v", fragbytes, id.Fragment())
			}

			if len(id.Path()) > 0 {
				t.Errorf("TestCreateIdsWithFragments: id.Path len failed expected %d: got %d", 0, len(id.Path()))
			}
		}
	}
}

func TestCreateIdsWithPathsAndFragments(t *testing.T) {
	for _, test := range testDomains {
		domainId, err := domain.ToId(test.schemeId, test.idRoot, test.incarnation, test.crcLength, test.versionType, true, true)

		if err != nil {
			t.Errorf("TestCreateIdsWithPathsAndFragments: ToId failed with err %s:", err)
			continue
		}

		if base62.Encode(domainId) != test.encodedPathAndFragId {
			t.Errorf("TestCreateIdsWithPathsAndFragments: ToId failed expected %s: got %s", test.encodedPathAndFragId, base62.Encode(domainId))
		}

		domain := domain.Wrap(domainId)

		var idversion version.Version

		switch domain.VersionType() {
		case versiontype.NUMERIC:
			idversion = version.NumericVersion(rand.Int63())
		case versiontype.SEMANTIC:
			idversion = &version.SemanticVersion{
				Major: uint64(rand.Int63n(50)),
				Minor: uint64(rand.Int63n(50)),
				Patch: uint64(rand.Int63n(100))}
		}

		idlen := int(rand.Int63n(24))
		idbytes := random.RandomBytes(idlen)

		id, err := identifier.New(domain, idbytes, idversion)

		if err != nil {
			t.Errorf("TestCreateIdsWithPathsAndFragments: identifier.New failed with err %s:", err)
			continue
		}

		if !id.IsValid() {
			t.Errorf("TestCreateIdsWithPathsAndFragments: identifier.New failed id not valid:")
		}

		if !bytes.Equal(idbytes, id.Id()) {
			t.Errorf("TestCreateIdsWithPathsAndFragments: id.Id %v: got %v", idbytes, id.Id())
		}

		if !bytes.Equal(id.IdRoot(), id.Id()) {
			t.Errorf("TestCreateIdsWithPathsAndFragments: id.Id/id.IdRoot %v: got %v", id.IdRoot(), id.Id())
		}

		if len(id.Path()) > 0 {
			t.Errorf("TestCreateIdsWithPathsAndFragments: id.Path len failed expected %d: got %d", 0, len(id.Path()))
		}

		if len(id.Fragment()) > 0 {
			t.Errorf("TestCreateIdsWithPathsAndFragments: id.Fragment len failed expected %d: got %d", 0, len(id.Fragment()))
		}

		pathlen := int(rand.Int63n(24))
		pathbytes := random.RandomBytes(pathlen)
		fraglen := int(rand.Int63n(24))
		fragbytes := random.RandomBytes(fraglen)

		id, err = identifier.New(domain, idbytes, idversion, pathbytes, fragbytes)

		if err != nil {
			t.Errorf("TestCreateIdsWithPathsAndFragments: identifier.New failed with err %s:", err)
			continue
		}

		if !id.IsValid() {
			t.Errorf("TestCreateIdsWithPathsAndFragments: identifier.New failed id not valid:")
		}

		if !bytes.Equal(idbytes, id.IdRoot()) {
			t.Errorf("TestCreateIdsWithPathsAndFragments: id.IdRoot expected: %v: got %v", idbytes, id.IdRoot())
		}

		if !bytes.Equal(fragbytes, id.Fragment()) {
			t.Errorf("TestCreateIdsWithPathsAndFragments: id.Fragment expected: %v: got %v", fragbytes, id.Fragment())
		}

		if !bytes.Equal(pathbytes, id.Path()) {
			t.Errorf("TestCreateIdsWithPathsAndFragments: id.Path %v: got %v", pathbytes, id.Path())
		}

		if idversion == nil {
			id, err = identifier.New(domain, idbytes, pathbytes, fragbytes)

			if err != nil {
				t.Errorf("TestCreateIdsWithPathsAndFragments: identifier.New failed with err %s:", err)
				continue
			}

			if !id.IsValid() {
				t.Errorf("TestCreateIdsWithPathsAndFragments: identifier.New failed id not valid:")
			}

			if !bytes.Equal(idbytes, id.IdRoot()) {
				t.Errorf("TestCreateIdsWithFragments: id.IdRoot expected: %v: got %v", idbytes, id.IdRoot())
				fmt.Println(id.Value())
				fmt.Println(fragbytes)
			}

			if !bytes.Equal(fragbytes, id.Fragment()) {
				t.Errorf("TestCreateIdsWithPathsAndFragments: id.Fragment expected: %v: got %v", fragbytes, id.Fragment())
			}

			if !bytes.Equal(pathbytes, id.Path()) {
				t.Errorf("TestCreateIdsWithPathsAndFragments: id.Path %v: got %v", pathbytes, id.Path())
			}
		}
	}
}
