package ospf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"syscall"
)

const (
	Version   = 2  // OSPF Version2
	HeaderLen = 24 // header length without payload
)

var errorHeaderTooShort = errors.New("OSPF packet header is to short")

type Header struct {
	Version        byte
	Type           byte
	PacketLength   int
	RouterID       net.IP
	AreaID         net.IP
	CheckSum       uint16
	AuType         uint16
	Authentication uint64
}

func (h *Header) Marshal() ([]byte, error) {
	if h == nil {
		return nil, syscall.EINVAL
	}
	if h.PacketLength < HeaderLen {
		return nil, errorHeaderTooShort
	}

	b := make([]byte, HeaderLen)
	b[0] = h.Version
	b[1] = h.Type
	binary.BigEndian.PutUint16(b[2:4], uint16(h.PacketLength))
	if ip := h.RouterID.To4(); ip != nil {
		copy(b[4:8], ip[:net.IPv4len])
	}
	if ip := h.AreaID.To4(); ip != nil {
		copy(b[8:12], ip[:net.IPv4len])
	}
	binary.BigEndian.PutUint16(b[12:14], uint16(h.CheckSum))
	binary.BigEndian.PutUint16(b[14:16], uint16(h.AuType))
	binary.BigEndian.PutUint64(b[16:24], uint64(h.Authentication))

	return b, nil
}

func UnMarshalOSPFHeader(b []byte) (*Header, []byte, error) {
	if len(b) < HeaderLen {
		return nil, nil, errors.New("OSPF Header is too short")
	}

	o := &Header{}
	o.Version = b[0]
	o.Type = b[1]
	o.PacketLength = int(binary.BigEndian.Uint16(b[2:4]))
	o.RouterID = net.IPv4(b[4], b[5], b[6], b[7])
	o.AreaID = net.IPv4(b[8], b[9], b[10], b[11])
	o.CheckSum = binary.BigEndian.Uint16(b[12:14])
	o.AuType = binary.BigEndian.Uint16(b[14:16])
	o.Authentication = binary.BigEndian.Uint64(b[16:24])
	return o, b[HeaderLen:], nil
}

var TypeToName = map[byte]string{
	1: "Hello",
	2: "Database Description",
	3: "Link State Request",
	4: "Link State Update",
	5: "Link State Acknowledge",
}

func (h *Header) String() string {
	var s string
	s += fmt.Sprintf("OSPF Packet: \n")
	s += fmt.Sprintf("         Version                : %d \n", h.Version)
	s += fmt.Sprintf("         Packet Type            : (%d)%s \n", h.Type, TypeToName[h.Type])
	s += fmt.Sprintf("         Packet Length          : %d \n", h.PacketLength)
	s += fmt.Sprintf("         Router ID              : %s \n", h.RouterID)
	s += fmt.Sprintf("         Area ID                : %s \n", h.AreaID)
	s += fmt.Sprintf("         CheckSum               : 0x%x \n", h.CheckSum)
	s += fmt.Sprintf("         AuType                 : 0x%x \n", h.AuType)
	s += fmt.Sprintf("         Authentication         : 0x%x \n", h.Authentication)
	return s
}
