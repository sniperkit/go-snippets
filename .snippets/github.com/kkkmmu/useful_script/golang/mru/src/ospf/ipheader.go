// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ospf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/net/ipv4"
	"net"
	"syscall"
)

var errIPHeaderTooShort = errors.New("IPv4 header is too short")
var errIPHeaderMissingAddress = errors.New("IPv4 header with no address")
var errIPHeaderBufferTooShort = errors.New("IPv4 header buffer is too short")

// A Header represents an IPv4 header.
type IP4H struct {
	Version  int              // protocol version
	Len      int              // header length
	TOS      int              // type-of-service
	TotalLen int              // packet total length
	ID       int              // identification
	Flags    ipv4.HeaderFlags // flags
	FragOff  int              // fragment offset
	TTL      int              // time-to-live
	Protocol int              // next protocol
	CheckSum int              // checksum
	Src      net.IP           // source address
	Dst      net.IP           // destination address
	Options  []byte           // options, extension headers
}

func (h *IP4H) String() string {
	if h == nil {
		return "<nil>"
	}

	var s string
	s += fmt.Sprintf("IPv4 IP4H: \n")
	s += fmt.Sprintf("        Version       : %d \n", h.Version)
	s += fmt.Sprintf("        IP4H Length : %d \n", h.Len)
	s += fmt.Sprintf("        TOS           : %d \n", h.TOS)
	s += fmt.Sprintf("        Total Length  : %d \n", h.TotalLen)
	s += fmt.Sprintf("        ID            : %d \n", h.ID)
	s += fmt.Sprintf("        Flags         : 0x%x \n", h.Flags)
	s += fmt.Sprintf("        FragOff       : 0x%x \n", h.FragOff)
	s += fmt.Sprintf("        TTL           : 0x%x \n", h.TTL)
	s += fmt.Sprintf("        Protocol      : %d \n", h.Protocol)
	s += fmt.Sprintf("        CheckSum      : 0x%x \n", h.CheckSum)
	s += fmt.Sprintf("        From          : %s \n", h.Src)
	s += fmt.Sprintf("        To            : %s \n", h.Dst)
	s += fmt.Sprintf("        Options       : %v \n", h.Options)
	return s
	//return fmt.Sprintf("ver=%d hdrlen=%d tos=%#x totallen=%d id=%#x flags=%#x fragoff=%#x ttl=%d proto=%d cksum=%#x src=%v dst=%v", h.Version, h.Len, h.TOS, h.TotalLen, h.ID, h.Flags, h.FragOff, h.TTL, h.Protocol, h.CheckSum, h.Src, h.Dst)
}

// Marshal returns the binary encoding of h.
func (h *IP4H) Marshal() ([]byte, error) {
	if h == nil {
		return nil, syscall.EINVAL
	}
	if h.Len < ipv4.HeaderLen {
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++")
		return nil, errIPHeaderTooShort
	}
	hdrlen := ipv4.HeaderLen + len(h.Options)
	b := make([]byte, hdrlen)
	b[0] = byte(Version<<4 | (hdrlen >> 2 & 0x0f))
	b[1] = byte(h.TOS)
	flagsAndFragOff := (h.FragOff & 0x1fff) | int(h.Flags<<13)
	binary.BigEndian.PutUint16(b[2:4], uint16(h.TotalLen))
	binary.BigEndian.PutUint16(b[6:8], uint16(flagsAndFragOff))
	binary.BigEndian.PutUint16(b[4:6], uint16(h.ID))
	b[8] = byte(h.TTL)
	b[9] = byte(h.Protocol)
	binary.BigEndian.PutUint16(b[10:12], uint16(h.CheckSum))
	if ip := h.Src.To4(); ip != nil {
		copy(b[12:16], ip[:net.IPv4len])
	}
	if ip := h.Dst.To4(); ip != nil {
		copy(b[16:20], ip[:net.IPv4len])
	} else {
		return nil, errIPHeaderMissingAddress
	}
	if len(h.Options) > 0 {
		copy(b[ipv4.HeaderLen:], h.Options)
	}
	return b, nil
}

// Parse parses b as an IPv4 header and sotres the result in h.
func (h *IP4H) Parse(b []byte) error {
	if h == nil || len(b) < ipv4.HeaderLen {
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++")
		return errIPHeaderTooShort
	}
	hdrlen := int(b[0]&0x0f) << 2
	if hdrlen > len(b) {
		return errIPHeaderBufferTooShort
	}
	h.Version = int(b[0] >> 4)
	h.Len = hdrlen
	h.TOS = int(b[1])
	h.ID = int(binary.BigEndian.Uint16(b[4:6]))
	h.TTL = int(b[8])
	h.Protocol = int(b[9])
	h.CheckSum = int(binary.BigEndian.Uint16(b[10:12]))
	h.Src = net.IPv4(b[12], b[13], b[14], b[15])
	h.Dst = net.IPv4(b[16], b[17], b[18], b[19])
	h.TotalLen = int(binary.BigEndian.Uint16(b[2:4]))
	h.FragOff = int(binary.BigEndian.Uint16(b[6:8]))
	h.Flags = ipv4.HeaderFlags(h.FragOff&0xe000) >> 13
	h.FragOff = h.FragOff & 0x1fff
	optlen := hdrlen - ipv4.HeaderLen
	if optlen > 0 && len(b) >= hdrlen {
		if cap(h.Options) < optlen {
			h.Options = make([]byte, optlen)
		} else {
			h.Options = h.Options[:optlen]
		}
		copy(h.Options, b[ipv4.HeaderLen:hdrlen])
	}
	return nil
}

// ParseIP4H parses b as an IPv4 header.
func ParseIP4H(b []byte) (*IP4H, error) {
	h := new(IP4H)
	if err := h.Parse(b); err != nil {
		return nil, err
	}
	return h, nil
}
