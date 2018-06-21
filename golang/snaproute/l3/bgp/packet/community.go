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

// bgp.go
package packet

import (
	"encoding/binary"
	"fmt"
)

type BGPPathAttrCommunity struct {
	BGPPathAttrBase
	Value []uint32
}

func (c *BGPPathAttrCommunity) Clone() BGPPathAttr {
	x := *c
	x.BGPPathAttrBase = c.BGPPathAttrBase.Clone()
	x.Value = make([]uint32, len(c.Value))
	copy(x.Value, c.Value)
	return &x
}

func (c *BGPPathAttrCommunity) Encode() ([]byte, error) {
	pkt, err := c.BGPPathAttrBase.Encode()
	if err != nil {
		return pkt, err
	}

	index := int(c.BGPPathAttrLen)
	for i, comm := range c.Value {
		binary.BigEndian.PutUint32(pkt[i*4+index:], comm)
	}
	return pkt, nil
}

func (c *BGPPathAttrCommunity) Decode(pkt []byte, data interface{}) error {
	err := c.BGPPathAttrBase.Decode(pkt, data)
	if err != nil {
		return err
	}

	if (c.Length % 4) != 0 {
		return BGPMessageError{BGPUpdateMsgError, BGPAttrLenError, pkt[:c.TotalLen()], "Bad Attribute Length"}
	}

	index := int(c.BGPPathAttrLen)
	num_comm := int(c.Length / 4)
	c.Value = make([]uint32, num_comm)
	for i := 0; i < num_comm; i++ {
		c.Value[i] = binary.BigEndian.Uint32(pkt[(i*4)+index:])
	}
	return nil
}

func (c *BGPPathAttrCommunity) New() BGPPathAttr {
	return &BGPPathAttrCommunity{}
}

func (c *BGPPathAttrCommunity) String() string {
	return fmt.Sprintf("{Community %v}", c.Value)
}

func (c *BGPPathAttrCommunity) AddCommunity(value uint32) {
	if c.Length <= 255 && (c.Length+4) > 255 {
		c.Flags = c.Flags | BGPPathAttrFlagExtendedLen
		c.BGPPathAttrLen = 4
	}
	c.Value = append(c.Value, 0)
	copy(c.Value[1:], c.Value[0:])
	c.Value[0] = value
	c.Length += 4
}

func NewBGPPathAttrCommunity() *BGPPathAttrCommunity {
	comm := &BGPPathAttrCommunity{
		BGPPathAttrBase: BGPPathAttrBase{
			Flags:          BGPPathAttrFlagOptional | BGPPathAttrFlagTransitive,
			Code:           BGPPathAttrTypeCommunity,
			BGPPathAttrLen: 3,
		},
		Value: make([]uint32, 0),
	}

	return comm
}

func AddCommunityToPathAttrs(pa []BGPPathAttr, value uint32) []BGPPathAttr {
	paCommunity := getTypeFromPathAttrs(pa, BGPPathAttrTypeCommunity)
	if paCommunity == nil {
		paCommunity = NewBGPPathAttrCommunity()
		pa = addPathAttrToPathAttrs(pa, paCommunity)
	}

	community := paCommunity.(*BGPPathAttrCommunity)
	community.AddCommunity(value)
	return pa
}

func GetCommunityValues(pa []BGPPathAttr) []uint32 {
	paCommunity := getTypeFromPathAttrs(pa, BGPPathAttrTypeCommunity)
	if paCommunity == nil {
		return nil
	}

	community := paCommunity.(*BGPPathAttrCommunity)
	return community.Value
}

type BGPPathAttrExtCommunity struct {
	BGPPathAttrBase
	Value []uint64
}

func (e *BGPPathAttrExtCommunity) Clone() BGPPathAttr {
	x := *e
	x.BGPPathAttrBase = e.BGPPathAttrBase.Clone()
	x.Value = make([]uint64, len(e.Value))
	copy(x.Value, e.Value)
	return &x
}

func (e *BGPPathAttrExtCommunity) Encode() ([]byte, error) {
	pkt, err := e.BGPPathAttrBase.Encode()
	if err != nil {
		return pkt, err
	}

	index := int(e.BGPPathAttrLen)
	for i, comm := range e.Value {
		binary.BigEndian.PutUint64(pkt[i*8+index:], comm)
	}
	return pkt, nil
}

func (e *BGPPathAttrExtCommunity) Decode(pkt []byte, data interface{}) error {
	err := e.BGPPathAttrBase.Decode(pkt, data)
	if err != nil {
		return err
	}

	if (e.Length % 8) != 0 {
		return BGPMessageError{BGPUpdateMsgError, BGPAttrLenError, pkt[:e.TotalLen()], "Bad Attribute Length"}
	}

	index := int(e.BGPPathAttrLen)
	num_comm := int(e.Length / 8)
	e.Value = make([]uint64, num_comm)
	for i := 0; i < num_comm; i++ {
		e.Value[i] = binary.BigEndian.Uint64(pkt[(i*8)+index:])
	}
	return nil
}

func (e *BGPPathAttrExtCommunity) New() BGPPathAttr {
	return &BGPPathAttrExtCommunity{}
}

func (e *BGPPathAttrExtCommunity) String() string {
	return fmt.Sprintf("{ExtCommunity %v}", e.Value)
}

func (e *BGPPathAttrExtCommunity) AddCommunity(value uint64) {
	if e.Length <= 255 && (e.Length+8) > 255 {
		e.Flags = e.Flags | BGPPathAttrFlagExtendedLen
		e.BGPPathAttrLen = 4
	}
	e.Value = append(e.Value, 0)
	copy(e.Value[1:], e.Value[0:])
	e.Value[0] = value
	e.Length += 8
}

func NewBGPPathAttrExtCommunity() *BGPPathAttrExtCommunity {
	comm := &BGPPathAttrExtCommunity{
		BGPPathAttrBase: BGPPathAttrBase{
			Flags:          BGPPathAttrFlagOptional | BGPPathAttrFlagTransitive,
			Code:           BGPPathAttrTypeExtCommunity,
			BGPPathAttrLen: 3,
		},
		Value: make([]uint64, 0),
	}

	return comm
}

func AddExtCommunityToPathAttrs(pa []BGPPathAttr, value uint64) []BGPPathAttr {
	paExtComm := getTypeFromPathAttrs(pa, BGPPathAttrTypeExtCommunity)
	if paExtComm == nil {
		paExtComm = NewBGPPathAttrExtCommunity()
		pa = addPathAttrToPathAttrs(pa, paExtComm)
	}

	extComm := paExtComm.(*BGPPathAttrExtCommunity)
	extComm.AddCommunity(value)
	return pa
}

func GetExtCommunityValues(pa []BGPPathAttr) []uint64 {
	paExtComm := getTypeFromPathAttrs(pa, BGPPathAttrTypeExtCommunity)
	if paExtComm == nil {
		return nil
	}

	extComm := paExtComm.(*BGPPathAttrExtCommunity)
	return extComm.Value
}
