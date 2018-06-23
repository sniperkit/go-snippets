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
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

type BfdSessionState int

const (
	STATE_ADMIN_DOWN BfdSessionState = 0
	STATE_DOWN       BfdSessionState = 1
	STATE_INIT       BfdSessionState = 2
	STATE_UP         BfdSessionState = 3
)

func (server *BFDServer) ConvertBfdSessionStateValToStr(state BfdSessionState) string {
	var stateStr string
	switch state {
	case STATE_ADMIN_DOWN:
		stateStr = "admin_down"
	case STATE_DOWN:
		stateStr = "down"
	case STATE_INIT:
		stateStr = "init"
	case STATE_UP:
		stateStr = "up"
	}
	return stateStr
}

type BfdSessionEvent int

const (
	REMOTE_DOWN       BfdSessionEvent = 1
	REMOTE_INIT       BfdSessionEvent = 2
	REMOTE_UP         BfdSessionEvent = 3
	TIMEOUT           BfdSessionEvent = 4
	ADMIN_DOWN        BfdSessionEvent = 5
	ADMIN_UP          BfdSessionEvent = 6
	REMOTE_ADMIN_DOWN BfdSessionEvent = 7
)

type BfdDiagnostic int

const (
	DIAG_NONE                 BfdDiagnostic = 0 // No Diagnostic
	DIAG_TIME_EXPIRED         BfdDiagnostic = 1 // Control Detection Time Expired
	DIAG_ECHO_FAILED          BfdDiagnostic = 2 // Echo Function Failed
	DIAG_NEIGHBOR_SIGNAL_DOWN BfdDiagnostic = 3 // Neighbor Signaled Session Down
	DIAG_FORWARD_PLANE_RESET  BfdDiagnostic = 4 // Forwarding Plane Reset
	DIAG_PATH_DOWN            BfdDiagnostic = 5 // Path Down
	DIAG_CONCAT_PATH_DOWN     BfdDiagnostic = 6 // Concatenated Path Down
	DIAG_ADMIN_DOWN           BfdDiagnostic = 7 // Administratively Down
	DIAG_REV_CONCAT_PATH_DOWN BfdDiagnostic = 8 // Reverse Concatenated Path Down
)

func (server *BFDServer) ConvertBfdSessionDiagValToStr(diag BfdDiagnostic) string {
	var diagStr string
	switch diag {
	case DIAG_NONE:
		diagStr = "None"
	case DIAG_TIME_EXPIRED:
		diagStr = "Control detection timer expired"
	case DIAG_ECHO_FAILED:
		diagStr = "Echo function failed"
	case DIAG_NEIGHBOR_SIGNAL_DOWN:
		diagStr = "Neighbor signaled session down"
	case DIAG_FORWARD_PLANE_RESET:
		diagStr = "Forwarding plane reset"
	case DIAG_PATH_DOWN:
		diagStr = "Path down"
	case DIAG_CONCAT_PATH_DOWN:
		diagStr = "Concatanated path down"
	case DIAG_ADMIN_DOWN:
		diagStr = "Administratively down"
	case DIAG_REV_CONCAT_PATH_DOWN:
		diagStr = "Reverse concatenated path down"
	}
	return diagStr
}

type BfdControlPacket struct {
	Version                   uint8
	Diagnostic                BfdDiagnostic
	State                     BfdSessionState
	Poll                      bool
	Final                     bool
	ControlPlaneIndependent   bool
	AuthPresent               bool
	Demand                    bool
	Multipoint                bool // Must always be false
	DetectMult                uint8
	MyDiscriminator           uint32
	YourDiscriminator         uint32
	DesiredMinTxInterval      time.Duration
	RequiredMinRxInterval     time.Duration
	RequiredMinEchoRxInterval time.Duration
	AuthHeader                *BfdAuthHeader
}

// Constants
const (
	DEFAULT_BFD_VERSION                   = 1
	DEFAULT_DETECT_MULTI                  = 3
	DEFAULT_DESIRED_MIN_TX_INTERVAL       = 250000
	DEFAULT_REQUIRED_MIN_RX_INTERVAL      = 250000
	DEFAULT_REQUIRED_MIN_ECHO_RX_INTERVAL = 0
	DEFAULT_CONTROL_PACKET_LEN            = 24
	DEST_PORT                             = 3784
	SRC_PORT                              = 49152
	DEST_PORT_LAG                         = 6784
	SRC_PORT_LAG                          = 49153
	STARTUP_TX_INTERVAL                   = 2000000
	STARTUP_RX_INTERVAL                   = 2000000
	TX_JITTER                             = 10 //Timer will be running at 0 to 10% less than TX_INTERVAL
)

// Flags in BFD Control packet
const (
	BFD_MP             = 0x01 // Multipoint
	BFD_DEMAND         = 0x02 // Demand mode
	BFD_AUTH_PRESENT   = 0x04 // Authentication present
	BFD_CP_INDEPENDENT = 0x08 // Control plane independent
	BFD_FINAL          = 0x10 // Final message, response to Poll
	BFD_POLL           = 0x20 // Poll message
)

/*
 * Create a control packet
 */
func (p *BfdControlPacket) CreateBfdControlPacket() ([]byte, error) {
	//var auth []byte
	//var err error
	var authLength uint8
	buf := bytes.NewBuffer([]uint8{})
	flags := uint8(0)
	length := uint8(DEFAULT_CONTROL_PACKET_LEN)

	binary.Write(buf, binary.BigEndian, (p.Version<<5 | (uint8(p.Diagnostic) & 0x1f)))

	if p.Poll {
		flags |= BFD_POLL
	}
	if p.Final {
		flags |= BFD_FINAL
	}
	if p.ControlPlaneIndependent {
		flags |= BFD_CP_INDEPENDENT
	}
	if p.AuthPresent && (p.AuthHeader != nil) {
		flags |= BFD_AUTH_PRESENT
		/*
			auth, err = p.AuthHeader.createBfdAuthHeader()
			if err != nil {
				return nil, err
			}
		*/
		authLength = p.AuthHeader.getBfdAuthenticationLength()
		length += authLength
		//length += uint8(len(auth))
	}
	if p.Demand {
		flags |= BFD_DEMAND
	}
	if p.Multipoint {
		flags |= BFD_MP
	}

	binary.Write(buf, binary.BigEndian, (uint8(p.State)<<6 | flags))
	binary.Write(buf, binary.BigEndian, p.DetectMult)
	binary.Write(buf, binary.BigEndian, length)

	binary.Write(buf, binary.BigEndian, p.MyDiscriminator)
	binary.Write(buf, binary.BigEndian, p.YourDiscriminator)
	binary.Write(buf, binary.BigEndian, uint32(p.DesiredMinTxInterval))
	binary.Write(buf, binary.BigEndian, uint32(p.RequiredMinRxInterval))
	binary.Write(buf, binary.BigEndian, uint32(p.RequiredMinEchoRxInterval))

	if authLength > 0 {
		binary.Write(buf, binary.BigEndian, p.AuthHeader.Type)
		binary.Write(buf, binary.BigEndian, authLength)
		binary.Write(buf, binary.BigEndian, p.AuthHeader.AuthKeyID)
		if p.AuthHeader.Type != BFD_AUTH_TYPE_SIMPLE {
			binary.Write(buf, binary.BigEndian, uint8(0))
			binary.Write(buf, binary.BigEndian, p.AuthHeader.SequenceNumber)
		}
		copiedBuf := bytes.NewBuffer(buf.Bytes())
		switch p.AuthHeader.Type {
		case BFD_AUTH_TYPE_SIMPLE:
			binary.Write(buf, binary.BigEndian, p.AuthHeader.AuthData)
		case BFD_AUTH_TYPE_KEYED_MD5, BFD_AUTH_TYPE_METICULOUS_MD5:
			var authData [16]byte
			binary.Write(copiedBuf, binary.BigEndian, p.AuthHeader.AuthData)
			authData = md5.Sum(copiedBuf.Bytes())
			binary.Write(buf, binary.BigEndian, authData)
			fmt.Println("MD5 sum ", authData)
		case BFD_AUTH_TYPE_KEYED_SHA1, BFD_AUTH_TYPE_METICULOUS_SHA1:
			var authData [20]byte
			binary.Write(copiedBuf, binary.BigEndian, p.AuthHeader.AuthData)
			authData = sha1.Sum(copiedBuf.Bytes())
			binary.Write(buf, binary.BigEndian, authData)
			fmt.Println("SHA1 sum ", authData)
		}
	}

	return buf.Bytes(), nil
}

/*
 * Decode the control packet
 */
func DecodeBfdControlPacket(data []byte) (*BfdControlPacket, error) {
	var err error
	packet := &BfdControlPacket{}

	packet.Version = uint8((data[0] & 0xE0) >> 5)
	packet.Diagnostic = BfdDiagnostic(data[0] & 0x1F)

	packet.State = BfdSessionState((data[1] & 0xD0) >> 6)

	// bit flags
	packet.Poll = (data[1]&0x20 != 0)
	packet.Final = (data[1]&0x10 != 0)
	packet.ControlPlaneIndependent = (data[1]&0x08 != 0)
	packet.AuthPresent = (data[1]&0x04 != 0)
	packet.Demand = (data[1]&0x02 != 0)
	packet.Multipoint = (data[1]&0x01 != 0)
	packet.DetectMult = uint8(data[2])

	length := uint8(data[3]) // No need to store this
	if uint8(len(data)) != length {
		err = errors.New("Packet length mis-match!")
		return nil, err
	}

	packet.MyDiscriminator = binary.BigEndian.Uint32(data[4:8])
	packet.YourDiscriminator = binary.BigEndian.Uint32(data[8:12])
	packet.DesiredMinTxInterval = time.Duration(binary.BigEndian.Uint32(data[12:16]))
	packet.RequiredMinRxInterval = time.Duration(binary.BigEndian.Uint32(data[16:20]))
	packet.RequiredMinEchoRxInterval = time.Duration(binary.BigEndian.Uint32(data[20:24]))

	if packet.AuthPresent {
		if len(data) > 24 {
			packet.AuthHeader, err = decodeBfdAuthHeader(data[24:])
		} else {
			err = errors.New("Header flag set, but packet too short!")
		}
	}

	return packet, err
}
