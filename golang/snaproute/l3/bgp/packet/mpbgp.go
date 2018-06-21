//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

// mpbgp.go
package packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	_ "l3/bgp/utils"
	"net"
	_ "strconv"
)

var BGPAFIToStructMap = map[AFI]MPNextHop{
	AfiIP:  &MPNextHopIP{},
	AfiIP6: &MPNextHopIP6{},
}

type MPNextHop interface {
	Clone() MPNextHop
	Encode([]byte) error
	Decode(pkt []byte) error
	Len() uint8
	New() MPNextHop
	String() string
	GetNextHop() net.IP
}

type MPNextHopIP struct {
	Length uint8
	Value  net.IP
}

func (i *MPNextHopIP) Clone() MPNextHop {
	x := *i
	x.Value = make(net.IP, len(i.Value), cap(i.Value))
	copy(x.Value, i.Value)
	return &x
}

func (i *MPNextHopIP) Encode(pkt []byte) error {
	pkt[0] = i.Length
	if i.Length != 4 && i.Length != 16 && i.Length != 32 {
		return errors.New(fmt.Sprintf("Wrong Next hop len %d", i.Length))
	}
	ipLen := net.IPv4len
	if i.Length == 16 || i.Length == 32 {
		ipLen = net.IPv6len
	}
	copy(pkt[1:], i.Value[cap(i.Value)-ipLen:])
	return nil
}

func (i *MPNextHopIP) Decode(pkt []byte) error {
	i.Length = pkt[0]
	if i.Length != 4 && i.Length != 16 && i.Length != 32 {
		return errors.New(fmt.Sprintf("Wrong Next hop len %d", i.Length))
	}
	ipLen := net.IPv4len
	if i.Length == 16 || i.Length == 32 {
		ipLen = net.IPv6len
	}

	i.Value = make(net.IP, ipLen)
	if ipLen == 4 {
		i.Value = net.IPv4(pkt[1], pkt[2], pkt[3], pkt[4])
	} else {
		copy(i.Value, pkt[1:])
	}
	return nil
}

func (i *MPNextHopIP) Len() uint8 {
	return i.Length + 1
}

func (i *MPNextHopIP) New() MPNextHop {
	return &MPNextHopIP{}
}

func (i *MPNextHopIP) String() string {
	return fmt.Sprintf("{NEXTHOP %v}", i.Value)
}

func (i *MPNextHopIP) GetNextHop() net.IP {
	return i.Value
}

func (i *MPNextHopIP) SetNextHop(ip net.IP) error {
	if len(ip) != 4 && len(ip) != 16 {
		return errors.New(fmt.Sprintf("Next hop IP address is not 4 bytes or 16 bytes, length=%d", len(ip)))
	}

	if ip.To4() == nil && ip.To16() == nil {
		return errors.New(fmt.Sprintf("Next hop IP address is NOT IPv4 or IPv6 address, ip=%s", ip))
	}

	i.Value = ip
	i.Length = uint8(net.IPv6len)
	if ip.To4() != nil {
		i.Length = uint8(net.IPv4len)
	}
	return nil
}

func NewMPNextHopIP() *MPNextHopIP {
	return &MPNextHopIP{
		Length: 0,
		Value:  net.IP{},
	}
}

type MPNextHopIP6 struct {
	*MPNextHopIP
	LinkLocal net.IP
}

func (i *MPNextHopIP6) Clone() MPNextHop {
	x := *i
	nextHopIP := i.MPNextHopIP.Clone()
	x.MPNextHopIP = nextHopIP.(*MPNextHopIP)
	x.Value = make(net.IP, len(i.Value), cap(i.Value))
	copy(x.Value, i.Value)
	return &x
}

func (i *MPNextHopIP6) Encode(pkt []byte) error {
	err := i.MPNextHopIP.Encode(pkt)
	if err != nil {
		return err
	}
	if i.LinkLocal != nil {
		ipLen := net.IPv6len
		copy(pkt[ipLen+1:], i.LinkLocal[cap(i.LinkLocal)-ipLen:])
	}
	return nil
}

func (i *MPNextHopIP6) Decode(pkt []byte) error {
	err := i.MPNextHopIP.Decode(pkt)
	if err != nil {
		return err
	}

	if i.Length == 32 {
		ipLen := net.IPv6len
		i.LinkLocal = make(net.IP, ipLen)
		copy(i.LinkLocal, pkt[ipLen+1:])
	}
	return nil
}

func (i *MPNextHopIP6) New() MPNextHop {
	return &MPNextHopIP{}
}

func (i *MPNextHopIP6) String() string {
	return fmt.Sprintf("{NEXTHOP %v}", i.Value)
}

func (i *MPNextHopIP6) SetGlobalNextHop(ip net.IP) error {
	if len(ip) != 16 {
		return errors.New(fmt.Sprintf("IPv6 next hop address is not 16 bytes, length =%d", len(ip)))
	}

	i.Value = ip
	if i.LinkLocal != nil {
		i.Length = uint8(len(ip)) * 2
	} else {
		i.Length = uint8(len(ip))
	}
	return nil
}

func (i *MPNextHopIP6) SetLinkLocalNextHop(ip net.IP) error {
	if len(ip) != 16 {
		return errors.New(fmt.Sprintf("IPv6 next hop address is not 16 bytes, length =%d", len(ip)))
	}

	i.LinkLocal = ip
	if i.LinkLocal != nil {
		i.Length = uint8(len(ip)) * 2
	} else {
		i.Length = uint8(len(ip))
	}
	return nil
}

func NewMPNextHopIP6() *MPNextHopIP6 {
	return &MPNextHopIP6{
		MPNextHopIP: &MPNextHopIP{
			Length: 0,
			Value:  net.IP{},
		},
	}
}

type MPNextHopUnknown struct {
	Length uint8
	Value  []byte
}

func (u *MPNextHopUnknown) Clone() MPNextHop {
	x := *u
	x.Value = make(net.IP, len(u.Value), cap(u.Value))
	copy(x.Value, u.Value)
	return &x
}

func (u *MPNextHopUnknown) Encode(pkt []byte) error {
	pkt[0] = u.Length
	copy(pkt[1:], u.Value)
	return nil
}

func (u *MPNextHopUnknown) Decode(pkt []byte) error {
	u.Length = pkt[0]
	u.Value = make([]byte, u.Length)
	copy(u.Value, pkt[1:])
	return nil
}

func (i *MPNextHopUnknown) Len() uint8 {
	return i.Length + 1
}

func (u *MPNextHopUnknown) New() MPNextHop {
	return &MPNextHopUnknown{}
}

func (u *MPNextHopUnknown) String() string {
	return fmt.Sprintf("{NEXTHOP %v}", u.Value)
}

func (u *MPNextHopUnknown) GetNextHop() net.IP {
	return u.Value
}

func (u *MPNextHopUnknown) SetNextHop(nextHop []byte) error {
	u.Length = uint8(len(nextHop))
	u.Value = make([]byte, u.Length)
	copy(u.Value, nextHop)
	return nil
}

func NewMPNextHopUnknown() *MPNextHopUnknown {
	return &MPNextHopUnknown{
		Length: 0,
		Value:  []byte{},
	}
}

func BGPGetMPNextHop(afi AFI) MPNextHop {
	var nextHop MPNextHop
	var ok bool
	if nextHop, ok = BGPAFIToStructMap[afi]; ok {
		nextHop = nextHop.New()
	} else {
		nextHop = &MPNextHopUnknown{}
	}
	return nextHop
}

type BGPPathAttrMPReachNLRI struct {
	BGPPathAttrBase
	AFI      AFI
	SAFI     SAFI
	NextHop  MPNextHop
	Reserved byte
	NLRI     []NLRI
}

func (r *BGPPathAttrMPReachNLRI) Clone() BGPPathAttr {
	x := *r
	x.BGPPathAttrBase = r.BGPPathAttrBase.Clone()
	x.NextHop = r.NextHop.Clone()
	x.NLRI = make([]NLRI, len(r.NLRI))
	for i, _ := range r.NLRI {
		x.NLRI[i] = r.NLRI[i].Clone()
	}
	return &x
}

func (r *BGPPathAttrMPReachNLRI) Encode() ([]byte, error) {
	pkt, err := r.BGPPathAttrBase.Encode()
	if err != nil {
		return pkt, nil
	}
	idx := int(r.BGPPathAttrBase.BGPPathAttrLen)

	binary.BigEndian.PutUint16(pkt[idx:idx+2], uint16(r.AFI))
	idx += 2
	pkt[idx] = uint8(r.SAFI)
	idx++

	err = r.NextHop.Encode(pkt[idx:])
	if err != nil {
		return pkt, err
	}
	idx += int(r.NextHop.Len())

	pkt[idx] = 0
	idx++

	for i := 0; i < len(r.NLRI); i++ {
		bytes, err := r.NLRI[i].Encode(r.AFI)
		if err != nil {
			return pkt, err
		}
		copy(pkt[idx:], bytes)
		idx += len(bytes)
	}
	return pkt, nil
}

func (r *BGPPathAttrMPReachNLRI) Decode(pkt []byte, data interface{}) error {
	err := r.BGPPathAttrBase.Decode(pkt, data)
	if err != nil {
		return err
	}

	idx := int(r.BGPPathAttrBase.BGPPathAttrLen)
	r.AFI = AFI(binary.BigEndian.Uint16(pkt[idx : idx+2]))
	r.SAFI = SAFI(pkt[idx+2])
	idx += 3

	nextHop := BGPGetMPNextHop(r.AFI)
	nextHop.Decode(pkt[idx:])
	r.NextHop = nextHop
	idx += int(nextHop.Len())

	r.Reserved = pkt[idx]
	idx++

	r.NLRI = make([]NLRI, 0)
	length := uint32(r.BGPPathAttrBase.Length) - 4 - uint32(r.NextHop.Len())
	_, err = decodeNLRI(pkt[idx:], &r.NLRI, length, r.AFI, r.SAFI, data)
	return err
}

func (r *BGPPathAttrMPReachNLRI) New() BGPPathAttr {
	return &BGPPathAttrMPReachNLRI{}
}

func (r *BGPPathAttrMPReachNLRI) SetNextHop(nextHop MPNextHop) {
	r.NextHop = nextHop
	r.BGPPathAttrBase.Length += uint16(r.NextHop.Len())
}

func (r *BGPPathAttrMPReachNLRI) AddNLRI(nlri NLRI) {
	r.NLRI = append(r.NLRI, nlri)
	r.BGPPathAttrBase.Length += uint16(nlri.Len())
}

func (r *BGPPathAttrMPReachNLRI) SetNLRIList(nlriList []NLRI) {
	r.NLRI = nlriList
	for _, nlri := range nlriList {
		r.BGPPathAttrBase.Length += uint16(nlri.Len())
	}
}

func NewBGPPathAttrMPReachNLRI() *BGPPathAttrMPReachNLRI {
	return &BGPPathAttrMPReachNLRI{
		BGPPathAttrBase: BGPPathAttrBase{
			Flags:          BGPPathAttrFlagOptional | BGPPathAttrFlagExtendedLen,
			Code:           BGPPathAttrTypeMPReachNLRI,
			Length:         4,
			BGPPathAttrLen: 4,
		},
		Reserved: 0,
		NLRI:     make([]NLRI, 0),
	}
}

type BGPPathAttrMPUnreachNLRI struct {
	BGPPathAttrBase
	AFI  AFI
	SAFI SAFI
	NLRI []NLRI
}

func (u *BGPPathAttrMPUnreachNLRI) Clone() BGPPathAttr {
	x := *u
	x.BGPPathAttrBase = u.BGPPathAttrBase.Clone()
	x.NLRI = make([]NLRI, len(u.NLRI))
	for i, _ := range u.NLRI {
		x.NLRI[i] = u.NLRI[i].Clone()
	}
	return &x
}

func (u *BGPPathAttrMPUnreachNLRI) Encode() ([]byte, error) {
	pkt, err := u.BGPPathAttrBase.Encode()
	if err != nil {
		return pkt, nil
	}
	idx := int(u.BGPPathAttrBase.BGPPathAttrLen)

	binary.BigEndian.PutUint16(pkt[idx:idx+2], uint16(u.AFI))
	idx += 2
	pkt[idx] = uint8(u.SAFI)
	idx++

	for i := 0; i < len(u.NLRI); i++ {
		bytes, err := u.NLRI[i].Encode(u.AFI)
		if err != nil {
			return pkt, err
		}
		copy(pkt[idx:], bytes)
		idx += len(bytes)
	}
	return pkt, nil
}

func (u *BGPPathAttrMPUnreachNLRI) Decode(pkt []byte, data interface{}) error {
	err := u.BGPPathAttrBase.Decode(pkt, data)
	if err != nil {
		return err
	}

	idx := int(u.BGPPathAttrBase.BGPPathAttrLen)
	u.AFI = AFI(binary.BigEndian.Uint16(pkt[idx : idx+2]))
	u.SAFI = SAFI(pkt[idx+2])
	idx += 3

	u.NLRI = make([]NLRI, 0)
	length := uint32(u.BGPPathAttrBase.Length) - 3
	_, err = decodeNLRI(pkt[idx:], &u.NLRI, length, u.AFI, u.SAFI, data)
	return err
}

func (u *BGPPathAttrMPUnreachNLRI) New() BGPPathAttr {
	return &BGPPathAttrMPUnreachNLRI{}
}

func (u *BGPPathAttrMPUnreachNLRI) AddNLRI(nlri NLRI) {
	u.NLRI = append(u.NLRI, nlri)
	u.BGPPathAttrBase.Length += uint16(nlri.Len())
}

func (u *BGPPathAttrMPUnreachNLRI) AddNLRIList(nlriList []NLRI) {
	for _, nlri := range nlriList {
		u.NLRI = append(u.NLRI, nlri)
		u.BGPPathAttrBase.Length += uint16(nlri.Len())
	}
}

func NewBGPPathAttrMPUnreachNLRI() *BGPPathAttrMPUnreachNLRI {
	return &BGPPathAttrMPUnreachNLRI{
		BGPPathAttrBase: BGPPathAttrBase{
			Flags:          BGPPathAttrFlagOptional | BGPPathAttrFlagExtendedLen,
			Code:           BGPPathAttrTypeMPUnreachNLRI,
			Length:         3,
			BGPPathAttrLen: 4,
		},
		NLRI: make([]NLRI, 0),
	}
}
