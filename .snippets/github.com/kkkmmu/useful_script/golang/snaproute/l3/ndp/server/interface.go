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
package server

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"l3/ndp/config"
	"l3/ndp/debug"
	"l3/ndp/packet"
	"net"
	"strings"
	"time"
)

const (
	NDP_PCAP_FILTER                              = "(ip6[6] == 0x3a) and (ip6[40] >= 133 && ip6[40] <= 136)"
	NDP_ETHER_SRC                                = " and not ether src "
	NDP_PCAP_TIMEOUT                             = 1 * time.Second
	NDP_PCAP_SNAPSHOTlEN                         = 1024
	NDP_PCAP_PROMISCUOUS                         = false
	MIN_DELAY_BETWEEN_RAS                  uint8 = 3 // RFC: 4861
	MAX_INITIAL_RTR_ADVERTISEMENTS         uint8 = 3
	MAX_INITIAL_RTR_ADVERT_INTERVAL        uint8 = 16
	ALL_NODES_MULTICAST_IPV6_ADDRESS             = "ff02::1"
	ALL_NODES_MULTICAST_LINK_LAYER_ADDRESS       = "33:33:00:00:00:01"

	NDP_DEFAULT_RTR_ADVERTISEMENT_INTERVAL uint8  = 5
	NDP_DEFAULT_RETRANSMIT_INTERVAL        uint32 = 1
	NDP_DEFAULT_REACHABLE_INTERVAL         uint32 = 30000
)

/* https://tools.ietf.org/html/rfc7346
   +------+--------------------------+-------------------------+
   | scop | NAME                     | REFERENCE               |
   +------+--------------------------+-------------------------+
   |  0   | Reserved                 | [RFC4291], RFC 7346     |
   |  1   | Interface-Local scope    | [RFC4291], RFC 7346     |
   |  2   | Link-Local scope         | [RFC4291], RFC 7346     |
   |  3   | Realm-Local scope        | [RFC4291], RFC 7346     |
   |  4   | Admin-Local scope        | [RFC4291], RFC 7346     |
   |  5   | Site-Local scope         | [RFC4291], RFC 7346     |
   |  6   | Unassigned               |                         |
   |  7   | Unassigned               |                         |
   |  8   | Organization-Local scope | [RFC4291], RFC 7346     |
   |  9   | Unassigned               |                         |
   |  A   | Unassigned               |                         |
   |  B   | Unassigned               |                         |
   |  C   | Unassigned               |                         |
   |  D   | Unassigned               |                         |
   |  E   | Global scope             | [RFC4291], RFC 7346     |
   |  F   | Reserved                 | [RFC4291], RFC 7346     |
   +------+--------------------------+-------------------------+
*/
var IPV6_MULTICAST_PREFIXES = []string{"ff00", "ff01", "ff02", "ff03", "ff04", "ff05", "ff06", "ff07",
	"ff08", "ff09", "ff0a", "ff0b", "ff0c", "ff0d", "ff0e", "ff0f"}

type PcapBase struct {
	// TX Pcap handler
	Tx *pcap.Handle
	// RX Pcap Handler for Each Port
	PcapHandle *pcap.Handle
	// at any give time there can be two users for Pcap..
	// if 0 then only start rx/tx
	// if 1 then only stop rx/tx
	PcapUsers uint8
}

type PktCounter struct {
	Send int64
	Rcvd int64
}

type Interface struct {
	PcapBase
	IntfRef           string
	IfIndex           int32
	IfType            int    // IfTypePort, IfTypeVlan
	IpAddr            string // CIDR Format
	LinkLocalIp       string // CIDR format
	globalScope       string // absolute
	linkScope         string // absolute
	MsgType           string
	OperState         string
	reachableTime     uint32
	retransTime       uint32
	routerLifeTime    uint16
	raRestransmitTime uint8 // @TODO: get it from user
	raTimer           *time.Timer
	initialRASend     uint8                   // on port up we have to send 3 RA before kicking in config timer
	Neighbor          map[string]NeighborInfo // key is NbrIp_NbrMac to handle move scenario's
	incompleteNbr     map[string]struct{}
	PktDataCh         chan config.PacketData
	counter           PktCounter
}

func (intf *Interface) addIP(ipAddr string) {
	if isLinkLocal(ipAddr) {
		intf.LinkLocalIp = ipAddr
		ip, _, err := net.ParseCIDR(intf.LinkLocalIp)
		if err != nil {
			debug.Logger.Err("Parsing link local ip failed", err)
		} else {
			intf.linkScope = ip.String()
		}
	} else {
		intf.IpAddr = ipAddr
		ip, _, err := net.ParseCIDR(intf.IpAddr)
		if err != nil {
			debug.Logger.Err("Parsing Global Scope ip failed", err)
		} else {
			intf.globalScope = ip.String()
		}
	}
}

func (intf *Interface) removeIP(ipAddr string) {
	if isLinkLocal(ipAddr) {
		intf.LinkLocalIp = ""
		intf.linkScope = ""
	} else {
		intf.IpAddr = ""
		intf.globalScope = ""
	}
}

/*
 * common init params between InitIntf and CreateIntf
 */
func (intf *Interface) commonInit(ipAddr string, pktCh chan config.PacketData, gCfg NdpConfig) {
	intf.addIP(ipAddr)
	// Pcap Init
	intf.PcapBase.PcapHandle = nil
	intf.PcapBase.PcapUsers = 0
	// Timers Value Init
	intf.retransTime = gCfg.RetransTime             //1       // config value ms
	intf.reachableTime = gCfg.ReachableTime         //30000 // config value ms
	intf.raRestransmitTime = gCfg.RaRestransmitTime //5 // config value s ADAM asked for 5 seconds :)
	intf.routerLifeTime = 1800                      // config value s
	intf.initialRASend = 0
	intf.raTimer = nil
	// Neighbor Init
	intf.PktDataCh = pktCh
	intf.Neighbor = make(map[string]NeighborInfo, 10)
	intf.incompleteNbr = make(map[string]struct{})

	// set counters to zero
	intf.counter.Send = 0
	intf.counter.Rcvd = 0
}

/*
 * Init Interface will be called during bootup when we do Get ALL ipv6 intfs
 */
func (intf *Interface) InitIntf(obj *config.IPIntfNotification, pktCh chan config.PacketData, gCfg NdpConfig) {
	intf.IntfRef = obj.IntfRef
	intf.IfIndex = obj.IfIndex
	intf.OperState = obj.Operation
	intf.commonInit(obj.IpAddr, pktCh, gCfg)
	debug.Logger.Debug("InitIntf port:", obj.IntfRef, "ifIndex:", obj.IfIndex, "GS:", intf.IpAddr, "LS:", intf.LinkLocalIp)
}

/*
 * De-Init Interface will be called during delete ip interface
 */
func (intf *Interface) DeInitIntf() []string {
	deleteEntries, _ := intf.DeleteAll()
	intf.PcapBase.PcapHandle = nil
	//intf.PcapBase.PcapCtrl = nil
	intf.PcapBase.PcapUsers = 0
	intf.removeIP(intf.IpAddr)
	intf.removeIP(intf.LinkLocalIp)
	// Timers Value De-Init
	intf.raTimer = nil
	// Delete Nbrmap
	intf.Neighbor = nil
	return deleteEntries
}

/*
 * CreateIntf is called during CreateIPInterface notification
 */
func (intf *Interface) CreateIntf(obj *config.IPIntfNotification, pktCh chan config.PacketData, gCfg NdpConfig) {
	intf.IntfRef = obj.IntfRef
	intf.IfIndex = obj.IfIndex
	intf.commonInit(obj.IpAddr, pktCh, gCfg)
	debug.Logger.Debug("Created IP interface", intf.IntfRef, "ifIndex:", intf.IfIndex,
		"GS:", intf.globalScope, "LS:", intf.linkScope, "Pcap Users are:", intf.PcapUsers)
}

/*
 * If Entry Already exists during CreateIPInterface then Update Interface will be called
 */
func (intf *Interface) UpdateIntf(ipAddr string) {
	intf.addIP(ipAddr)
	debug.Logger.Debug("UpdateIntf port:", intf.IntfRef, "ifIndex:", intf.IfIndex, "GS:", intf.IpAddr, "LS:", intf.LinkLocalIp)
}

/* set if type for the l3 port
 */
func (intf *Interface) SetIfType(ifType int) {
	intf.IfType = ifType
}

func (intf *Interface) deleteNbrList() ([]string, error) {
	if intf.PcapBase.PcapHandle == nil && intf.PcapBase.PcapUsers == 0 {
		intf.StopRATimer()
		deleteEntries, err := intf.FlushNeighbors()
		return deleteEntries, err
	}
	return make([]string, 0), nil
}

/*
 * DeleteIntf will kill pcap, flush neighbors and then stop all timers
 */
func (intf *Interface) DeleteIntf(ipAddr string) ([]string, error) {
	debug.Logger.Debug("Deleting Interface Called for:", intf.IntfRef, intf.IpAddr, intf.LinkLocalIp)
	intf.DeletePcap()
	return intf.deleteNbrList()
}

/*
 * Delete All will delete ip address and then remove entire pcap
 */
func (intf *Interface) DeleteAll() ([]string, error) {
	intf.DeletePcap()
	intf.DeletePcap()
	return intf.deleteNbrList()
}

/*
 * API: will create pcap handler for each port
 *		1) check if pcap users are > 0.. if so then just add pcap user and move on
 *		2) if no pcap users then check for PcapHandler and then create a new pcap handler
 *		3) Check if PcapCtrl is created or not..
 */
func (intf *Interface) CreatePcap() (err error) {
	// RX Pcap handler for interface
	if intf.PcapBase.PcapHandle == nil {
		name := intf.IntfRef
		intf.PcapBase.PcapHandle, err = pcap.OpenLive(name, NDP_PCAP_SNAPSHOTlEN, NDP_PCAP_PROMISCUOUS, NDP_PCAP_TIMEOUT)
		if err != nil {
			debug.Logger.Err("Creating Pcap Handler failed for interface:", name, "Error:", err)
			return err
		}
		err = intf.PcapBase.PcapHandle.SetBPFFilter(NDP_PCAP_FILTER)
		if err != nil {
			debug.Logger.Err("Creating BPF Filter failed Error", err)
			intf.PcapBase.PcapHandle = nil
			return err
		}
	}
	intf.addPcapUser()
	debug.Logger.Info("Total pcap user for", intf.IntfRef, "to", intf.PcapBase.PcapUsers)
	return err
}

/*
 *  Create a TX pcap handler for sending packets on timer expiry
 */
func (intf *Interface) CreateTXPcap() (err error) {
	if intf.PcapBase.Tx == nil {
		intf.PcapBase.Tx, err = pcap.OpenLive(intf.IntfRef, NDP_PCAP_SNAPSHOTlEN, NDP_PCAP_PROMISCUOUS, NDP_PCAP_TIMEOUT)
		if err != nil {
			debug.Logger.Err("Creating TX Pcap Handler failed for interface:", intf.IntfRef, "Error:", err)
			intf.DeletePcap()
			return
		}
		debug.Logger.Info("TX Pcap created for interface:", intf.IntfRef)
	}
	return nil
}

func (intf *Interface) DeleteTXPcap() {
	// delete tx channel only if RX channel is closed && there are no pcap users
	if intf.PcapBase.PcapHandle == nil && intf.PcapBase.PcapUsers == 0 {
		if intf.PcapBase.Tx != nil {
			intf.PcapBase.Tx.Close()
			intf.PcapBase.Tx = nil
		}
	}
}

/*
 * API: add pcap users
 */
func (intf *Interface) addPcapUser() {
	intf.PcapBase.PcapUsers += 1
}

/*
 * API: add pcap users
 */
func (intf *Interface) deletePcapUser() {
	intf.PcapBase.PcapUsers -= 1
}

/*
 * DeletePcap Handler
 *	1) fpPort1 has one ip address, bypass the check and delete pcap
 *	2) fpPort1 has two ip address
 *		a) 2003::2/64 	- Global Scope
 *		b) fe80::123/64 - Link Scope
 *		In this case we will get two Notification for port down from the chip, one is for
 *		Global Scope Ip and second is for Link Scope..
 *		On first Notification NDP will update pcap users and move on. Only when second delete
 *		notification comes then NDP will delete pcap
 */
func (intf *Interface) DeletePcap() {
	if intf.PcapBase.PcapHandle == nil {
		// create ip interface but state down will not have pcap handler created
		//debug.Logger.Debug("No pcap created or it might have been deleted during l2 port down returning early")
		return
	}
	intf.deletePcapUser()
	debug.Logger.Debug("Total pcap user for", intf.IntfRef, "reduced to", intf.PcapBase.PcapUsers)
	if intf.PcapBase.PcapUsers == 0 {
		debug.Logger.Debug("No More Pcap users and hence close pcap Handler for port:", intf.IntfRef)
		// once go routine is exited, delete pcap handler
		if intf.PcapBase.PcapHandle != nil {
			// Inform go routine spawned for intf to exit..
			intf.PcapBase.PcapHandle.Close()
			intf.PcapBase.PcapHandle = nil
		}
		// deleted ctrl channel to avoid any memory usage
		intf.PcapBase.PcapUsers = 0 // set to zero
		// flushing the counter values after the pcap is deleted
		intf.counter.Send = 0
		intf.counter.Rcvd = 0
	}
}

func (intf *Interface) writePkt(pkt []byte) error {
	//if intf.PcapBase.PcapHandle != nil {
	if intf.PcapBase.Tx != nil {
		//err := intf.PcapBase.PcapHandle.WritePacketData(pkt)
		err := intf.PcapBase.Tx.WritePacketData(pkt)
		if err != nil {
			debug.Logger.Err("Sending Packet failed error:", err)
			return errors.New("Sending Packet Failed")
		}
	} else {
		debug.Logger.Warning("Pcap deleted for interface:", intf.IntfRef)
		return errors.New(fmt.Sprintln("Pcap deleted for interface:", intf.IntfRef))
	}
	return nil
}

/*
 * Receive Ndp Packet and push it on the pktCh
 */
func (intf *Interface) ReceiveNdpPkts(pktCh chan *RxPktInfo) error {
	if intf.PcapBase.PcapHandle == nil {
		debug.Logger.Err("pcap handler for port:", intf.IntfRef, "is not valid. ABORT!!!!")
		return errors.New(fmt.Sprintln("pcap handler for port:", intf.IntfRef, "is not valid. ABORT!!!!"))
	}
	src := gopacket.NewPacketSource(intf.PcapBase.PcapHandle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		select {
		case pkt, ok := <-in:
			if ok {
				pktCh <- &RxPktInfo{pkt, intf.IfIndex}
			} else {
				debug.Logger.Debug("Pcap closed as in is invalid exiting go routine for port:", intf.IntfRef)
				return nil
			}
		}
	}
	return nil
}

/*
 *  On physical link down.. server will request to delete all the neighbors from the cache..
 *  We will iterate over all the neighbors, stop its timer and delete the neighbor cache
 *  @NOTE: Always start from the bottom most layer as we never know when golang will run its garbage collector
 */
func (intf *Interface) FlushNeighbors() ([]string, error) {
	// during link local down we will have CIDR format ip
	debug.Logger.Debug("Deleting all neighbor entries for interface:", intf.IntfRef)
	deleteEntries := make([]string, 0)
	for nbrKey, nbr := range intf.Neighbor {
		nbr.DeInit()
		//deleteEntries = append(deleteEntries, nbr.IpAddr)
		deleteEntries = append(deleteEntries, nbrKey)
		debug.Logger.Debug("Deleting neighbor:", nbrKey)
		delete(intf.Neighbor, nbrKey)
	}
	// do not delete link information here... only if IP interface is deleted then we need to delete
	// link information
	return deleteEntries, nil
}

/*
 * flush neighbors per ip address
 */
func (intf *Interface) FlushNeighborPerIp(nbrKey, ipAddr string) ([]string, error) {
	debug.Logger.Debug("flushing nbr entry for ipAddr:", ipAddr, "key:", nbrKey)
	deleteEntries := make([]string, 0)
	nbr, exists := intf.Neighbor[nbrKey]
	if !exists {
		return deleteEntries, errors.New("No Neighbor found for:" + nbrKey)
	}
	nbr.DeInit()
	//deleteEntries = append(deleteEntries, ipAddr)
	deleteEntries = append(deleteEntries, nbrKey)
	delete(intf.Neighbor, nbrKey)
	return deleteEntries, nil
}

/*
 *  helper function for creating key based on nd Information
 */
func (intf *Interface) createNbrKey(ndInfo *packet.NDInfo) (nbrkey string) {
	if ndInfo.SrcIp == intf.globalScope || ndInfo.SrcIp == intf.linkScope {
		// use destination ip as index to neighbor information
		nbrkey = ndInfo.DstMac + "_" + ndInfo.DstIp + "_" + ndInfo.LearnedIntfRef
	} else {
		nbrkey = ndInfo.SrcMac + "_" + ndInfo.SrcIp + "_" + ndInfo.LearnedIntfRef
	}
	debug.Logger.Debug("nbrKey created for ipv6 adj is:", nbrkey)
	return nbrkey
}

/*
 * process nd will be called during received message
 */
func (intf *Interface) ProcessND(ndInfo *packet.NDInfo) (*config.NeighborConfig, NDP_OPERATION) {
	intf.counter.Rcvd++
	if intf.Neighbor == nil {
		debug.Logger.Alert("!!!!Neighbor Initialization for intf:", intf.IntfRef, "didn't happen properly!!!!!")
		intf.Neighbor = make(map[string]NeighborInfo, 10)
	}
	// if tx is closed then do not learn the packet
	if intf.PcapBase.Tx == nil {
		debug.Logger.Debug("TX channel is closed for:", intf.IntfRef, "hence ignoring incoming packet")
		return nil, IGNORE
	}
	debug.Logger.Debug("processing nd info:", *ndInfo)
	switch ndInfo.PktType {
	case layers.ICMPv6TypeNeighborSolicitation:
		return intf.processNS(ndInfo)
	case layers.ICMPv6TypeNeighborAdvertisement:
		return intf.processNA(ndInfo)
	case layers.ICMPv6TypeRouterAdvertisement:
		return intf.processRA(ndInfo)
	case layers.ICMPv6TypeRouterSolicitation:
		// @TODO: not supported
	}

	return nil, IGNORE
}

/*
 * send neighbor discover messages on timer expiry
 */
func (intf *Interface) SendND(pktData config.PacketData, mac string) NDP_OPERATION {
	switch pktData.SendPktType {
	case layers.ICMPv6TypeNeighborSolicitation:
		return intf.SendNS(mac, pktData.NeighborMac, pktData.NeighborIp, pktData.FastProbe)
	case layers.ICMPv6TypeNeighborAdvertisement:
		// @TODO: implement this
	case layers.ICMPv6TypeRouterAdvertisement:
		intf.SendRA(mac)
	case layers.ICMPv6TypeRouterSolicitation:
		// @TODO: ignore router solicitation
	}
	return IGNORE
}

/*
 *  Update timer values
 */
func (intf *Interface) UpdateTimer(gCfg NdpConfig) {
	intf.reachableTime = gCfg.ReachableTime
	intf.raRestransmitTime = gCfg.RaRestransmitTime
	intf.retransTime = gCfg.RetransTime
}

/*
 *  Get Neighbor Information
 */
func (intf *Interface) PopulateNeighborInfo(nbr NeighborInfo, nbrState *config.NeighborEntry) {
	debug.Logger.Debug("Neighbor Information in NDP is:", nbr)
	nbrState.IpAddr = nbr.IpAddr
	nbrState.MacAddr = nbr.LinkLayerAddress
	baseTime := (int64)(nbr.BaseReachableTimer * 60 * 1000)
	baseReachableTime := time.Duration(baseTime) * time.Millisecond
	elapsedTime := time.Since(nbr.pktRcvdTime)
	expiryTime := baseReachableTime - elapsedTime
	nbrState.ExpiryTimeLeft = expiryTime.String()
	nbrState.ReceivedPackets = nbr.counter.Rcvd
	nbrState.SendPackets = nbr.counter.Send
	switch nbr.State {
	case INCOMPLETE:
		nbrState.State = "Incomplete"
	case REACHABLE:
		nbrState.State = "Reachable"
	case STALE:
		nbrState.State = "Stale"
	case DELAY:
		nbrState.State = "Delay"
	case PROBE:
		nbrState.State = "Probe"
	}
}

/*
 *   Interface validator for nbrKey generated
 */
func (intf *Interface) validNbrKey(nbrKey string) bool {
	// mac_ip_intf
	nSplit := strings.Split(nbrKey, "_")
	// ip
	splitString := strings.Split(nSplit[1], ":")
	for _, value := range IPV6_MULTICAST_PREFIXES {
		if strings.Contains(strings.ToLower(splitString[0]), value) {
			return false
		}
	}
	return true
}

/*
 *   Api to Send Unicast Solicitation Message for all neighbor entries to do refresh before the expiry
 *   Time
 */
func (intf *Interface) RefreshAllNeighbors(mac string) {
	debug.Logger.Info("Refresh Action for All Neighbors by intferface:", intf.IntfRef)
	for _, nbr := range intf.Neighbor {
		debug.Logger.Info("Refreshing Neighbor:", nbr.LinkLayerAddress, nbr.IpAddr)
		intf.SendNS(mac, nbr.LinkLayerAddress, nbr.IpAddr, false /* isFastProbe */)
	}
}

/*
 *   Api to delete one neighbor entry request by the action
 */
func (intf *Interface) DeleteNeighbor(nbrEntry config.NeighborConfig) ([]string, error) {
	return intf.FlushNeighborPerIp(createNeighborKey(nbrEntry.MacAddr, nbrEntry.IpAddr, nbrEntry.Intf), nbrEntry.IpAddr)
}

/*
 *   UpdateNbrEntry will be called during mac move.
 *   Args:    1) oldNbrKey - old nbr that needs to be deleted...
 *	      2) newNbrKey - new nbr that needs to be initialized...
 *   return:  linkscope old key & link scope ip that got deleted with following reason:
 *	      * oldNbrKey : mac ip port
 *                          for this mac & port combination get link scope ip and do update on that entry also
 *			    return this newly formed linkScopeNbrKey to the caller so that server state can be
 *			    updated
 */
func (intf *Interface) UpdateNbrEntry(oldNbrKey, newNbrKey string) (string, string) {
	// old nbr getting deleted first
	nbr, exists := intf.Neighbor[oldNbrKey]
	if exists {
		debug.Logger.Info("Deleting neighbor:", nbr)
		nbr.DeInit()
		delete(intf.Neighbor, oldNbrKey)
	}
	// new nbr getting added
	nbr, exists = intf.Neighbor[newNbrKey]
	if !exists {
		nbr.InitCache(intf.reachableTime, intf.retransTime, newNbrKey, intf.PktDataCh, intf.IfIndex)
		nbr.RchTimer()
		nbr.State = REACHABLE
		nbr.updatePktRxStateInfo()
		intf.Neighbor[newNbrKey] = nbr
	}

	oldNbrKeySS := splitNeighborKey(oldNbrKey)
	// search for linkScopeNbrKey on the same old mac & port as we need to update that entry also
	for nbrKey, nbr := range intf.Neighbor {
		splitString := splitNeighborKey(nbrKey)
		// key is mac_ip_intf
		if oldNbrKeySS[0] == splitString[0] && oldNbrKeySS[2] == splitString[2] && isLinkLocal(splitString[1]) {
			// found the linkScope Old Entry that we want to delete
			debug.Logger.Info("Deleting Link Scope neighbor:", nbr, "key:", nbrKey)
			nbr.DeInit()
			// new nbr getting added
			newNbrKeySS := splitNeighborKey(newNbrKey)
			// new nbr mac & port but old entry ipv6 address
			linkScopeNbrKey := createNeighborKey(newNbrKeySS[0], splitString[1], newNbrKeySS[2])
			nbr, exists = intf.Neighbor[linkScopeNbrKey]
			if !exists {
				nbr.InitCache(intf.reachableTime, intf.retransTime, linkScopeNbrKey, intf.PktDataCh, intf.IfIndex)
				nbr.RchTimer()
				nbr.State = REACHABLE
				nbr.updatePktRxStateInfo()
				intf.Neighbor[linkScopeNbrKey] = nbr
			}
			delete(intf.Neighbor, nbrKey)
			return nbrKey, splitString[1]
		}
	}

	return "", ""
}

func (intf *Interface) updateFilter(macAddr string) {
	if intf.PcapBase.PcapHandle != nil {
		err := intf.PcapBase.PcapHandle.SetBPFFilter(getNewFilter(macAddr))
		if err != nil {
			debug.Logger.Err("failed to update interface:", intf.IntfRef, "filter, err:", err)
		}
	}
}

func (intf *Interface) resetFilter() {
	if intf.PcapBase.PcapHandle != nil {
		err := intf.PcapBase.PcapHandle.SetBPFFilter(NDP_PCAP_FILTER)
		if err != nil {
			debug.Logger.Err("Reseting BPF Filter failed for interface:", intf.IntfRef, "Error", err)
		}
	}
}
