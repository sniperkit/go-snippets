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

package server

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type BfdAuthHeader struct {
	Type           AuthenticationType
	AuthKeyID      uint8
	SequenceNumber uint32
	AuthData       []byte
}

type AuthenticationType uint8

const (
	BFD_AUTH_TYPE_RESERVED        AuthenticationType = 0 // Reserved
	BFD_AUTH_TYPE_SIMPLE          AuthenticationType = 1 // Simple Password
	BFD_AUTH_TYPE_KEYED_MD5       AuthenticationType = 2 // Keyed MD5
	BFD_AUTH_TYPE_METICULOUS_MD5  AuthenticationType = 3 // Meticulous Keyed MD5
	BFD_AUTH_TYPE_KEYED_SHA1      AuthenticationType = 4 // Keyed SHA1
	BFD_AUTH_TYPE_METICULOUS_SHA1 AuthenticationType = 5 // Meticulous Keyed SHA1
)

func (server *BFDServer) ConvertBfdAuthTypeStrToVal(authType string) AuthenticationType {
	var authVal AuthenticationType
	switch authType {
	case "simple":
		authVal = BFD_AUTH_TYPE_SIMPLE
	case "keyedmd5":
		authVal = BFD_AUTH_TYPE_KEYED_MD5
	case "metmd5":
		authVal = BFD_AUTH_TYPE_METICULOUS_MD5
	case "keyedsha1":
		authVal = BFD_AUTH_TYPE_KEYED_SHA1
	case "metsha1":
		authVal = BFD_AUTH_TYPE_METICULOUS_SHA1
	}
	return authVal
}

func (server *BFDServer) ConvertBfdAuthTypeValToStr(authType AuthenticationType) string {
	var authStr string
	switch authType {
	case BFD_AUTH_TYPE_SIMPLE:
		authStr = "simple"
	case BFD_AUTH_TYPE_KEYED_MD5:
		authStr = "keyedmd5"
	case BFD_AUTH_TYPE_METICULOUS_MD5:
		authStr = "metmd5"
	case BFD_AUTH_TYPE_KEYED_SHA1:
		authStr = "keyedsha1"
	case BFD_AUTH_TYPE_METICULOUS_SHA1:
		authStr = "metsha1"
	}
	return authStr
}

/*
 * Create the Auth header section
 */
func (h *BfdAuthHeader) createBfdAuthHeader() ([]byte, error) {
	buf := bytes.NewBuffer([]uint8{})
	var length uint8

	if h.Type != BFD_AUTH_TYPE_SIMPLE {
		length = uint8(len(h.AuthData) + 8)
	} else {
		length = uint8(len(h.AuthData) + 3)
	}

	binary.Write(buf, binary.BigEndian, h.Type)
	binary.Write(buf, binary.BigEndian, length)
	binary.Write(buf, binary.BigEndian, h.AuthKeyID)

	if h.Type != BFD_AUTH_TYPE_SIMPLE {
		binary.Write(buf, binary.BigEndian, uint8(0))
		binary.Write(buf, binary.BigEndian, h.SequenceNumber)
	}

	binary.Write(buf, binary.BigEndian, h.AuthData)

	return buf.Bytes(), nil
}

func (h *BfdAuthHeader) getBfdAuthenticationLength() uint8 {
	var length uint8
	switch h.Type {
	case BFD_AUTH_TYPE_SIMPLE:
		length = uint8(len(h.AuthData) + 3)
	case BFD_AUTH_TYPE_KEYED_MD5, BFD_AUTH_TYPE_METICULOUS_MD5:
		length = uint8(24)
	case BFD_AUTH_TYPE_KEYED_SHA1, BFD_AUTH_TYPE_METICULOUS_SHA1:
		length = uint8(28)
	}
	return length
}

/*
 * Decode the Auth header section
 */
func decodeBfdAuthHeader(data []byte) (*BfdAuthHeader, error) {
	var err error
	h := &BfdAuthHeader{}

	h.Type = AuthenticationType(data[0])
	length := uint8(data[1])

	if length > 0 {
		h.AuthKeyID = uint8(data[2])

		switch h.Type {
		case BFD_AUTH_TYPE_SIMPLE:
			h.AuthData = data[3:]
			break
		case BFD_AUTH_TYPE_KEYED_MD5, BFD_AUTH_TYPE_METICULOUS_MD5:
			h.SequenceNumber = binary.BigEndian.Uint32(data[4:8])
			h.AuthData = data[8:]
			if len(h.AuthData) != 16 {
				err = errors.New("Invalid MD5 Auth Key/Digest length!")
			}
		case BFD_AUTH_TYPE_KEYED_SHA1, BFD_AUTH_TYPE_METICULOUS_SHA1:
			h.SequenceNumber = binary.BigEndian.Uint32(data[4:8])
			h.AuthData = data[8:]
			if len(h.AuthData) != 20 {
				err = errors.New("Invalid SHA1 Auth Key/Hash length!")
			}
		default:
			err = errors.New("Unsupported Authentication type!")
		}
	}

	if err != nil {
		return nil, err
	}

	return h, nil
}
