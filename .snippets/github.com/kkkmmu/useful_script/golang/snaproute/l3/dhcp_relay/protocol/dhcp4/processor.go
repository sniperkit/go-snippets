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

package dhcp4

import (
	"dhcprelayd"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/ipv4"
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
	ClientConn    *ipv4.PacketConn

	PeerAddrIntfMap map[string]string
	PcapHandles     map[int]*pcap.Handle

	// States
	ClientStateSlice     []*dhcprelayd.DHCPRelayClientState
	ClientStateMap       map[string]*dhcprelayd.DHCPRelayClientState
	IntfStateSlice       []*dhcprelayd.DHCPRelayIntfState
	IntfStateMap         map[string]*dhcprelayd.DHCPRelayIntfState
	IntfServerStateSlice []*dhcprelayd.DHCPRelayIntfServerState
	IntfServerStateMap   map[string]*dhcprelayd.DHCPRelayIntfServerState
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
	pProc.PcapHandles = make(map[int]*pcap.Handle)
	pProc.ClientStateSlice = []*dhcprelayd.DHCPRelayClientState{}
	pProc.ClientStateMap = make(map[string]*dhcprelayd.DHCPRelayClientState)
	pProc.IntfStateSlice = []*dhcprelayd.DHCPRelayIntfState{}
	pProc.IntfStateMap = make(map[string]*dhcprelayd.DHCPRelayIntfState)
	pProc.IntfServerStateSlice = []*dhcprelayd.DHCPRelayIntfServerState{}
	pProc.IntfServerStateMap = make(map[string]*dhcprelayd.DHCPRelayIntfServerState)

	pProc.EnabledMutex.Lock()
	pProc.EnabledFlag = false
	pProc.EnabledMutex.Unlock()

	return pProc
}

func (pProc *Processor) GetEnabledFlag() bool {
	return false
}
func (pProc *Processor) SetEnabledFlag() {}

func (pProc *Processor) GetClientStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPRelayClientState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.ClientStateSlice)
	if fromIdx < 0 || fromIdx >= length || count <= 0 {
		return 0, 0, false, []*dhcprelayd.DHCPRelayClientState{}
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

	result := make([]*dhcprelayd.DHCPRelayClientState, actualCount)
	copy(result, pProc.ClientStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor) GetClientState(
	macAddr string) (*dhcprelayd.DHCPRelayClientState, bool) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	if clientState, ok := pProc.ClientStateMap[macAddr]; ok {
		return clientState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor) GetIntfStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPRelayIntfState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.IntfStateSlice)
	if fromIdx < 0 || fromIdx >= length || count <= 0 {
		return 0, 0, false, []*dhcprelayd.DHCPRelayIntfState{}
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

	result := make([]*dhcprelayd.DHCPRelayIntfState, actualCount)
	copy(result, pProc.IntfStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor) GetIntfState(
	ifName string) (*dhcprelayd.DHCPRelayIntfState, bool) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	if intfState, ok := pProc.IntfStateMap[ifName]; ok {
		return intfState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor) GetIntfServerStateSlice(
	fromIdx, count int) (int, int, bool, []*dhcprelayd.DHCPRelayIntfServerState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	var nextIdx int
	var more bool
	var actualCount int
	length := len(pProc.IntfServerStateSlice)
	if fromIdx < 0 || fromIdx >= length || count <= 0 {
		return 0, 0, false, []*dhcprelayd.DHCPRelayIntfServerState{}
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

	result := make([]*dhcprelayd.DHCPRelayIntfServerState, actualCount)
	copy(result, pProc.IntfServerStateSlice[fromIdx:fromIdx+actualCount])
	return nextIdx, actualCount, more, result
}

func (pProc *Processor) GetIntfServerState(
	ifName string, serverAddr string) (*dhcprelayd.DHCPRelayIntfServerState, bool) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfServerStateKey := ifName + "_" + serverAddr
	if intfServerState, ok := pProc.IntfServerStateMap[intfServerStateKey]; ok {
		return intfServerState, true
	} else {
		return nil, false
	}
}

func (pProc *Processor) initClientState(
	clientAddr string) *dhcprelayd.DHCPRelayClientState {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	clientStateKey := clientAddr
	clientState, ok := pProc.ClientStateMap[clientStateKey]
	if !ok {
		clientState = &dhcprelayd.DHCPRelayClientState{}
		clientState.MacAddr = clientAddr
		pProc.ClientStateSlice = append(pProc.ClientStateSlice, clientState)
		pProc.ClientStateMap[clientStateKey] = clientState
	}
	return clientState
}

func (pProc *Processor) initIntfState(
	ifName string) *dhcprelayd.DHCPRelayIntfState {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfState, ok := pProc.IntfStateMap[ifName]
	if !ok {
		intfState = &dhcprelayd.DHCPRelayIntfState{}
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
	ifName string, serverAddr string) *dhcprelayd.DHCPRelayIntfServerState {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfServerStateKey := ifName + "_" + serverAddr
	intfServerState, ok := pProc.IntfServerStateMap[intfServerStateKey]
	if !ok {
		intfServerState = &dhcprelayd.DHCPRelayIntfServerState{}
		intfServerState.IntfRef = ifName
		intfServerState.ServerIp = serverAddr
		pProc.IntfServerStateSlice = append(
			pProc.IntfServerStateSlice, intfServerState)
		pProc.IntfServerStateMap[intfServerStateKey] = intfServerState
	}
	return intfServerState
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

func (pProc *Processor) CreateConn() bool {
	saddr := net.UDPAddr{
		IP:   net.ParseIP(""),
		Port: DHCP_SERVER_PORT,
	}
	var err error
	pProc.ClientHandler, err = net.ListenUDP("udp4", &saddr)
	if err != nil {
		Logger.Err("Opening udp port for client --> server failed", err)
		return false
	}
	pProc.ClientConn = ipv4.NewPacketConn(pProc.ClientHandler)
	controlFlag := ipv4.FlagSrc | ipv4.FlagDst | ipv4.FlagInterface
	err = pProc.ClientConn.SetControlMessage(controlFlag, true)
	if err != nil {
		Logger.Err("Setting control flag for client failed..", err)
		return false
	}
	Logger.Debug("Client Connection opened successfully")
	return true
}

func (pProc *Processor) setUpstreamInState(
	mType MessageType,
	inPkt DhcpRelayAgentPacket,
	requestedIp string,
	intfState *dhcprelayd.DHCPRelayIntfState,
	clientState *dhcprelayd.DHCPRelayClientState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfState.TotalDhcpClientRx++
	clientState.ClientRequests++

	switch mType {
	case DhcpDiscover:
		clientState.ClientDiscover = time.Now().String()
	case DhcpRequest:
		clientState.ClientRequest = time.Now().String()
		clientState.RequestedIp = requestedIp
	}
}

func (pProc *Processor) setDownstreamInState(
	mType MessageType,
	inPkt DhcpRelayAgentPacket,
	clientState *dhcprelayd.DHCPRelayClientState,
	intfState *dhcprelayd.DHCPRelayIntfState) {

	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	intfState.TotalDhcpServerRx++
	clientState.ServerResponses++

	switch mType {
	case DhcpOffer:
		clientState.OfferedIp = inPkt.GetYIAddr().String()
		clientState.ServerOffer = time.Now().String()
	case DhcpACK:
		clientState.AcceptedIp = inPkt.GetYIAddr().String()
		clientState.ServerAck = time.Now().String()
	}
}

func (pProc *Processor) DhcpRelayAgentSendClientOptPacket(
	inIntfProp *infra.IPv4IntfProperty,
	inReq DhcpRelayAgentPacket, reqOptions DhcpRelayAgentOptions,
	mt MessageType, intfState *dhcprelayd.DHCPRelayIntfState,
	clientState *dhcprelayd.DHCPRelayClientState) {

	// Create Packet
	outPacket := DhcpRelayAgentCreateNewPacket(Request, inReq)
	if inReq.GetGIAddr().String() == DHCP_NO_IP {
		outPacket.SetGIAddr(net.ParseIP(inIntfProp.IpAddr))
	} else {
		Logger.Debug("DRA: Relay Agent " + inReq.GetGIAddr().String() +
			" requested for DHCP for HOST " + inReq.GetCHAddr().String())
		outPacket.SetGIAddr(inReq.GetGIAddr())
	}

	requestedIp, serverIp := DhcpRelayAgentAddOptionsToPacket(reqOptions,
		mt, &outPacket)
	if serverIp == "" {
		Logger.Warning("DRA: no server ip.. dropping the request")
		pProc.StateMutex.Lock() // State obj lock
		intfState.TotalDrops++
		pProc.StateMutex.Unlock() // State obj unlock
		return
	}
	// get host + server state entry for updating the state
	// Create server ip address + port number
	serverIpPort := serverIp + ":" + strconv.Itoa(DHCP_SERVER_PORT)
	Logger.Debug("DRA: Sending " + ParseMessageTypeToString(mt) +
		" packet to " + serverIpPort)
	serverAddr, err := net.ResolveUDPAddr("udp", serverIpPort)
	if err != nil {
		Logger.Err("DRA: couldn't resolved udp addr for and err is", err)
		pProc.StateMutex.Lock() // State obj lock
		intfState.TotalDrops++
		pProc.StateMutex.Unlock() // State obj unlock
		return
	}
	// Pad to minimum size of dhcp packet
	outPacket.PadToMinSize()
	// send out the packet...
	pProc.setUpstreamInState(mt, inReq, requestedIp, intfState, clientState)
	intfServerState := pProc.initIntfServerState(inIntfProp.IfRef, serverIp)
	_, err = pProc.ClientHandler.WriteToUDP(outPacket, serverAddr)
	if err != nil {
		Logger.Debug("DRA: WriteToUDP failed with error:", err)
		pProc.StateMutex.Lock() // State obj lock
		intfState.TotalDrops++
		pProc.StateMutex.Unlock() // State obj unlock
		return
	}
	pProc.StateMutex.Lock() // State obj lock
	intfState.TotalDhcpServerTx++
	clientState.ServerRequests++
	intfServerState.Request++
	pProc.StateMutex.Unlock() // State obj unlock
	Logger.Debug("DRA: Create & Send of PKT successfully to server", serverIp)
}

func DhcpRelayAgentAddOptionsToPacket(reqOptions DhcpRelayAgentOptions, mt MessageType,
	outPacket *DhcpRelayAgentPacket) (string, string) {
	outPacket.AddDhcpOptions(OptionDHCPMessageType, []byte{byte(mt)})
	var dummyDup map[DhcpOptionCode]int
	var reqIp string
	var serverIp = ""
	dummyDup = make(map[DhcpOptionCode]int, len(reqOptions))
	for i := 0; i < len(reqOptions); i++ {
		opt := reqOptions.SelectOrderOrAll(reqOptions[DhcpOptionCode(i)])
		for _, option := range opt {
			_, ok := dummyDup[option.Code]
			if ok {
				continue
			}
			switch option.Code {
			case OptionRequestedIPAddress:
				reqIp = net.IPv4(option.Value[0], option.Value[1],
					option.Value[2], option.Value[3]).String()
				break
			case OptionServerIdentifier:
				serverIp = net.IPv4(option.Value[0], option.Value[1],
					option.Value[2], option.Value[3]).String()
				break
			}
			outPacket.AddDhcpOptions(option.Code, option.Value)
			dummyDup[option.Code] = 9999
		}
	}
	return reqIp, serverIp
}

func (pProc *Processor) DhcpRelayAgentSendPacketToDhcpClient(
	inReq DhcpRelayAgentPacket, outIfName string,
	reqOptions DhcpRelayAgentOptions, mt MessageType,
	serverIp net.IP) {

	var outPacket DhcpRelayAgentPacket
	outPacket = DhcpRelayAgentCreateNewPacket(Reply, inReq)
	DhcpRelayAgentAddOptionsToPacket(reqOptions, mt, &outPacket)
	// Pad to minimum size of dhcp packet
	outPacket.PadToMinSize()

	netIntf, err := net.InterfaceByName(outIfName)
	if err != nil {
		Logger.Debug("Could not find interface by name", outIfName)
		return
	}
	ifIdx, ok := pProc.InfraMgr.GetIPv4IntfIndex(outIfName)
	if !ok {
		Logger.Debug("Cannot get non IPv4 IfIndex,", ifIdx)
		return
	}
	ipv4Intf, ok := pProc.InfraMgr.GetIPv4Intf(ifIdx)
	if !ok {
		Logger.Debug("Cannot get non IPv4Intf for IfIndex,", ifIdx)
		return
	}
	eth := &layers.Ethernet{
		SrcMAC:       netIntf.HardwareAddr,
		DstMAC:       outPacket.GetCHAddr(),
		EthernetType: layers.EthernetTypeIPv4,
	}
	ipv4 := &layers.IPv4{
		SrcIP:    net.ParseIP(ipv4Intf.IpAddr),
		DstIP:    outPacket.GetYIAddr(),
		Version:  4,
		Protocol: layers.IPProtocolUDP,
		TTL:      64,
	}
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(DHCP_SERVER_PORT),
		DstPort: layers.UDPPort(DHCP_CLIENT_PORT),
	}
	udp.SetNetworkLayerForChecksum(ipv4)

	goOpts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	buffer := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, goOpts, eth, ipv4, udp,
		gopacket.Payload(outPacket))

	intfState := pProc.initIntfState(outIfName)
	clientState := pProc.initClientState(outPacket.GetCHAddr().String())
	intfServerState := pProc.initIntfServerState(outIfName, serverIp.String())
	pProc.StateMutex.Lock()
	pcapHdl, ok := pProc.PcapHandles[ifIdx]
	pProc.StateMutex.Unlock()
	if !ok {
		Logger.Debug("DRA: opening pcap handle for", outIfName)
		pcapHdl, err = pcap.OpenLive(outIfName, snapshot_len,
			promiscuous, timeout)
		if err != nil {
			Logger.Err("DRA: opening pcap for", outIfName, "failed with Error:", err)
			pProc.StateMutex.Lock()
			intfState.TotalDrops++
			pProc.StateMutex.Unlock()
			return
		}
		pProc.StateMutex.Lock()
		pProc.PcapHandles[ifIdx] = pcapHdl
		pProc.StateMutex.Unlock()
	}
	err = pcapHdl.WritePacketData(buffer.Bytes())
	if err != nil {
		Logger.Debug("DRA: WritePacketData failed with error:", err)
		pProc.StateMutex.Lock()
		intfState.TotalDrops++
		pProc.StateMutex.Unlock()
		return
	}
	pProc.StateMutex.Lock()
	intfState.TotalDhcpClientTx++
	intfServerState.Responses++
	clientState.ClientResponses++
	pProc.StateMutex.Unlock()
	Logger.Debug("DRA: Create & Send of PKT successfully to client")
}

func (pProc *Processor) DhcpRelayAgentSendDiscoverPacket(
	inIntfProp *infra.IPv4IntfProperty,
	inReq DhcpRelayAgentPacket, reqOptions DhcpRelayAgentOptions,
	mt MessageType, intfState *dhcprelayd.DHCPRelayIntfState,
	clientState *dhcprelayd.DHCPRelayClientState) {

	Logger.Debug("DRA: Sending Discover Request")
	draIntf, ok := pProc.InfraMgr.GetActiveDRAv4Intf(inIntfProp.IfIndex)
	if !ok {
		Logger.Debug(
			"Dropping DHCP packet from unconfigured interface", inIntfProp.IfIndex)
		pProc.StateMutex.Lock() // State obj lock
		intfState.TotalDrops++
		pProc.StateMutex.Unlock() // State obj unlock
		return
	}
	pProc.setUpstreamInState(mt, inReq, "", intfState, clientState)

	txOk := false
	for i := 0; i < len(draIntf.ServerIp); i++ {
		serverIpPort := draIntf.ServerIp[i] + ":" +
			strconv.Itoa(DHCP_SERVER_PORT)
		Logger.Debug("DRA: Sending DHCP PACKET to server: " + serverIpPort)
		serverAddr, err := net.ResolveUDPAddr("udp", serverIpPort)
		if err != nil {
			Logger.Err("DRA: couldn't resolved udp addr for and err is", err)
			continue
		}

		outPacket := DhcpRelayAgentCreateNewPacket(Request, inReq)
		if inReq.GetGIAddr().String() == DHCP_NO_IP {
			outPacket.SetGIAddr(net.ParseIP(inIntfProp.IpAddr))
		} else {
			Logger.Debug("DRA: Relay Agent " + inReq.GetGIAddr().String() +
				" requested for DHCP for HOST " + inReq.GetCHAddr().String())
			outPacket.SetGIAddr(inReq.GetGIAddr())
		}

		DhcpRelayAgentAddOptionsToPacket(reqOptions, mt, &outPacket)
		// Pad to minimum size of dhcp packet
		outPacket.PadToMinSize()
		// send out the packet...
		intfServerState := pProc.initIntfServerState(inIntfProp.IfRef, draIntf.ServerIp[i])
		_, err = pProc.ClientHandler.WriteToUDP(outPacket, serverAddr)
		if err != nil {
			Logger.Debug("DRA: WriteToUDP failed with error:", err)
			continue
		}
		txOk = true
		pProc.StateMutex.Lock()
		intfState.TotalDhcpServerTx++
		clientState.ServerRequests++
		clientState.ServerIp = draIntf.ServerIp[i]
		intfServerState.Request++
		pProc.StateMutex.Unlock()
		Logger.Debug("DRA: Create & Send of PKT successfully to server", draIntf.ServerIp[i])
	}
	if !txOk {
		pProc.StateMutex.Lock() // State obj lock
		intfState.TotalDrops++
		pProc.StateMutex.Unlock() // State obj unlock
		return
	}
}

func (pProc *Processor) DhcpRelayAgentSendPacketToDhcpServer(
	inIntfProp *infra.IPv4IntfProperty, inReq DhcpRelayAgentPacket,
	reqOptions DhcpRelayAgentOptions, mt MessageType,
	intfState *dhcprelayd.DHCPRelayIntfState,
	clientState *dhcprelayd.DHCPRelayClientState) {

	switch mt {
	case DhcpDiscover:
		pProc.DhcpRelayAgentSendDiscoverPacket(
			inIntfProp, inReq, reqOptions, mt, intfState, clientState)
		break
	case DhcpRequest, DhcpDecline, DhcpRelease, DhcpInform:
		pProc.DhcpRelayAgentSendClientOptPacket(
			inIntfProp, inReq, reqOptions, mt, intfState, clientState)
		break
	}
}

func (pProc *Processor) validateRxUpstream(
	inIfIdx int, inIfName string, srcAddr net.IP) (*infra.IPv4IntfProperty, bool) {

	if _, ok := pProc.InfraMgr.GetActiveDRAv4Intf(inIfIdx); !ok {
		Logger.Debug(
			"Dropping DHCP packet from unconfigured interface", inIfName)
		return nil, false
	}
	ipv4Intf, ok := pProc.InfraMgr.GetIPv4Intf(inIfIdx)
	if !ok {
		Logger.Debug("Cannot get non IPv4Intf for index,", inIfIdx)
		return nil, false
	}
	return ipv4Intf, true
}

func (pProc *Processor) DhcpRelayAgentSendPacket(inIfIdx int, inIfName string,
	inReq DhcpRelayAgentPacket, srcAddr net.IP, reqOptions DhcpRelayAgentOptions,
	mType MessageType) {

	switch mType {
	case DhcpDiscover, DhcpRequest, DhcpDecline, DhcpRelease, DhcpInform:
		// Use obtained logical id to find the global interface object
		inIfProp, ok := pProc.validateRxUpstream(inIfIdx, inIfName, srcAddr)
		if !ok {
			return
		}
		clientMacAddr := inReq.GetCHAddr().String()
		intfState := pProc.initIntfState(inIfName)
		clientState := pProc.initClientState(clientMacAddr)
		pProc.PeerAddrIntfMap[clientMacAddr] = inIfName
		// Send Packet
		pProc.DhcpRelayAgentSendPacketToDhcpServer(
			inIfProp, inReq, reqOptions, mType, intfState, clientState)
	case DhcpOffer, DhcpACK, DhcpNAK:
		// Get the interface from reverse mapping to send the unicast
		// packet...
		outIfName, ok := pProc.PeerAddrIntfMap[inReq.GetCHAddr().String()]
		if !ok {
			Logger.Err("DRA: cache for linux interface for " +
				inReq.GetCHAddr().String() + " not present")
			return
		}
		clientMacAddr := inReq.GetCHAddr().String()
		intfState := pProc.initIntfState(inIfName)
		clientState := pProc.initClientState(clientMacAddr)
		pProc.setDownstreamInState(mType, inReq, clientState, intfState)
		pProc.DhcpRelayAgentSendPacketToDhcpClient(inReq,
			outIfName, reqOptions, mType, srcAddr)
	default:
		Logger.Debug("DRA: any new message type")
	}
	return
}

func (pProc *Processor) MakeSendTxPkt(
	inIfIdx int, inIfName string, inPktBuf []byte, srcAddr net.IP) bool {

	// var buf []byte = make([]byte, 1500)
	bytesRead := len(inPktBuf)
	if bytesRead < DHCP_PACKET_MIN_BYTES {
		// This is not dhcp packet as the minimum size is 240
		//intfState.TotalDrops++
		//dhcprelayIntfStateMap[intfId] = intfState
		return false
	}
	inReq, reqOptions, mType := DhcpRelayAgentDecodeInPkt(inPktBuf, bytesRead)
	pProc.DhcpRelayAgentSendPacket(inIfIdx, inIfName, inReq, srcAddr, reqOptions,
		mType)
	return true
}

func (pProc *Processor) RxTx() {
	buf := make([]byte, 1500)
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
		ifIdx, ok := pProc.InfraMgr.GetIPv4IntfIndex(inIf.Name)
		if !ok {
			Logger.Err("DRA: IPv4 Intf not found for Intf", inIf.Name)
			continue
		}
		srcUdpAddr, _ := net.ResolveUDPAddr("udp4", srcAddr.String())
		pProc.MakeSendTxPkt(ifIdx, inIf.Name, buf[:bytesRead], srcUdpAddr.IP)
	}
	pProc.EnabledMutex.Lock()
	pProc.EnabledFlag = false
	pProc.EnabledMutex.Unlock()
}

func (pProc *Processor) openPcapHandler(ifIdx int) {
	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	ipv4Intf, ok := pProc.InfraMgr.GetIPv4Intf(ifIdx)
	if !ok {
		Logger.Debug("DRA: Unable to find interface with index", ifIdx)
	}
	pcapHdl, ok := pProc.PcapHandles[ifIdx]
	if ok {
		pcapHdl.Close()
	}
	Logger.Debug("DRA: opening pcap handle for", ipv4Intf.IfRef)
	pcapHdl, err := pcap.OpenLive(ipv4Intf.IfRef, snapshot_len,
		promiscuous, timeout)
	if err != nil {
		Logger.Err("DRA: opening pcap for", ifIdx, "failed with Error:", err)
		return
	}
	pProc.PcapHandles[ifIdx] = pcapHdl
}

func (pProc *Processor) closePcapHandler(ifIdx int) {
	defer pProc.StateMutex.Unlock()
	pProc.StateMutex.Lock()

	pcapHdl, ok := pProc.PcapHandles[ifIdx]
	if ok {
		pcapHdl.Close()
		delete(pProc.PcapHandles, ifIdx)
	}
}

func (pProc *Processor) ProcessCreateDRAIntf(ifName string) {
	pProc.initIntfState(ifName)
}

func (pProc *Processor) ProcessDeleteDRAIntf(ifName string) {
	pProc.deleteIntfState(ifName)
}

func (pProc *Processor) ProcessActiveDRAIntf(ifIdx int) {
	pProc.openPcapHandler(ifIdx)
}

func (pProc *Processor) ProcessInactiveDRAIntf(ifIdx int) {
	pProc.closePcapHandler(ifIdx)
}
