package ospf

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type LSU struct {
	Header    Header
	LSANumber uint32
	LSAs      []LSA
}

func UnMarshalLSU(b []byte, length int) (*LSU, error) {
	if length < 4 {
		return nil, errors.New("LSUPacket is too short")
	}

	var lsu = &LSU{}
	lsu.LSANumber = binary.BigEndian.Uint32(b[0:4])
	if lsu.LSANumber != uint32(0) {
		lsas, err := UnMarshalLSAFromLSU(b[4:], lsu.LSANumber, uint32(length-4))
		if err != nil {
			return nil, errors.New("Cannot parse LSAs from LSU")
		}
		lsu.LSAs = lsas
	}
	return lsu, nil
}

func (lsu *LSU) String() string {
	if lsu == nil {
		return ""
	}
	var s string
	s += fmt.Sprintf("OSPF LS Update Packet: \n")
	s += fmt.Sprintf("         LSA Count:             : %d\n", lsu.LSANumber)
	for _, l := range lsu.LSAs {
		switch l.(type) {
		case RouterLSA:
			r, _ := l.(RouterLSA)
			s += fmt.Sprintf("             %s\n", r)
		case NetworkLSA:
			r, _ := l.(NetworkLSA)
			s += fmt.Sprintf("             %s\n", r)
		case NetworkSummaryLSA:
			r, _ := l.(NetworkSummaryLSA)
			s += fmt.Sprintf("             %s\n", r)
		case ASBRSummaryLSA:
			r, _ := l.(ASBRSummaryLSA)
			s += fmt.Sprintf("             %s\n", r)
		case ExternalLSA:
			r, _ := l.(ExternalLSA)
			s += fmt.Sprintf("             %s\n", r)
		default:
			panic("Unknown LSA type")
		}
	}

	return s
}
