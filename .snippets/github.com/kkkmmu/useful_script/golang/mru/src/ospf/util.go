package ospf

import (
	"fmt"
	"golang.org/x/net/ipv4"
)

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func PrettyIPv4Header(h *ipv4.Header) string {
	if h == nil {
		return "<nil>"
	}

	var s string
	s += fmt.Sprintf("IPv4 Header: \n")
	s += fmt.Sprintf("         Version                : %d \n", h.Version)
	s += fmt.Sprintf("         IP4H Length            : %d \n", h.Len)
	s += fmt.Sprintf("         TOS                    : %d \n", h.TOS)
	s += fmt.Sprintf("         Total Length           : %d \n", h.TotalLen)
	s += fmt.Sprintf("         ID                     : %d \n", h.ID)
	s += fmt.Sprintf("         Flags                  : 0x%x \n", h.Flags)
	s += fmt.Sprintf("         FragOff                : 0x%x \n", h.FragOff)
	s += fmt.Sprintf("         TTL                    : 0x%x \n", h.TTL)
	s += fmt.Sprintf("         Protocol               : %d \n", h.Protocol)
	s += fmt.Sprintf("         CheckSum               : 0x%x \n", h.Checksum)
	s += fmt.Sprintf("         From                   : %s \n", h.Src)
	s += fmt.Sprintf("         To                     : %s \n", h.Dst)
	s += fmt.Sprintf("         Options                : %v \n", h.Options)
	return s
}
