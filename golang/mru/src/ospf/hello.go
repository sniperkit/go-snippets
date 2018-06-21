package ospf

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Hello struct {
	Header                 Header
	NetworkMask            net.IPMask
	HelloInterval          uint16
	Options                uint8
	RtrPri                 uint8
	RouterDeadInterval     uint32
	DesignatedRouter       net.IP
	BackupDesignatedRouter net.IP
	Neighbors              []net.IP
}

func (h Hello) Marshal() ([]byte, error) {
	b := make([]byte, 20+len(h.Neighbors)*4)
	copy(b[0:4], h.NetworkMask)

	binary.BigEndian.PutUint16(b[4:6], uint16(h.HelloInterval))
	b[6] = h.Options
	b[7] = h.RtrPri
	binary.BigEndian.PutUint32(b[8:12], uint32(h.RouterDeadInterval))

	if ip := h.DesignatedRouter.To4(); ip != nil {
		copy(b[12:16], ip[:net.IPv4len])
	}

	if ip := h.BackupDesignatedRouter.To4(); ip != nil {
		copy(b[16:20], ip[:net.IPv4len])
	}

	for i, n := range h.Neighbors {
		if ip := n.To4(); ip != nil {
			copy(b[20+i*4:20+(i+1)*4], ip[:net.IPv4len])
		}
	}

	return b, nil
}

func UnMarshalHello(b []byte, length int) (*Hello, error) {
	if length < 20 {
		return nil, fmt.Errorf("OSPF Hello packet is too short with length: %d", length)
	}

	hello := &Hello{}
	hello.NetworkMask = net.IPv4Mask(b[0], b[1], b[2], b[3])
	hello.HelloInterval = binary.BigEndian.Uint16(b[4:6])
	hello.Options = b[6]
	hello.RtrPri = b[7]
	hello.RouterDeadInterval = binary.BigEndian.Uint32(b[8:12])
	hello.DesignatedRouter = net.IPv4(b[12], b[13], b[14], b[15])
	hello.BackupDesignatedRouter = net.IPv4(b[16], b[17], b[18], b[19])
	hello.Neighbors = make([]net.IP, 0, (length-20)/4)

	for i := 0; i < (length-20)/4; i++ {
		hello.Neighbors = append(hello.Neighbors, net.IPv4(b[20+i*4], b[20+i*4+1], b[20+i*4+2], b[20+i*4+3]))
	}

	return hello, nil
}

func (h *Hello) String() string {
	if h == nil {
		return ""
	}
	var s string
	s += fmt.Sprintf("OSPF Hello Packet: \n")
	s += fmt.Sprintf("         NetworkMask            : %s\n", net.IPv4(h.NetworkMask[0], h.NetworkMask[1], h.NetworkMask[2], h.NetworkMask[3]))
	s += fmt.Sprintf("         HelloInterval          : %d\n", h.HelloInterval)
	s += fmt.Sprintf("         Options                : 0x%x\n", h.Options)
	s += fmt.Sprintf("         RtrPri                 : %d\n", h.RtrPri)
	s += fmt.Sprintf("         RouterDeadInterval     : %d\n", h.RouterDeadInterval)
	s += fmt.Sprintf("         DesignatedRouter       : %s\n", h.DesignatedRouter)
	s += fmt.Sprintf("         BackupDesignatedRouter : %s\n", h.BackupDesignatedRouter)
	s += fmt.Sprintf("         Attached Neighbors:    : \n")
	for _, n := range h.Neighbors {
		s += fmt.Sprintf("             %s\n", n)
	}

	return s
}

func (h Hello) IsOSPF() bool {
	return true
}

func (h Hello) UnMarshal(data []byte) (interface{}, error) {
	return nil, nil
}
