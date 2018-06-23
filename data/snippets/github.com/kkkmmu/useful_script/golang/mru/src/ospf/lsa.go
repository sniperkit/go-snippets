package ospf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

const (
	Router         = 1
	Network        = 2
	NetworkSummary = 3
	ASBRSummary    = 4
	External       = 5
)

type LSA interface {
	IsLSA() bool
	Marshal() ([]byte, error)
}

type LSAIdentity struct {
	LSType            uint32
	LSID              net.IP
	AdvertisingRouter net.IP
}

type LSAHeader struct {
	LSAge             uint16
	Options           uint8
	LSType            uint8
	LinkStateID       net.IP
	AdvertisingRouter net.IP
	LSSequenceNumber  uint32
	LSCheckSum        uint16
	Length            uint16
}

func UnMarshalLSAHeader(b []byte) (*LSAHeader, error) {
	lsh := &LSAHeader{}
	lsh.LSAge = binary.BigEndian.Uint16(b[0:2])
	lsh.Options = b[2]
	lsh.LSType = b[3]
	lsh.LinkStateID = net.IPv4(b[4], b[5], b[6], b[7])
	lsh.AdvertisingRouter = net.IPv4(b[8], b[9], b[10], b[11])
	lsh.LSSequenceNumber = binary.BigEndian.Uint32(b[12:16])
	lsh.LSCheckSum = binary.BigEndian.Uint16(b[16:18])
	lsh.Length = binary.BigEndian.Uint16(b[18:20])

	return lsh, nil
}

var LSATypeToName = map[uint8]string{
	1: "Router LSA",
	2: "Network LSA",
	3: "Network Summary LSA",
	4: "ASBR Summary LSA",
	5: "AS External LSA",
}

func (lsh *LSAHeader) String() string {
	var s string
	s += fmt.Sprintf("----------------------------------------------------------------------\n")
	s += fmt.Sprintf("        Type:                   :%d(%s)\n", lsh.LSType, LSATypeToName[lsh.LSType])
	s += fmt.Sprintf("        LinkStateID             : %s\n", lsh.LinkStateID)
	s += fmt.Sprintf("        AdvertisingRouter       : %s\n", lsh.AdvertisingRouter)
	s += fmt.Sprintf("        Age                     : %d\n", lsh.LSAge)
	s += fmt.Sprintf("        LSCheckSum              : 0x%x\n", lsh.LSCheckSum)
	s += fmt.Sprintf("        LSSequenceNumber        : 0x%x\n", lsh.LSSequenceNumber)
	s += fmt.Sprintf("        Options                 : 0x%x\n", lsh.Options)
	s += fmt.Sprintf("        Length                  : %d\n", lsh.Length)

	return s
}

var vebValueToString = map[uint8]string{
	0: "Intra area normal router",
	1: "Area Boarder Router",
	2: "AS Boundary Router",
	3: "Area Boarder Router, AS Boundary Router",
	4: "Virtuallink Endpoint",
	5: "Virtuallink Endpoint, Area Boarder Router",
	6: "Virutallink Endpoint, AS Boundary Router",
	7: "Virtuallink Endpoint, Area Boarder Router, AS Boundary Router",
}

type RouterLSA struct {
	Header    LSAHeader
	VEB       uint8
	Reserved  uint8
	LinkCount uint16
	Links     []LinkState
}

func (rl RouterLSA) IsLSA() bool {
	return true
}

func (rl RouterLSA) Marshal() ([]byte, error) {
	return nil, nil
}

func (rl RouterLSA) String() string {
	var s string
	s += fmt.Sprintf("   Router LSA:\n")
	s += fmt.Sprintf("                 V|E|B           : 0x%x(%s)\n", rl.VEB, vebValueToString[rl.VEB])
	s += fmt.Sprintf("                 Reverved        : 0x%x\n", rl.Reserved)
	s += fmt.Sprintf("                 Link Count      : %d\n", rl.LinkCount)
	for _, l := range rl.Links {
		s += fmt.Sprintf("                                %s\n", l)
	}

	return s
}

type LinkState struct {
	LinkID      net.IP
	LinkData    net.IP
	Type        uint8
	NumberOfTOS uint8
	Metric      uint16
}

type NetworkLSA struct {
	Header         LSAHeader
	NetworkMask    net.IP
	AttachedRouter []net.IP
}

func (nl NetworkLSA) String() string {
	var s string

	s += fmt.Sprintf("Network LSA: \n")
	s += fmt.Sprintf("            Network Mask        : %s\n", nl.NetworkMask)
	s += fmt.Sprintf("            AttachedRouter      : \n")
	for _, ar := range nl.AttachedRouter {
		s += fmt.Sprintf("                                : %s\n", ar)
	}

	return s
}

func (nl NetworkLSA) IsLSA() bool {
	return true
}

func (nl NetworkLSA) Marshal() ([]byte, error) {
	return nil, nil
}

type NetworkSum struct {
	NetworkMask net.IP
	metric      uint32 /* 24 bit by standard */
	TOS         uint8
	TOSMetric   uint32 /* 24 bit by standard */
}

type NetworkSummaryLSA struct {
	Header   LSAHeader
	Networks []NetworkSum
}

func (nsl NetworkSummaryLSA) IsLSA() bool {
	return true
}

func (nsl NetworkSummaryLSA) Marshal() ([]byte, error) {
	return nil, nil
}

type ASBRSum struct {
	NetworkMask net.IP
	metric      uint32 /* 24 bit by standard */
	TOS         uint8
	TOSMetric   uint32 /* 24 bit by standard */
}

type ASBRSummaryLSA struct {
	Header LSAHeader
	ASBRs  []ASBRSum
}

func (abl ASBRSummaryLSA) IsLSA() bool {
	return true
}

func (abl ASBRSummaryLSA) Marshal() ([]byte, error) {
	return nil, nil
}

type ExternalNetwork struct {
	NetworkMask       net.IP
	EBit              bool
	Metric            uint32 /* 24 bit by standard */
	ForWardingAddress net.IP
	ExternalRouteTag  uint32
}

type ExternalLSA struct {
	Header   LSAHeader
	Networks []ExternalNetwork
}

func (asel ExternalLSA) IsLSA() bool {
	return true
}

func (asel ExternalLSA) Marshal() ([]byte, error) {
	return nil, nil
}

func UnMarshalLSAFromLSU(b []byte, count, length uint32) ([]LSA, error) {
	if length < 20 {
		return nil, errors.New("LSU's LSA field is too short")
	}

	var off = 0
	lsas := make([]LSA, 0, count)
	for i := 0; i < int(count); i++ {
		lsah, err := UnMarshalLSAHeader(b)
		if err != nil {
			log.Fatal("Cannot parse LSA header")
		}

		switch lsah.LSType {
		case Router:
			lsa, err := UnMarshalRouterLSA(b[off:off+int(lsah.Length)], lsah)
			if err != nil {
				log.Fatal("Cannot Parase Router LSA")
				continue
			}
			lsa.Header = *lsah
			lsas = append(lsas, *lsa)
		case Network:
			lsa, err := UnMarshalNetworkLSA(b[off:off+int(lsah.Length)], lsah)
			if err != nil {
				log.Fatal("Cannot Parase Network LSA")
				continue
			}
			lsa.Header = *lsah
			lsas = append(lsas, *lsa)
		case NetworkSummary:
			lsa, err := UnMarshalNetworkSummaryLSA(b[off:off+int(lsah.Length)], lsah)
			if err != nil {
				log.Fatal("Cannot Parase Network Summary LSA")
				continue
			}
			lsa.Header = *lsah
			lsas = append(lsas, *lsa)
		case ASBRSummary:
			lsa, err := UnMarshalASBRSummaryLSA(b[off:off+int(lsah.Length)], lsah)
			if err != nil {
				log.Fatal("Cannot Parase ASBR Summary LSA")
				continue
			}
			lsa.Header = *lsah
			lsas = append(lsas, *lsa)
		case External:
			lsa, err := UnMarshalExternalLSA(b[off:off+int(lsah.Length)], lsah)
			if err != nil {
				log.Fatal("Cannot Parase External LSA")
				continue
			}
			lsa.Header = *lsah
			lsas = append(lsas, *lsa)
		default:
			panic("Unknow LSA type")
		}
		off += int(lsah.Length)
	}

	return lsas, nil
}

func UnMarshalLinkState(b []byte) (*LinkState, error) {
	var ls = &LinkState{}
	ls.LinkID = net.IPv4(b[0], b[1], b[2], b[3])
	ls.LinkData = net.IPv4(b[4], b[5], b[6], b[7])
	ls.Type = b[8]
	ls.NumberOfTOS = b[9]
	ls.Metric = binary.BigEndian.Uint16(b[10:12])
	return ls, nil
}

var linkTypeToName = map[uint8]string{
	1: "Point-to-Point connect to another router",
	2: "Connection to a transit network",
	3: "Connection to a stub network",
	4: "Virtual link",
}

func (ls LinkState) String() string {
	var s string
	s += fmt.Sprintf("Link %s: ID(%s) Data(%s) Type(%d:%s) Metric(%d)\n", ls.LinkID, ls.LinkID, ls.LinkData, ls.Type, linkTypeToName[ls.Type], ls.Metric)
	return s
}

func UnMarshalRouterLSA(b []byte, lsah *LSAHeader) (*RouterLSA, error) {
	var rl = &RouterLSA{}
	rl.VEB = b[20]
	rl.Reserved = b[21]
	rl.LinkCount = binary.BigEndian.Uint16(b[22:24])

	count := (lsah.Length - 24) / 12
	if count > rl.LinkCount {
		count = rl.LinkCount
	}

	if count != 0 {
		rl.Links = make([]LinkState, 0, count)
		for i := 0; i < int(count); i++ {
			if ls, err := UnMarshalLinkState(b[24+i*12 : 24+(i+1)*12]); err != nil {
				log.Fatal("Cannot Decode link state from router LSA")
				continue
			} else {
				rl.Links = append(rl.Links, *ls)
			}
		}
	}
	return rl, nil
}

func UnMarshalNetworkLSA(b []byte, lsah *LSAHeader) (*NetworkLSA, error) {
	var nl = &NetworkLSA{}
	atrcount := (lsah.Length - 24) / 4
	nl.NetworkMask = net.IPv4(b[20], b[21], b[22], b[23])
	nl.AttachedRouter = make([]net.IP, 0, atrcount)

	for i := 0; i < int(atrcount); i++ {
		nl.AttachedRouter = append(nl.AttachedRouter, net.IPv4(b[24+i*4], b[25+i*4], b[26+i*4], b[27+i*4]))
	}
	return nl, nil
}

func UnMarshalNetworkSummaryLSA(b []byte, lsah *LSAHeader) (*NetworkSummaryLSA, error) {
	return nil, errors.New("Cannot parse Network Summary lsa")
}

func UnMarshalASBRSummaryLSA(b []byte, lsah *LSAHeader) (*ASBRSummaryLSA, error) {
	return nil, errors.New("Cannot parse ASBR Summary lsa")
}

func UnMarshalExternalLSA(b []byte, lsah *LSAHeader) (*ExternalLSA, error) {
	return nil, errors.New("Cannot parse Exteral lsa")
}
