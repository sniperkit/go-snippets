package identifier_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/OneOfOne/xxhash"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/ids/identitydomain"
	"github.com/distributed-vision/go-resources/init/idsinit"
	"github.com/distributed-vision/go-resources/util/hton"
	"github.com/distributed-vision/go-resources/util/ntoh"
	"github.com/distributed-vision/go-resources/util/random"
	"github.com/howeyc/crc16"
	"github.com/sigurn/crc8"
)

func TestMain(m *testing.M) {
	idsinit.Init()
	os.Exit(m.Run())
}

func TestCreateXX64Id(t *testing.T) {

	domain, err := domain.Get(context.Background(), domain.Selector{Id: domain.MustDecodeId(encodertype.BASE62, "1", "0")})

	if err != nil {
		t.Errorf("TestCreateXX64Id: domain.Get failed: %s", err)
		return
	}

	value := random.RandomBytes(64)
	hash := xxhash.Checksum64S(value, 0xCAFEBABE)
	idbuf := make([]byte, 8)

	id, err := identifier.New(domain, hton.U64(idbuf, 0, hash))

	if err != nil {
		t.Errorf("TestCreateXX64Id: identifier.New failed: %s", err)
		return
	}

	if !bytes.Equal(id.DomainId(), domain.Id()) {
		t.Errorf("TestDomainWithCrc Failed: id.DomainId != domain.Id: expected: '%s' got '%s'", id.DomainId(), domain.Id())
	}

	if !bytes.Equal(id.Id(), idbuf) {
		t.Errorf("TestDomainWithCrc Failed: id.Id: expected: '%s' got '%s'", idbuf, id.Id())
	}

	if !bytes.Equal(id.SchemeId(), domain.SchemeId()) {
		t.Errorf("TestDomainWithCrc Failed: id.SchemeId != domain.SchemeId: expected: '%s' got '%s'", id.SchemeId(), domain.SchemeId())
	}

	if id.Checksum() != nil {
		t.Errorf("TestDomainWithCrc Failed: unexpected checksum")
	}
}

var crc8Table *crc8.Table = crc8.MakeTable(crc8.CRC8_MAXIM)
var crc16Table *crc16.Table = crc16.MakeTable(crc16.IBM)

func TestCreateXX64IdWithCrc(t *testing.T) {
	domain, err := domain.Get(context.Background(), domain.Selector{Id: domain.MustDecodeId(encodertype.BASE62, "1", "0")})

	if err != nil {
		t.Errorf("TestCreateXX64IdWithCrc: domain.Get failed: %s", err)
		return
	}

	crc8d, err := identitydomain.WithCrc(domain, 8)

	if err != nil {
		t.Errorf("TestCreateXX64IdWithCrc: identitydomain.WithCrc failed: %s", err)
		return
	}

	value := random.RandomBytes(64)
	hash := xxhash.Checksum64S(value, 0xCAFEBABE)
	idbuf := make([]byte, 8)

	id, err := identifier.New(domain, hton.U64(idbuf, 0, hash))

	if err != nil {
		t.Errorf("TestCreateXX64IdWithCrc: identifier.New failed: %s", err)
		return
	}

	crcid, err := identifier.New(crc8d, idbuf)

	if err != nil {
		t.Errorf("TestCreateXX64IdWithCrc: identifier.New failed: %s", err)
		return
	}

	cs := crc8.Checksum(crcid.Value()[:len(crcid.Value())-1], crc8Table)

	if !bytes.Equal(id.DomainId(), domain.Id()) {
		t.Errorf("TestCreateXX64IdWithCrc Failed: id.DomainId != domain.Id: expected: '%s' got '%s'", id.DomainId(), domain.Id())
	}
	if !bytes.Equal(id.Id(), idbuf) {
		t.Errorf("TestCreateXX64IdWithCrc Failed: id.Id: expected: '%s' got '%s'", idbuf, id.Id())
	}
	if !bytes.Equal(id.SchemeId(), domain.SchemeId()) {
		t.Errorf("TestCreateXX64IdWithCrc Failed: id.SchemeId != domain.SchemeId: expected: '%s' got '%s'", id.SchemeId(), domain.SchemeId())
	}
	if id.Checksum() != nil {
		t.Errorf("TestCreateXX64IdWithCrc Failed: unexpected checksum")
	}
	if !id.IsValid() {
		t.Errorf("TestCreateXX64IdWithCrc Failed: invalid id")
	}

	if !bytes.Equal(crcid.DomainId(), crc8d.Id()) {
		t.Errorf("TestCreateXX64IdWithCrc Failed: crcid.DomainId != crc8d.Id: expected: '%s' got '%s'", crcid.DomainId(), crc8d.Id())
	}
	if !bytes.Equal(crcid.DomainIdRoot(), crc8d.IdRoot()) {
		t.Errorf("TestCreateXX64IdWithCrc Failed: crcid.DomainIdRoot != crc8d.IdRoot: expected: '%s' got '%s'", crcid.DomainIdRoot(), crc8d.IdRoot())
	}
	if !bytes.Equal(crcid.DomainIdRoot(), id.DomainIdRoot()) {
		t.Errorf("TestCreateXX64IdWithCrc Failed: crcid.DomainIdRoot != id.DomainIdRoot: expected: '%s' got '%s'", crcid.DomainIdRoot(), id.DomainIdRoot())
	}
	if !bytes.Equal(crcid.Id(), idbuf) {
		t.Errorf("TestCreateXX64IdWithCrc Failed: crcid.Id: expected: '%s' got '%s'", idbuf, crcid.Id())
	}

	if !bytes.Equal(crcid.SchemeId(), domain.SchemeId()) {
		t.Errorf("TestCreateXX64IdWithCrc Failed: crcid.SchemeId != domain.SchemeId: expected: '%s' got '%s'", crcid.SchemeId(), domain.SchemeId())
	}

	if ntoh.U8(crcid.Checksum(), 0) != cs {
		t.Errorf("TestCreateXX64IdWithCrc Failed: crcid.Checksum: expected: '%v' got '%v'", cs, crcid.Checksum())
	}

	if !crcid.IsValid() {
		t.Errorf("TestCreateXX64IdWithCrc Failed: invalid crcid")
	}

	// if !id.toString('base62', 'base62')).equal(crcid.toString('base62', 'base62')) {

	//}
}

func TestCreateXX64IdWithIncarnationAndCrc(t *testing.T) {
	domain, err := domain.Get(context.Background(), domain.Selector{Id: domain.MustDecodeId(encodertype.BASE62, "1", "0")})

	if err != nil {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc: domain.Get failed: %s", err)
		return
	}

	crc16_30, err := identitydomain.WithIncarnation(domain, 30, 16)

	value := random.RandomBytes(64)
	hash := xxhash.Checksum64S(value, 0xCAFEBABE)
	idbuf := make([]byte, 8)

	id, err := identifier.New(domain, hton.U64(idbuf, 0, hash))

	crcid, err := identifier.New(crc16_30, idbuf)

	cs := crc16.Checksum(crcid.Value()[:len(crcid.Value())-2], crc16Table)

	if !bytes.Equal(id.DomainId(), domain.Id()) {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: id.DomainId != domain.Id: expected: '%s' got '%s'", id.DomainId(), domain.Id())
	}
	if !bytes.Equal(id.Id(), idbuf) {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: id.Id: expected: '%s' got '%s'", idbuf, id.Id())
	}
	if !bytes.Equal(id.SchemeId(), domain.SchemeId()) {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: id.SchemeId != domain.SchemeId: expected: '%s' got '%s'", id.SchemeId(), domain.SchemeId())
	}
	if id.Checksum() != nil {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: unexpected checksum: %v", id.Checksum())
	}
	if !id.IsValid() {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: invalid id")
	}

	if !bytes.Equal(crcid.DomainId(), crc16_30.Id()) {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: crcid.DomainId != crc16_30.Id: expected: '%v' got '%v'", crcid.DomainId(), crc16_30.Id())
	}
	if !bytes.Equal(crcid.DomainIdRoot(), crc16_30.IdRoot()) {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: crcid.DomainIdRoot != crc16_30.IdRoot: expected: '%s' got '%s'", crcid.DomainIdRoot(), crc16_30.IdRoot())
	}
	if !bytes.Equal(crcid.DomainIdRoot(), id.DomainIdRoot()) {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: crcid.DomainIdRoot != id.DomainIdRoot: expected: '%s' got '%s'", crcid.DomainIdRoot(), id.DomainIdRoot())
	}
	if !bytes.Equal(crcid.Id(), idbuf) {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: crcid.Id: expected: '%s' got '%s'", idbuf, crcid.Id())
	}

	if !bytes.Equal(crcid.SchemeId(), domain.SchemeId()) {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: crcid.SchemeId != domain.SchemeId: expected: '%s' got '%s'", crcid.SchemeId(), domain.SchemeId())
	}
	if ntoh.U16(crcid.Checksum(), 0) != cs {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: crcid.Checksum: expected: '%v' got '%v'", cs, ntoh.U16(crcid.Checksum(), 0))
	}
	if !crcid.IsValid() {
		t.Errorf("TestCreateXX64IdWithIncarnationAndCrc Failed: invalid crcid")
	}

	// if !id.toString('base62', 'base62')).equal(crcid.toString('base62', 'base62')) {

	//}
}
