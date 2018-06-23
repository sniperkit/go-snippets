//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package dhcp6

import (
	"dhcprelayd"
	//	"errors"
	"fmt"
	"golang.org/x/net/ipv6"
	"l3/dhcp_relay/infra"
	"net"
	"strconv"
	"sync"
	"time"
	"utils/logging"
)

var Logger logging.LoggerIntf

type Processor struct {
	Logger        logging.LoggerIntf
	InfraMgr      *infra.InfraMgr
	ClientHandler *net.UDPConn
	ClientConn    *ipv6.PacketConn

	PeerAddrIntfMap map[string]string

	// States
	ClientStateSlice     []*dhcprelayd.DHCPv6RelayClientState
	ClientStateMap       map[string]*dhcprelayd.DHCPv6RelayClientState
	IntfStateSlice       []*dhcprelayd.DHCPv6RelayIntfState
	IntfStateMap         map[string]*dhcprelayd.DHCPv6RelayIntfState
	IntfServerStateSlice []*dhcprelayd.DHCPv6RelayIntfServerState
	IntfServerStateMap   map[string]*dhcprelayd.DHCPv6RelayIntfServerState
	StateMutex           sync.Mutex

	EnabledFlag  bool
	EnabledMutex sync.Mutex
}

type ProcessorInitParams struct {
	Logger   logging.LoggerIntf
	InfraMgr *infra.InfraMgr
}

func NewProcessor(initParams *ProcessorInitParams) *Processor {
	Logger = initParams.Logger
	pProc := &Processor{
		Logger:   initParams.Logger,
		InfraMgr: initParams.InfraMgr,
	}

	pProc.PeerAddrIntfMap = make(map[string]string)
	pProc.ClientStateSlice = []*dhcprelayd.DHCPv6RelayClientState{}
	pProc.ClientStateMap = make(map[string]*dhcprelayd.DHCPv6RelayClientState)
	pProc.IntfStateSlice = []*dhcprelayd.DHCPv6RelayIntfState{}
	pProc.IntfStateMap = make(map[string]*dhcprelayd.DHCPv6RelayIntfState)
	pProc.IntfServerStateSlice = []*dhcprelayd.DHCPv6RelayIntfServerState{}
	pProc.IntfServerStateMap = make(map[string]*dhcprelayd.DHCPv6RelayIntfServerState)

	pProc.EnabledMutex.Lock()
	pProc.EnabledFlag = false
	pProc.EnabledMutex.Unlock()

	return pProc
}

func (pProc *Processor) GetEnabledFlag() bool {
	defer pProc.EnabledMutex.Unlock()
	pProc.EnabledMutex.Lock()
	return pProc.EnabledFlag
}

// Don't use this ever (Used for mock testing)
func (pProc *Processor) SetEnabledFlag() {
	defer pProc.EnabledMutex.Unlock()
	pProc.EnabledMutex.Lock()
	pProc.EnabledFlag = true
}

func (pProc *Processor) GetClientStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPv6RelayClientState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.ClientStateSlice)
	if fromIdx < 0 || fromIdx >= length || count <= 0 {
		return 0, 0, false, []*dhcprelayd.DHCPv6RelayClientState{}
	}
	if fromIdx+count >= length {
		actualCount = length - fromIdx
		nextIdx = 0
		more = false
	} else {
		actualCount = count
		nextIdx = fromIdx + count
		more = true
	}

	result := make([]*dhcprelayd.DHCPv6RelayClientState, actualCount)
	copy(result, pProc.ClientStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor) GetClientState(
	macAddr string) (*dhcprelayd.DHCPv6RelayClientState, bool) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	if clientState, ok := pProc.ClientStateMap[macAddr]; ok {
		return clientState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor) GetIntfStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPv6RelayIntfState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.IntfStateSlice)
	if fromIdx < 0 || fromIdx >= length || count <= 0 {
		return 0, 0, false, []*dhcprelayd.DHCPv6RelayIntfState{}
	}
	if fromIdx+count >= length {
		actualCount = length - fromIdx
		nextIdx = 0
		more = false
	} else {
		actualCount = count
		nextIdx = fromIdx + count
		more = true
	}

	result := make([]*dhcprelayd.DHCPv6RelayIntfState, actualCount)
	copy(result, pProc.IntfStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor) GetIntfState(
	ifName string) (*dhcprelayd.DHCPv6RelayIntfState, bool) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	if intfState, ok := pProc.IntfStateMap[ifName]; ok {
		return intfState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor) GetIntfServerStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPv6RelayIntfServerState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.IntfServerStateSlice)
	if fromIdx < 0 || fromIdx >= length || count <= 0 {
		return 0, 0, false, []*dhcprelayd.DHCPv6RelayIntfServerState{}
	}
	if fromIdx+count >= length {
		actualCount = length - fromIdx
		nextIdx = 0
		more = false
	} else {
		actualCount = count
		nextIdx = fromIdx + count
		more = true
	}

	result := make([]*dhcprelayd.DHCPv6RelayIntfServerState, actualCount)
	copy(result, pProc.IntfServerStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor) GetIntfServerState(
	ifName string, serverAddr string) (*dhcprelayd.DHCPv6RelayIntfServerState, bool) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfServerStateKey := ifName + "_" + serverAddr
	if intfServerState, ok := pProc.IntfServerStateMap[intfServerStateKey]; ok {
		return intfServerState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor) StartRxTx() {
	Logger.Debug("pProc: StartRxTx being called")

	defer pProc.EnabledMutex.Unlock()
	pProc.EnabledMutex.Lock()

	if pProc.EnabledFlag {
		return
	}
	if !pProc.CreateConn() {
		return
	}
	go pProc.RxTx()
	pProc.EnabledFlag = true
}

func (pProc *Processor) StopRxTx() {
	Logger.Debug("pProc: StopRxTx being called")

	defer pProc.EnabledMutex.Unlock()
	pProc.EnabledMutex.Lock()
	// Close UDP socket connection of Rx thread
	if pProc.ClientConn != nil {
		err := pProc.ClientConn.Close()
		if err != nil {
			Logger.Err("Error while closing socket")
		}
		pProc.ClientConn = nil
	}
	pProc.EnabledFlag = false
}

func (pProc *Processor) initClientState(
	clientAddr string) *dhcprelayd.DHCPv6RelayClientState {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	clientStateKey := clientAddr
	clientState, ok := pProc.ClientStateMap[clientStateKey]
	if !ok {
		clientState = &dhcprelayd.DHCPv6RelayClientState{}
		clientState.MacAddr = clientAddr
		pProc.ClientStateSlice = append(pProc.ClientStateSlice, clientState)
		pProc.ClientStateMap[clientStateKey] = clientState
	}
	return clientState
}

func (pProc *Processor) initIntfState(
	ifName string) *dhcprelayd.DHCPv6RelayIntfState {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfState, ok := pProc.IntfStateMap[ifName]
	if !ok {
		intfState = &dhcprelayd.DHCPv6RelayIntfState{}
		intfState.IntfRef = ifName
		pProc.IntfStateSlice = append(pProc.IntfStateSlice, intfState)
		pProc.IntfStateMap[ifName] = intfState
	}
	return intfState
}

func (pProc *Processor) deleteIntfState(
	ifName string) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	_, ok := pProc.IntfStateMap[ifName]
	if !ok {
		return
	}
	sliceEntIdx := -1
	for i, intfState := range pProc.IntfStateSlice {
		if intfState.IntfRef == ifName {
			sliceEntIdx = i
			break
		}
	}
	if sliceEntIdx != -1 {
		pProc.IntfStateSlice = append(pProc.IntfStateSlice[:sliceEntIdx],
			pProc.IntfStateSlice[sliceEntIdx+1:]...)
	}
	delete(pProc.IntfStateMap, ifName)
}

func (pProc *Processor) initIntfServerState(
	ifName string, serverAddr string) *dhcprelayd.DHCPv6RelayIntfServerState {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfServerStateKey := ifName + "_" + serverAddr
	intfServerState, ok := pProc.IntfServerStateMap[intfServerStateKey]
	if !ok {
		intfServerState = &dhcprelayd.DHCPv6RelayIntfServerState{}
		intfServerState.IntfRef = ifName
		intfServerState.ServerIp = serverAddr
		pProc.IntfServerStateSlice = append(
			pProc.IntfServerStateSlice, intfServerState)
		pProc.IntfServerStateMap[intfServerStateKey] = intfServerState
	}
	return intfServerState
}

func (pProc *Processor) setDownstreamInState(
	mType MsgType,
	dhcpOptions DhcpOptionMap,
	clientState *dhcprelayd.DHCPv6RelayClientState,
	intfState *dhcprelayd.DHCPv6RelayIntfState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfState.TotalDhcpServerRx++
	if mType == RELAY_REPL {
		return
	}
	clientState.ServerResponses++

	switch mType {
	case ADVERTISE:
		if dOpt, ok := dhcpOptions[OPTION_IA_NA]; ok {
			ipAddrField := OptionIANA(dOpt).GetIpAddr()
			if ipAddrField == nil {
				return
			}
			clientState.OfferedIp = net.IP(ipAddrField).String()
		} else if dOpt, ok := dhcpOptions[OPTION_IA_TA]; ok {
			ipAddrField := OptionIATA(dOpt).GetIpAddr()
			if ipAddrField == nil {
				return
			}
			clientState.OfferedIp = net.IP(ipAddrField).String()
		}
		clientState.ServerAdvertise = time.Now().String()
	case REPLY:
		if dOpt, ok := dhcpOptions[OPTION_IA_NA]; ok {
			ipAddrField := OptionIANA(dOpt).GetIpAddr()
			if ipAddrField == nil {
				return
			}
			clientState.AcceptedIp = net.IP(ipAddrField).String()
		} else if dOpt, ok := dhcpOptions[OPTION_IA_TA]; ok {
			ipAddrField := OptionIATA(dOpt).GetIpAddr()
			if ipAddrField == nil {
				return
			}
			clientState.AcceptedIp = net.IP(ipAddrField).String()
		}
		clientState.ServerReply = time.Now().String()
	case RECONFIGURE:
		clientState.ServerReconfigure = time.Now().String()
	}
}

// Downstream: Server -> Client
func (pProc *Processor) SendPktDownstream(
	outPkt []byte, dhcpOptions DhcpOptionMap, peerAddr net.IP, srcAddr net.IP,
	clientState *dhcprelayd.DHCPv6RelayClientState) {

	mType := MsgType(outPkt[0])
	peerAddrString := peerAddr.String()
	var peerPort int
	if mType == RELAY_REPL {
		peerPort = DHCP_SERVER_PORT
	} else {
		peerPort = DHCP_CLIENT_PORT
	}

	outIfName, ok := pProc.PeerAddrIntfMap[peerAddrString]
	if !ok {
		Logger.Debug("DRA: No out interface found for peer Ip", peerAddrString)
		return
	}
	outIfIdx, ok := pProc.InfraMgr.GetIPv6IntfIndex(outIfName)
	if !ok {
		Logger.Debug("DRA: No IPv6 interface found for", outIfIdx)
		return
	}

	intfState := pProc.initIntfState(outIfName)
	intfServerState := pProc.initIntfServerState(outIfName, srcAddr.String())
	if mType != RELAY_REPL {
		pProc.setDownstreamInState(mType, dhcpOptions, clientState, intfState)
	} else {
		pProc.setDownstreamInState(mType, nil, nil, intfState)
	}

	var destIpPortString string
	if peerAddr.IsLinkLocalUnicast() {
		destIpPortString = fmt.Sprintf(
			"[%s%%%s]:%s", peerAddrString, outIfName, strconv.Itoa(peerPort),
		)
	} else {
		destIpPortString = fmt.Sprintf(
			"[%s]:%s", peerAddrString, strconv.Itoa(peerPort),
		)
	}

	Logger.Debug("DRA: destIpPortString", destIpPortString)
	pProc.StateMutex.Lock()
	if pProc.SendPkt(outPkt, destIpPortString) {
		intfState.TotalDhcpClientTx++
		intfServerState.Responses++
		if mType != RELAY_REPL {
			clientState.ClientResponses++
		}
	} else {
		intfState.TotalDrops++
	}
	pProc.StateMutex.Unlock()
}

func (pProc *Processor) setUpstreamInState(
	mType MsgType,
	dhcpOptions DhcpOptionMap,
	intfState *dhcprelayd.DHCPv6RelayIntfState,
	clientState *dhcprelayd.DHCPv6RelayClientState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfState.TotalDhcpClientRx++
	if mType == RELAY_FORW {
		return
	}
	clientState.ClientRequests++

	switch mType {
	case SOLICIT:
		clientState.ClientSolicit = time.Now().String()
	case REQUEST:
		if dOpt, ok := dhcpOptions[OPTION_IA_NA]; ok {
			ipAddrField := OptionIANA(dOpt).GetIpAddr()
			if ipAddrField == nil {
				return
			}
			clientState.RequestedIp = net.IP(ipAddrField).String()
		} else if dOpt, ok := dhcpOptions[OPTION_IA_TA]; ok {
			ipAddrField := OptionIATA(dOpt).GetIpAddr()
			if ipAddrField == nil {
				return
			}
			clientState.RequestedIp = net.IP(ipAddrField).String()
		}
		clientState.ClientRequest = time.Now().String()
	case CONFIRM:
		clientState.ClientConfirm = time.Now().String()
	case RENEW:
		clientState.ClientRenew = time.Now().String()
	case REBIND:
		clientState.ClientRebind = time.Now().String()
	case DECLINE:
		clientState.ClientDecline = time.Now().String()
	case RELEASE:
		clientState.ClientRelease = time.Now().String()
	case INFO_REQ:
		clientState.ClientInfoRequest = time.Now().String()
	}
}

func (pProc *Processor) SendPkt(outPkt []byte, destIpPort string) bool {
	destAddr, _ := net.ResolveUDPAddr("udp6", destIpPort)
	_, err := pProc.ClientHandler.WriteToUDP(outPkt, destAddr)
	if err != nil {
		Logger.Debug("DRA: WriteToUDP failed with error:", err)
		return false
	}
	return true
}

// Upstream: Client -> Server
func (pProc *Processor) SendPktUpstream(
	outPkt []byte, mType MsgType, inIfIdx int, inIfName string, srcAddr net.IP,
	intfState *dhcprelayd.DHCPv6RelayIntfState,
	clientState *dhcprelayd.DHCPv6RelayClientState) {

	Logger.Debug("Sending packet upstream")
	draIntf, ok := pProc.InfraMgr.GetActiveDRAv6Intf(inIfIdx)
	if !ok {
		pProc.StateMutex.Lock() // State obj lock
		intfState.TotalDrops++
		pProc.StateMutex.Unlock() // State obj unlock
		Logger.Debug(
			"Dropping DHCP packet from unconfigured interface", inIfIdx)
		return
	}

	txOk := false
	for _, destIpAddr := range draIntf.ServerIp {
		destIpPortString := fmt.Sprintf(
			"[%s]:%s", destIpAddr, strconv.Itoa(DHCP_SERVER_PORT))
		intfServerState := pProc.initIntfServerState(inIfName, destIpAddr)
		pProc.StateMutex.Lock() // State obj lock
		if pProc.SendPkt(outPkt, destIpPortString) {
			txOk = true
			intfState.TotalDhcpServerTx++
			if mType != RELAY_FORW {
				clientState.ServerRequests++
			}
			intfServerState.Request++
		}
		pProc.StateMutex.Unlock() // State obj unlock
	}
	for _, ifRef := range draIntf.UpstreamIntfs {
		_, ok := pProc.InfraMgr.GetIPv6LLIntfIndex(ifRef)
		if !ok {
			Logger.Debug("DRA: No IPv6 link-Local interface found for", ifRef)
			continue
		}
		//Link local multicast addr
		destIpAddr := fmt.Sprintf("%s%%%s", ALL_DRA_SERVERS_ADDR, ifRef)
		destIpPortString := fmt.Sprintf(
			"[%s]:%s", destIpAddr, strconv.Itoa(DHCP_SERVER_PORT))
		intfServerState := pProc.initIntfServerState(inIfName, destIpAddr)
		pProc.StateMutex.Lock() // State obj lock
		if pProc.SendPkt(outPkt, destIpPortString) {
			txOk = true
			intfState.TotalDhcpServerTx++
			if mType != RELAY_FORW {
				clientState.ServerRequests++
			}
			intfServerState.Request++
		}
		pProc.StateMutex.Unlock() // State obj unlock
	}
	if !txOk {
		pProc.StateMutex.Lock() // State obj lock
		intfState.TotalDrops++
		pProc.StateMutex.Unlock() // State obj unlock
	}
}

func (pProc *Processor) validateRxUpstream(
	inIfIdx int, inIfName string, srcAddr net.IP) (net.IP, bool) {

	if _, ok := pProc.InfraMgr.GetActiveDRAv6Intf(inIfIdx); !ok {
		Logger.Debug(
			"Dropping DHCP packet from unconfigured interface", inIfName)
		return nil, false
	}
	ipv6Intf, ok := pProc.InfraMgr.GetIPv6Intf(inIfIdx)
	if !ok {
		infoMsg := fmt.Sprintln("Cannot get non linklocal IPv6Intf for interface,", inIfIdx)
		Logger.Debug(infoMsg)
		return nil, false
	}
	intfAddr := net.ParseIP(ipv6Intf.IpAddr)
	return intfAddr, true
}

func (pProc *Processor) MakeSendTxPkt(
	inIfIdx int, inIfName string, inPktBuf []byte, srcAddr net.IP) bool {

	outPkt := DhcpRelayPacket(make([]byte, 34, 1500))
	mType := MsgType(inPktBuf[0])
	switch mType {
	case SOLICIT, REQUEST, CONFIRM, RENEW, REBIND, DECLINE, RELEASE,
		INFO_REQ: // Packet Upstream

		inIfAddr, ok := pProc.validateRxUpstream(inIfIdx, inIfName, srcAddr)
		if !ok {
			return false
		}
		intfState := pProc.initIntfState(inIfName)
		clientMacAddr := DhcpPacket(inPktBuf).GetClientHwAddr()
		if clientMacAddr == nil {
			Logger.Err("Cannot get client mac address ")
			return false
		}
		clientState := pProc.initClientState(clientMacAddr.String())
		dhcpOptions := DhcpPacket(inPktBuf).ParseOptions([]OptionType{})
		pProc.setUpstreamInState(mType, dhcpOptions, intfState, clientState)
		pProc.PeerAddrIntfMap[srcAddr.String()] = inIfName

		outPkt.SetMsgType(int8(RELAY_FORW))
		outPkt.SetHopCount(1)
		outPkt.SetLinkAddrField(inIfAddr)
		outPkt.SetPeerAddrField(srcAddr)
		outPkt := outPkt.AddRelayMsgOption(inPktBuf)
		pProc.SendPktUpstream(
			outPkt, mType, inIfIdx, inIfName,
			srcAddr, intfState, clientState)

	case RELAY_FORW: // Packet Upstream
		inIfAddr, ok := pProc.validateRxUpstream(inIfIdx, inIfName, srcAddr)
		if !ok {
			return false
		}
		intfState := pProc.initIntfState(inIfName)
		pProc.setUpstreamInState(mType, nil, intfState, nil)
		pProc.PeerAddrIntfMap[srcAddr.String()] = inIfName

		tmpPkt := DhcpRelayPacket(inPktBuf)
		if tmpPkt.GetHopCount() >= HOP_COUNT_LIMIT {
			return false
		}
		outPkt.SetMsgType(int8(RELAY_FORW))
		outPkt.SetHopCount(tmpPkt.GetHopCount() + 1)
		if srcAddr.IsLinkLocalUnicast() {
			outPkt.SetLinkAddrField(inIfAddr)
		}
		outPkt.SetPeerAddrField(srcAddr)
		outPkt := outPkt.AddRelayMsgOption(inPktBuf)
		pProc.SendPktUpstream(
			outPkt, mType, inIfIdx, inIfName,
			srcAddr, intfState, nil)

	case RELAY_REPL: // Packet Downstream
		inPkt := DhcpRelayPacket(inPktBuf)
		peerAddr := net.IP(inPkt.GetPeerAddrField())
		newOutPkt := inPkt.GetDRAOptionField().GetPayload()
		outMsgType := MsgType(newOutPkt[0])
		if outMsgType != RELAY_REPL {
			clientMacAddr := DhcpPacket(inPktBuf).GetClientHwAddr()
			if clientMacAddr == nil {
				Logger.Err("Cannot get client mac address ")
				return false
			}
			clientState := pProc.initClientState(clientMacAddr.String())
			dhcpOptions := DhcpPacket(newOutPkt).ParseOptions([]OptionType{})
			pProc.SendPktDownstream(
				newOutPkt, dhcpOptions, peerAddr, srcAddr, clientState)
		} else {
			pProc.SendPktDownstream(
				newOutPkt, nil, peerAddr, srcAddr, nil)
		}
	}
	return true
}

func (pProc *Processor) RxTx() {
	var buf DhcpPacket = make([]byte, 1500)
	for {
		Logger.Debug("DRA: Calling ReadFrom")
		bytesRead, cm, srcAddr, err := pProc.ClientConn.ReadFrom(buf)
		if err != nil {
			Logger.Err("DRA: reading buffer failed")
			break
		}
		Logger.Debug("DRA: Received Packet from ", srcAddr)
		inIf, err := net.InterfaceByIndex(cm.IfIndex)
		if err != nil {
			Logger.Err("DRA: Linux interface not found for index", cm.IfIndex)
		}
		ifIdx, ok := pProc.InfraMgr.GetIPv6IntfIndex(inIf.Name)
		ifIdxLL, okLL := pProc.InfraMgr.GetIPv6LLIntfIndex(inIf.Name)
		var inIfIdx int
		if ok {
			inIfIdx = ifIdx
		} else if okLL {
			inIfIdx = ifIdxLL
		} else {
			Logger.Err("DRA: IPv6 Intf not found for Intf", inIf.Name)
			continue
		}
		srcUdpAddr, _ := net.ResolveUDPAddr("udp6", srcAddr.String())
		pProc.MakeSendTxPkt(inIfIdx, inIf.Name, buf[:bytesRead], srcUdpAddr.IP)
	}
	pProc.EnabledMutex.Lock()
	pProc.EnabledFlag = false
	pProc.EnabledMutex.Unlock()
}

func (pProc *Processor) CreateConn() bool {
	saddr := net.UDPAddr{
		IP:   net.ParseIP(""),
		Port: DHCP_SERVER_PORT,
	}
	var err error
	pProc.ClientHandler, err = net.ListenUDP("udp6", &saddr)
	if err != nil {
		Logger.Err("Opening udp port for client --> server failed", err)
		return false
	}
	pProc.ClientConn = ipv6.NewPacketConn(pProc.ClientHandler)
	controlFlag := ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface
	err = pProc.ClientConn.SetControlMessage(controlFlag, true)
	if err != nil {
		Logger.Err("Setting control flag for client failed..", err)
		return false
	}
	Logger.Debug("Client Connection opened successfully")
	return true
}

func (pProc *Processor) RegisterMcast(ifRef string) {
	ifObj, err := net.InterfaceByName(ifRef)
	if err != nil {
		Logger.Debug("Cannot find interface:", ifRef)
		return
	}
	Logger.Debug("Registering Mcast for intf", ifRef)
	err = pProc.ClientConn.JoinGroup(
		ifObj, &net.UDPAddr{IP: net.ParseIP(ALL_DRA_SERVERS_ADDR)},
	)
	if err != nil {
		Logger.Debug(
			"Cannot join DRA multicast group for interface:",
			ifRef, "reason:", err,
		)
		return
	}
}

func (pProc *Processor) DeregisterMcast(ifRef string) {
	ifObj, err := net.InterfaceByName(ifRef)
	if err != nil {
		Logger.Debug("Cannot find interface:", ifRef)
		return
	}
	Logger.Debug("Deregistering Mcast for intf", ifRef)

	err = pProc.ClientConn.LeaveGroup(
		ifObj, &net.UDPAddr{IP: net.ParseIP(ALL_DRA_SERVERS_ADDR)},
	)
	if err != nil {
		Logger.Debug(
			"Cannot leave DRA multicast group for interface:",
			ifRef, "reason:", err,
		)
		return
	}
}

func (pProc *Processor) ProcessCreateDRAIntf(ifName string) {
	pProc.initIntfState(ifName)
}

func (pProc *Processor) ProcessDeleteDRAIntf(ifName string) {
	pProc.deleteIntfState(ifName)
}

func (pProc *Processor) ProcessActiveDRAIntf(ifIdx int) {
	intf, ok := pProc.InfraMgr.GetIPv6Intf(ifIdx)
	if !ok {
		Logger.Debug("No Intf found", intf.IfRef)
		return
	}
	pProc.RegisterMcast(intf.IfRef)
}

func (pProc *Processor) ProcessInactiveDRAIntf(ifIdx int) {
	intf, ok := pProc.InfraMgr.GetIPv6Intf(ifIdx)
	if !ok {
		Logger.Debug("No Intf found", intf.IfRef)
		return
	}
	pProc.DeregisterMcast(intf.IfRef)
}
