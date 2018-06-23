package identitydomain_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/identitydomain"
	"github.com/distributed-vision/go-resources/ids/scheme"
	"github.com/distributed-vision/go-resources/init/idsinit"
)

func TestMain(m *testing.M) {
	//fmt.Printf("init...\n")
	idsinit.Init()
	//fmt.Printf("run...\n")
	os.Exit(m.Run())
}

type hashDomainTest struct {
	schemeId string
	idRoot   string
	name     string
	domainId string
}

var hashDomainTests = []hashDomainTest{
	hashDomainTest{"1", "0", "xx64", "H7A"},
	hashDomainTest{"1", "1", "xx32", "H7B"},
	hashDomainTest{"1", "2", "ripemd160", "H7C"},
	hashDomainTest{"1", "3", "sha256", "H7D"},
}

func TestDomainGet(t *testing.T) {

	for _, test := range hashDomainTests {

		schemeId, err := base62.Decode(test.schemeId)

		if err != nil {
			t.Errorf("TestDomainGet Failed: Can't decode schemeId")
		}

		domainId, err := base62.Decode(test.idRoot)

		if err != nil {
			t.Errorf("TestDomainGet Failed: Can't decode domainId")
		}

		schemeNameSelector := scheme.Selector{Name: "global hash", Opts: scheme.SelectorOpts{
			IgnoreCase:       true,
			IgnoreWhitespace: true}}

		domainNameSelector := domain.Selector{Name: test.name, Opts: domain.SelectorOpts{
			IgnoreCase:       true,
			IgnoreWhitespace: true}}

		idScheme, err := scheme.Get(context.Background(), scheme.Selector{Id: schemeId})

		if err != nil {
			t.Errorf("TestDomainGet Failed to resolve Scheme: %s", err)
		} else {
			if base62.Encode(idScheme.Id()) != test.schemeId {
				t.Errorf("TestDomainGet Failed: scheme.Id: expected: '%s' got '%s'", test.schemeId, base62.Encode(idScheme.Id()))
			}

			domain, derr := domain.Get(context.Background(), domain.Selector{SchemeId: schemeId, IdRoot: domainId})

			if derr != nil {
				t.Errorf("TestDomainGet Failed to resolve Domain: %s", derr)
			} else {
				if base62.Encode(domain.Id()) != test.domainId {
					t.Errorf("TestDomainGet Failed: domain.Id: expected: '%s' got '%s'", test.domainId, base62.Encode(domain.Id()))
				}

				if domain.Name() != test.name {
					t.Errorf("TestDomainGet Failed: domain.Name: expected: '%s' got '%s'", test.name, domain.Name())
				}

				if base62.Encode(domain.IdRoot()) != test.idRoot {
					t.Errorf("TestDomainGet Failed: domain.IdRoot: expected: '%s' got '%s'", test.idRoot, base62.Encode(domain.IdRoot()))

				}

				if base62.Encode(domain.Scheme().Id()) != test.schemeId {
					t.Errorf("TestDomainGet Failed: domain.Scheme.Id: expected: '%s' got '%s'", test.schemeId, base62.Encode(domain.Scheme().Id()))
				}
			}
		}

		idScheme, err = scheme.Get(context.Background(), schemeNameSelector)

		if err != nil {
			t.Errorf("TestDomainGet Failed to resolve Scheme: %s", err)
		} else {
			if base62.Encode(idScheme.Id()) != test.schemeId {
				t.Errorf("TestDomainGet Failed: scheme.Id: expected: '%s' got '%s'", test.schemeId, base62.Encode(idScheme.Id()))
			}

			domain, derr := domain.Get(context.Background(), domainNameSelector)

			if derr != nil {
				t.Errorf("TestDomainGet Failed to resolve Domain: %s", derr)
			} else {
				if base62.Encode(domain.Id()) != test.domainId {
					t.Errorf("TestDomainGet Failed: domain.Id: expected: '%s' got '%s'", test.domainId, base62.Encode(domain.Id()))
				}

				if domain.Name() != test.name {
					t.Errorf("TestDomainGet Failed: domain.Name: expected: '%s' got '%s'", test.name, domain.Name())

				}

				if base62.Encode(domain.IdRoot()) != test.idRoot {
					t.Errorf("TestDomainGet Failed: domain.IdRoot: expected: '%s' got '%s'", test.idRoot, base62.Encode(domain.IdRoot()))

				}

				if base62.Encode(domain.Scheme().Id()) != test.schemeId {
					t.Errorf("TestDomainGet Failed: domain.Scheme.Id: expected: '%s' got '%s'", test.schemeId, base62.Encode(domain.Scheme().Id()))
				}
			}
		}
	}
}

func TestDomainWithIncarnations(t *testing.T) {

	schemeId, err := base62.Decode("1")

	if err != nil {
		t.Errorf("TestDomainWithIncarnations Failed: Can't decode schemeId")
	}

	domainId, err := base62.Decode("3")

	if err != nil {
		t.Errorf("TestDomainWithIncarnations Failed: Can't decode domainId")
	}

	scheme, err := scheme.Get(context.Background(), scheme.Selector{Id: schemeId})

	if err != nil {
		t.Errorf("TestDomainWithIncarnations Failed to resolve Scheme: %s", err)
	} else {
		if !bytes.Equal(scheme.Id(), schemeId) {
			t.Errorf("TestDomainWithIncarnations Failed: scheme.Id: expected: '%s' got '%s'", schemeId, base62.Encode(scheme.Id()))
		}

		domain, derr := domain.Get(context.Background(), domain.Selector{SchemeId: schemeId, IdRoot: domainId})

		if derr != nil {
			t.Errorf("TestDomainWithIncarnations Failed to resolve Domain: %s", derr)
		} else {
			incarnation0, err := identitydomain.WithIncarnation(domain, 0, 0)

			if err != nil {
				t.Errorf("TestDomainWithIncarnations Failed to create incarnation0: %s", derr)
			}

			if base62.Encode(incarnation0.Id()) != "5tbcmW" {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation0.Id: expected: '%s' got '%s'", "5tbcmW", base62.Encode(incarnation0.Id()))
			}

			if base62.Encode(incarnation0.IdRoot()) != "3" {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation0.IdRoot: expected: '%s' got '%s'", "3", base62.Encode(incarnation0.IdRoot()))
			}

			if incarnation0.Incarnation() == nil || *incarnation0.Incarnation() != 0 {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation0.Incarnation: expected: '%v' got '%v'", 0, incarnation0.Incarnation())
			}

			if incarnation0.Name() != "sha256" {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation0.Name: expected: '%s' got '%s'", "sha256", incarnation0.Name())
			}

			incarnation300, err := identitydomain.WithIncarnation(domain, 300, 0)

			if err != nil {
				t.Errorf("TestDomainWithIncarnations Failed to create incarnation300: %s", derr)
			}

			if base62.Encode(incarnation300.Id()) != "OQJ9m7I" {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300.Id: expected: '%s' got '%s'", "OQJ9m7I", base62.Encode(incarnation300.Id()))
			}

			if base62.Encode(incarnation300.IdRoot()) != "3" {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300.IdRoot: expected: '%s' got '%s'", "3", base62.Encode(incarnation300.IdRoot()))
			}

			if incarnation300 == nil || *incarnation300.Incarnation() != 300 {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300.Incarnation: expected: '%v' got '%v'", 300, incarnation300.Incarnation())
			}

			if incarnation300.Name() != "sha256" {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300.Name: expected: '%s' got '%s'", "sha256", incarnation300.Name())
			}

			incarnation300000, err := identitydomain.WithIncarnation(domain, 300000, 0)

			if err != nil {
				t.Errorf("TestDomainWithIncarnations Failed to create incarnation300000: %s", derr)
			}

			if base62.Encode(incarnation300000.Id()) != "6kzTrVNgSO" {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300000.Id: expected: '%s' got '%s'", "6kzTrVNgSO", base62.Encode(incarnation300000.Id()))
			}

			if base62.Encode(incarnation300000.IdRoot()) != "3" {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300000.IdRoot: expected: '%s' got '%s'", "3", base62.Encode(incarnation300000.IdRoot()))
			}

			if incarnation300000.Incarnation() == nil || *incarnation300000.Incarnation() != 300000 {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300000.Incarnation: expected: '%v' got '%v'", 300000, incarnation300000.Incarnation())
			}

			if incarnation300000.Name() != "sha256" {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300000.Name: expected: '%s' got '%s'", "sha256", incarnation300000.Name())
			}

			if !domain.IsRoot() {
				t.Errorf("TestDomainWithIncarnations Failed: !domain.IsRoot")
			}

			if incarnation0.IsRoot() {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300 IsRoot")
			}

			if incarnation300.IsRoot() {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300 IsRoot")
			}

			if incarnation300000.IsRoot() {
				t.Errorf("TestDomainWithIncarnations Failed: incarnation300000 IsRoot")
			}

			if !domain.IsRootOf(incarnation0) {
				t.Errorf("TestDomainWithIncarnations Failed: !domain.IsRootOf incarnation0")
			}

			if !domain.IsRootOf(incarnation300) {
				t.Errorf("TestDomainWithIncarnations Failed: !domain.IsRootOf incarnation300")
			}

			if !domain.IsRootOf(incarnation300000) {
				t.Errorf("TestDomainWithIncarnations Failed: !domain.IsRootOf incarnation300000")
			}

			if !domain.Matches(incarnation0) {
				t.Errorf("TestDomainWithIncarnations Failed: !domain.Matches incarnation0")
			}

			if !domain.Matches(incarnation300) {
				t.Errorf("TestDomainWithIncarnations Failed: !domain.Matches incarnation300")
			}

			if !domain.Matches(incarnation300000) {
				t.Errorf("TestDomainWithIncarnations Failed: !domain.Matches incarnation300000")
			}

			if !incarnation0.Matches(domain) {
				t.Errorf("TestDomainWithIncarnations Failed: !incarnation0.Matches domain")
			}

			if !incarnation0.Matches(incarnation300) {
				t.Errorf("TestDomainWithIncarnations Failed: !incarnation0.Matches incarnation300")
			}

			if !incarnation0.Matches(incarnation300000) {
				t.Errorf("TestDomainWithIncarnations Failed: !incarnation0.Matches incarnation300000")
			}

			if !incarnation300.Matches(domain) {
				t.Errorf("TestDomainWithIncarnations Failed: !incarnation300.Matches domain")
			}

			if !incarnation300.Matches(incarnation0) {
				t.Errorf("TestDomainWithIncarnations Failed: !incarnation300.Matches incarnation0")
			}

			if !incarnation300.Matches(incarnation300000) {
				t.Errorf("TestDomainWithIncarnations Failed: !incarnation300.Matches incarnation300000")
			}

			if !incarnation300000.Matches(domain) {
				t.Errorf("TestDomainWithIncarnations Failed: !incarnation300000.Matches domain")
			}

			if !incarnation300000.Matches(incarnation0) {
				t.Errorf("TestDomainWithIncarnations Failed: !incarnation300000.Matches incarnation0")
			}

			if !incarnation300000.Matches(incarnation300) {
				t.Errorf("TestDomainWithIncarnations Failed: !incarnation300000.Matches incarnation300")
			}
		}
	}
}

func TestDomainWithCrc(t *testing.T) {
	schemeId, err := base62.Decode("1")

	if err != nil {
		t.Errorf("TestDomainWithCrc Failed: Can't decode schemeId: %s", err)
	}

	domainId, err := base62.Decode("2")

	if err != nil {
		t.Errorf("TestDomainWithCrc Failed: Can't decode domainId: %s", err)
	}

	scheme, err := scheme.Get(context.Background(), scheme.Selector{Id: schemeId})

	if err != nil {
		t.Errorf("TestDomainWithCrc Failed to resolve Scheme: %s", err)
	} else {
		if !bytes.Equal(scheme.Id(), schemeId) {
			t.Errorf("TestDomainWithCrc Failed: scheme.Id: expected: '%s' got '%s'", schemeId, base62.Encode(scheme.Id()))
		}

		domain, derr := domain.Get(context.Background(), domain.Selector{SchemeId: schemeId, IdRoot: domainId})

		if derr != nil {
			t.Errorf("TestDomainWithCrc Failed to resolve Domain: %s", derr)
		} else {
			crc8, err := identitydomain.WithCrc(domain, 8)

			if err != nil {
				t.Errorf("TestDomainWithCrc Failed: Can't create domain with CRC: %s", err)
			}

			if base62.Encode(crc8.Id()) != "1QGxu" {
				t.Errorf("TestDomainWithCrc Failed: crc8.Id: expected: '%s' got '%s'", "1QGxu", base62.Encode(crc8.Id()))
			}
			if base62.Encode(crc8.IdRoot()) != "2" {
				t.Errorf("TestDomainWithCrc Failed: crc8.IdRoot: expected: '%s' got '%s'", "2", base62.Encode(crc8.IdRoot()))
			}
			if crc8.CrcLength() != 8 {
				t.Errorf("TestDomainWithCrc Failed: crc8.CrcLength: expected: '%v' got '%v'", 8, crc8.CrcLength())
			}
			if crc8.Name() != "ripemd160" {
				t.Errorf("TestDomainWithCrc Failed: crc8.Name: expected: '%s' got '%s'", "ripemd160", crc8.Name())
			}

			crc16, err := identitydomain.WithCrc(domain, 16)

			if err != nil {
				t.Errorf("TestDomainWithCrc Failed: Can't create domain with CRC: %s", err)
			}

			if base62.Encode(crc16.Id()) != "1QHEQ" {
				t.Errorf("TestDomainWithCrc Failed: crc16.Id: expected: '%s' got '%s'", "1QHEQ", base62.Encode(crc16.Id()))
			}
			if base62.Encode(crc16.IdRoot()) != "2" {
				t.Errorf("TestDomainWithCrc Failed: crc16.IdRoot: expected: '%s' got '%s'", "2", base62.Encode(crc16.IdRoot()))
			}
			if crc16.CrcLength() != 16 {
				t.Errorf("TestDomainWithCrc Failed: crc16.CrcLength: expected: '%v' got '%v'", 16, crc16.CrcLength())
			}
			if crc16.Name() != "ripemd160" {
				t.Errorf("TestDomainWithCrc Failed: crc16.Name: expected: '%s' got '%s'", "ripemd160", crc16.Name())
			}

			crc32, err := identitydomain.WithCrc(domain, 32)

			if err != nil {
				t.Errorf("TestDomainWithCrc Failed: Can't create domain with CRC: %s", err)
			}

			if base62.Encode(crc32.Id()) != "1QHUw" {
				t.Errorf("TestDomainWithCrc Failed: crc32.Id: expected: '%s' got '%s'", "1QHUw", base62.Encode(crc32.Id()))
			}
			if base62.Encode(crc32.IdRoot()) != "2" {
				t.Errorf("TestDomainWithCrc Failed: crc32.IdRoot: expected: '%s' got '%s'", "2", base62.Encode(crc32.IdRoot()))
			}
			if crc32.CrcLength() != 32 {
				t.Errorf("TestDomainWithCrc Failed: crc32.CrcLength: expected: '%v' got '%v'", 32, crc32.CrcLength())
			}
			if crc32.Name() != "ripemd160" {
				t.Errorf("TestDomainWithCrc Failed: crc32.Name: expected: '%s' got '%s'", "ripemd160", crc32.Name())
			}

			if !domain.IsRoot() {
				t.Errorf("TestDomainWithCrc Failed: !domain.IsRoot")
			}

			if crc8.IsRoot() {
				t.Errorf("TestDomainWithCrc Failed: crc8 IsRoot")
			}
			if crc16.IsRoot() {
				t.Errorf("TestDomainWithCrc Failed: crc16 IsRoot")
			}
			if crc32.IsRoot() {
				t.Errorf("TestDomainWithCrc Failed: crc32 IsRoot")
			}

			if !domain.IsRootOf(crc8) {
				t.Errorf("TestDomainWithCrc Failed: !domain.IsRootOf crc8")
			}
			if !domain.IsRootOf(crc16) {
				t.Errorf("TestDomainWithCrc Failed: !domain.IsRootOf crc16")
			}
			if !domain.IsRootOf(crc32) {
				t.Errorf("TestDomainWithCrc Failed: !domain.IsRootOf crc32")
			}

			if !domain.Matches(crc8) {
				t.Errorf("TestDomainWithCrc Failed: !domain.Matches crc8")
			}
			if !domain.Matches(crc16) {
				t.Errorf("TestDomainWithCrc Failed: !domain.Matches crc16")
			}
			if !domain.Matches(crc32) {
				t.Errorf("TestDomainWithCrc Failed: !domain.Matches crc32")
			}

			if !crc8.Matches(domain) {
				t.Errorf("TestDomainWithCrc Failed: !crc8.Matches domain")
			}
			if !crc8.Matches(crc16) {
				t.Errorf("TestDomainWithCrc Failed: !crc8.Matches crc16")
			}
			if !crc8.Matches(crc32) {
				t.Errorf("TestDomainWithCrc Failed: !crc8.Matches crc16")
			}

			if !crc16.Matches(domain) {
				t.Errorf("TestDomainWithCrc Failed: !crc16.Matches domain")
			}
			if !crc16.Matches(crc8) {
				t.Errorf("TestDomainWithCrc Failed: !crc16.Matches crc8")
			}
			if !crc16.Matches(crc32) {
				t.Errorf("TestDomainWithCrc Failed: !crc16.Matches crc32")
			}

			if !crc32.Matches(domain) {
				t.Errorf("TestDomainWithCrc Failed: !crc32.Matches domain")
			}
			if !crc32.Matches(crc8) {
				t.Errorf("TestDomainWithCrc Failed: !crc32.Matches crc8")
			}
			if !crc32.Matches(crc32) {
				t.Errorf("TestDomainWithCrc Failed: !crc32.Matches crc32")
			}
		}
	}
}
