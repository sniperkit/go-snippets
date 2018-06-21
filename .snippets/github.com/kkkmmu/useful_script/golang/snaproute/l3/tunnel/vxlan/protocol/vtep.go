// vtepdb.go
// File contains db and struct which holds the vtep db.  Also it contains the logic
// for a small FSM which will wait till Next Hop and Next Hop MAC are resolved.
// The VTEP Encap and Decap functions are defined
// The pcap listener/sender for the VTEP is also defined
package vxlan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// constants for the status of the VTEP based on FSM and provisioning
const (
	VtepStatusUp             vtepStatus = "UP"
	VtepStatusDown                      = "DOWN"
	VtepStatusAdminDown                 = "ADMIN DOWN"
	VtepStatusIncomplete                = "INCOMPLETE VTEP PROV"
	VtepStatusDetached                  = "ICOMPLETE VTEP VXLAN NOT PROV"
	VtepStatusIntfUnknown               = "SRC INTF UNKNOWN"
	VtepStatusNextHopUnknown            = "NEXT HOP UKNOWN"
	VtepStatusArpUnresolved             = "ARP UNRESOLVED"
	VtepStatusConfigPending             = "CONFIG PENDING"
)

// VtepDbKey
// Holds the key for the VtepDB
type VtepDbKey struct {
	Name  string
	Vni   uint32
	DstIp string
}

// vtepStatus
// defines the type for the status
type vtepStatus string

type VtepCounters struct {
	Rxpkts           uint64
	Rxbytes          uint64
	Rxdroppkts       uint64
	Rxdropbytes      uint64
	Rxfwdpkts        uint64
	Rxfwdbytes       uint64
	Txpkts           uint64
	Txbytes          uint64
	Txfwdpkts        uint64
	Txfwdbytes       uint64
	Txdroppkts       uint64
	Txdropbytes      uint64
	Lastrxdropreason string
	Lasttxdropreason string
}

type VtepDbEntry struct {
	// reference to the vxlan db, and value used in encap/decap
	Vni uint32
	// Key used for external reference
	VtepConfigName string
	// name of this vtep interface
	VtepName string
	// interface name which vtep will get src ip/mac info from
	SrcIfName string
	// interface id which vtep will get src ip/mac info from
	SrcIfIndex int32
	// dst UDP port
	UDP uint16
	// TTL used in the ip header in the vxlan header
	TTL uint8
	// TOS used in ip header in the vxlan header
	TOS uint8
	// MTU of the vtep
	MTU uint16
	// Source Ip used in ip header in the vxlan header
	SrcIp net.IP
	// Destination Ip used in ip header in the vxlan header
	DstIp net.IP
	// Vlan to be used to in the ethernet header in the vxlan header
	VlanId uint16
	// Source MAC to be used to in the ethernet header in the vxlan header
	SrcMac net.HardwareAddr
	// Destination MAC to be used to in the ethernet header in the vxlan header
	DstMac net.HardwareAddr
	// interface id for which packets vxlan packets should rx/tx for this vtep
	VtepIfIndex int32
	// Next Hop Ip used to find the next hop MAC which will be used as the DstMac of the vxlan header
	NextHop VtepNextHopInfo

	// Drop unknown customer macs
	FilterUnknownCustVlan bool

	// Enable/Disable state
	Enable bool

	// handle name used to rx/tx packets to linux if, also known as the Int version of the vEth dev
	VtepHandleName string
	// handle used to rx/tx packets to linux if
	handle     *pcap.Handle
	taghandles map[uint16]*pcap.Handle

	counters VtepCounters

	VxlanVtepMachineFsm *VxlanVtepMachine

	/*
		nexthopchan      chan VtepNextHopInfo
		macchan          chan net.HardwareAddr
		hwconfig         chan bool
		killroutine      chan bool
		intfinfochan     chan VxlanIntfInfo
		vteplistenerchan chan string
	*/
	// number of ticks before hw was able to come up
	ticksTillConfig int
	ticksToPollArp  int

	retrytimer *time.Timer

	// wait group used to help sync on cleanup of FSM
	wg sync.WaitGroup
}

type VtepNextHopInfo struct {
	Ip      net.IP
	IfIndex int32
	IfName  string
}

// pcap handle (vtep) per source ip defined
type VtepVniSrcIpEntry struct {
	// handle used to rx/tx packets from other applications
	handle *pcap.Handle
}

type SrcIfIndexEntry struct {
	IfIndex int32
	// handle used to rx/tx packets from/to linux if
	handle *pcap.Handle
}

// vtep id to vtep data
var vtepDB map[VtepDbKey]*VtepDbEntry
var vtepDbList []*VtepDbEntry

// vni + customer mac to vtepId
//var fdbDb map[vtepVniCMACToVtepKey]VtepDbKey

// db to hold vni ip to pcap handle
var vtepAppPcap []VtepVniSrcIpEntry

var VxlanVtepSrcIp net.IP
var VxlanVtepSrcNetMac net.HardwareAddr
var VxlanVtepSrcMac [6]uint8
var VxlanVtepRxTx = CreateVtepRxTx

func (vtep *VtepDbEntry) GetStats() VtepCounters {
	return vtep.counters
}

func GetVtepDB() map[VtepDbKey]*VtepDbEntry {
	return vtepDB
}

func GetVtepDBList() []*VtepDbEntry {
	return vtepDbList
}

func GetVtepDBEntry(key *VtepDbKey) *VtepDbEntry {
	if vtep, ok := vtepDB[*key]; ok {
		return vtep
	}
	return nil
}

func DeleteVtepDbEntryFromList(i int) {
	logger.Info("DeleteVtepDbEntryFromList", i, vtepDbList)
	j := i + 1
	copy(vtepDbList[i:], vtepDbList[j:])
	for k, n := len(vtepDbList)-j+i, len(vtepDbList); k < n; k++ {
		vtepDbList[k] = nil // or the zero value of T
	}
	vtepDbList = vtepDbList[:len(vtepDbList)-j+i]
}

func GetVtepDbListEntry(idx int32, vxlan **VtepDbEntry) bool {
	if int(idx) < len(vtepDbList) {
		*vxlan = vtepDbList[idx]
		return true
	}
	return false
}

/* TODO may need to keep a table to map customer macs to vtep
type srcMacVtepMap struct {
	SrcMac      net.HardwareAddr
	VtepIfIndex int32
}
*/

var vtepNameIdList = make([]int, 0)
var vtepNameIdCnt = 1

func GenInternalVtepName() string {
	if len(vtepNameIdList) == 0 {
		name := fmt.Sprintf("Vtep%d", vtepNameIdCnt)
		vtepNameIdCnt++
		return name
	}
	x := vtepNameIdList[0]
	vtepNameIdList = append(vtepNameIdList[:0], vtepNameIdList[1:]...)
	return fmt.Sprintf("Vtep%d", x)
}

func FreeGenInternalVtepName(vtepName string) {
	id, err := strconv.Atoi(strings.TrimLeft(vtepName, "Vtep"))
	if err == nil {
		foundEntry := false
		for _, i := range vtepNameIdList {
			if i == id {
				foundEntry = true
				break
			}
		}

		if foundEntry {
			//logger.Err("Error Deleting Vtep%d ignoring", id)
			return
		}

		vtepNameIdList = append(vtepNameIdList, id)
	}
}

func NewVtepDbEntry(c *VtepConfig) *VtepDbEntry {
	vtepName := GenInternalVtepName()
	vtep := &VtepDbEntry{
		Vni:            c.Vni,
		VtepConfigName: c.VtepName,
		VtepName:       vtepName,
		VtepHandleName: vtepName + "Int",
		SrcIfName:      c.SrcIfName,
		UDP:            c.UDP,
		TTL:            uint8(c.TTL),
		TOS:            uint8(c.TOS),
		MTU:            uint16(c.MTU),
		DstIp:          c.TunnelDstIp,
		SrcIp:          c.TunnelSrcIp,
		SrcMac:         c.TunnelSrcMac,
		DstMac:         c.TunnelDstMac,
		VlanId:         c.VlanId,
		Enable:         true,
		taghandles:     make(map[uint16]*pcap.Handle, 0),
	}
	if c.InnerVlanHandlingMode == 0 {
		vtep.FilterUnknownCustVlan = true
	} else {
		vtep.FilterUnknownCustVlan = false
	}

	return vtep
}

func CreateVtep(c *VtepConfig) *VtepDbEntry {

	vtep := saveVtepConfigData(c)
	if VxlanGlobalStateGet() == VXLAN_GLOBAL_ENABLE &&
		c.Enable {
		// lets start the FSM
		vtep.VxlanVtepMachineMain()
		vtep.VxlanVtepMachineFsm.BEGIN()
	}
	logger.Info(fmt.Sprintln("Vtep CreateVtep", *vtep))

	return vtep
}

func DeProvisionVtep(vtep *VtepDbEntry, del bool) {
	logger.Info("Calling DeprovisionVtep")
	// delete vtep resources in hw
	if vtep.VxlanVtepMachineFsm != nil {

		if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() >= VxlanVtepStateNextHopInfo {
			for _, client := range ClientIntf {
				client.UnRegisterReachability(vtep.DstIp)
			}
		}

		if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() >= VxlanVtepStateResolveNextHopMac {
			for _, client := range ClientIntf {
				client.UnresolveNextHopMac(vtep.NextHop.Ip, vtep.NextHop.IfName)
			}
		}

		if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() == VxlanVtepStateStart ||
			vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() == VxlanVtepStateHwConfig {
			for _, client := range ClientIntf {
				client.DeleteVtep(vtep)
			}
			if vtep.handle != nil {
				vtep.handle.Close()
				vtep.handle = nil
			}
			for vlan, handle := range vtep.taghandles {
				handle.Close()
				delete(vtep.taghandles, vlan)
			}
			// need to check the ref count on the port
			VxlanDelPortRxTx(vtep.NextHop.IfName, vtep.UDP)
		}
	}

	// clear out the information which was discovered for this VTEP
	vtep.NextHop.Ip = nil
	vtep.NextHop.IfIndex = 0
	vtep.NextHop.IfName = ""
	if vtep.SrcIfName != "" {
		vtep.SrcIp = nil
	}
	vtep.DstMac, _ = net.ParseMAC("00:00:00:00:00:00")

	if !del {
		if vtep.Enable {
			vxlan := GetVxlanDBEntry(vtep.Vni)
			if vxlan.Enable && vtep.VxlanVtepMachineFsm != nil {
				vtep.VxlanVtepMachineFsm.BEGIN()
			}
		} else {
			vtep.retrytimer.Stop()
		}
	}
}

func ReProvisionVtep(vtep *VtepDbEntry) {

	vtep.VxlanVtepMachineFsm.BEGIN()

	if vtep.Enable {
		// restart the timer on deprovisioning as we will retry each of the
		// state transitions again
		vtep.retrytimer.Reset(retrytime)
	}
}

func DeleteVtep(c *VtepConfig) {

	key := &VtepDbKey{
		Name:  c.VtepName,
		Vni:   c.Vni,
		DstIp: c.TunnelDstIp.String(),
	}

	vtep := GetVtepDBEntry(key)
	if vtep != nil {
		FreeGenInternalVtepName(vtep.VtepName)
		if (VxlanGlobalStateGet() == VXLAN_GLOBAL_ENABLE ||
			VxlanGlobalStateGet() == VXLAN_GLOBAL_DISABLE_PENDING) &&
			c.Enable {

			if vtep.VxlanVtepMachineFsm != nil {
				vtep.VxlanVtepMachineFsm.Stop()
				vtep.VxlanVtepMachineFsm = nil
				//DeProvisionVtep(vtep, true)

			}
			if vtep.retrytimer != nil {
				vtep.retrytimer.Stop()
				vtep.retrytimer = nil
			}
		}
		if VxlanGlobalStateGet() == VXLAN_GLOBAL_ENABLE {
			for idx, vtep := range vtepDbList {
				if vtep.VtepConfigName == c.VtepName &&
					vtep.Vni == c.Vni &&
					vtep.DstIp.String() == c.TunnelDstIp.String() {
					DeleteVtepDbEntryFromList(idx)
				}
			}

			delete(vtepDB, *key)
		}
	}
}

func saveVtepConfigData(c *VtepConfig) *VtepDbEntry {
	key := &VtepDbKey{
		Name:  c.VtepName,
		Vni:   c.Vni,
		DstIp: c.TunnelDstIp.String(),
	}
	vtep := GetVtepDBEntry(key)
	if vtep == nil {
		vtep = NewVtepDbEntry(c)
		vtepDB[*key] = vtep
		vtepDbList = append(vtepDbList, vtep)
	} else {
		vtep.SrcIfName = c.SrcIfName
		vtep.UDP = c.UDP
		vtep.TTL = uint8(c.TTL)
		vtep.TOS = uint8(c.TOS)
		vtep.MTU = uint16(c.MTU)
		vtep.DstIp = c.TunnelDstIp
		vtep.SrcIp = c.TunnelSrcIp
		vtep.SrcMac = c.TunnelSrcMac
		vtep.DstMac = c.TunnelDstMac
		vtep.VlanId = c.VlanId
		vtep.Enable = c.Enable

		if c.InnerVlanHandlingMode == 0 {
			vtep.FilterUnknownCustVlan = true
		} else {
			vtep.FilterUnknownCustVlan = false
		}
		vtepDB[*key] = vtep
		for idx, v := range vtepDbList {
			if vtep.VtepConfigName == v.VtepConfigName &&
				vtep.Vni == v.Vni &&
				vtep.DstIp.String() == v.DstIp.String() {
				DeleteVtepDbEntryFromList(idx)
			}
		}

		vtepDbList = append(vtepDbList, vtep)
	}
	return vtep
}

func SaveVtepSrcMacSrcIp(paramspath string) {
	var cfgFile cfgFileJson
	asicdconffilename := paramspath + "asicdConf.json"
	cfgFileData, err := ioutil.ReadFile(asicdconffilename)
	if err != nil {
		logger.Info("Error reading config file - asicdConf.json")
		return
	}
	err = json.Unmarshal(cfgFileData, &cfgFile)
	if err != nil {
		logger.Info("Error parsing config file")
		return
	}

	VxlanVtepSrcNetMac, _ := net.ParseMAC(cfgFile.SwitchMac)
	VxlanVtepSrcMac = [6]uint8{VxlanVtepSrcNetMac[0], VxlanVtepSrcNetMac[1], VxlanVtepSrcNetMac[2], VxlanVtepSrcNetMac[3], VxlanVtepSrcNetMac[4], VxlanVtepSrcNetMac[5]}

}

func CreateVtepRxTx(vtep *VtepDbEntry) {
	vtep.createVtepSenderListener()
}

func (vtep *VtepDbEntry) SetIfIndex(ifindex int32) {
	vtep.VtepIfIndex = ifindex
}

func (vtep *VtepDbEntry) createUntaggedVtepSenderListener() error {
	// TODO need to revisit the timeout interval in case of processing lots of
	// data frames
	handle, err := pcap.OpenLive(vtep.VtepHandleName, 65536, true, 50*time.Millisecond)
	if err != nil {
		logger.Err(fmt.Sprintf("%s: Error opening pcap.OpenLive for %s err=%s", vtep.VtepName, vtep.VtepHandleName, err))
		return err
	}
	logger.Info("Creating VXLAN Listener for intf ", vtep.VtepHandleName)
	vtep.handle = handle
	src := gopacket.NewPacketSource(vtep.handle, layers.LayerTypeEthernet)
	in := src.Packets()

	go func(rxchan chan gopacket.Packet, vlan uint16) {
		for {
			select {
			// packets received from applications which should be sent out
			case packet, ok := <-rxchan:
				if ok {
					//logger.Info("Rx untagged pkg", packet)
					if !vtep.filterPacket(packet, vlan) {
						vtep.encapAndDispatchPkt(packet)
					}
				} else {
					// channel closed
					return
				}
			}
		}
	}(in, 0)
	return nil
}

func (vtep *VtepDbEntry) createTaggedVtepSenderListener() error {
	vxlan := GetVxlanDBEntry(vtep.Vni)

	for _, vlan := range vxlan.VlanId {
		vlanhandlename := fmt.Sprintf("%s.%d", vtep.VtepHandleName, vlan)
		handle, err := pcap.OpenLive(vlanhandlename, 65536, true, 50*time.Millisecond)
		if err != nil {
			logger.Err(fmt.Sprintf("%s: Error opening pcap.OpenLive for %s err=%s", vtep.VtepName, vtep.VtepHandleName, err))
			return err
		}
		logger.Info("Creating VXLAN Listener for intf ", vlanhandlename)
		vtep.taghandles[vlan] = handle
		src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
		in := src.Packets()

		go func(rxchan chan gopacket.Packet, vlan uint16) {
			for {
				select {
				// packets received from applications which should be sent out
				case packet, ok := <-rxchan:
					if ok {
						//logger.Info("Rx tagged pkg", packet)
						if !vtep.filterPacket(packet, vlan) {
							vtep.encapAndDispatchPkt(packet)
						}
					} else {
						// channel closed
						return
					}
				}
			}
		}(in, vlan)
	}
	return nil
}

// createVtepSenderListener:
// This will listen for packets from the linux stack on the VtepHandleName
// Similarly if the MAC was learned against this VTEP traffic will be transmited
// back to the linux stack from this interface.
func (vtep *VtepDbEntry) createVtepSenderListener() error {

	// need to create a listener for each bridge instance within
	// the vni, thus one for untagged and one for each tagged
	// member
	vxlan := GetVxlanDBEntry(vtep.Vni)
	if vxlan != nil {
		if len(vxlan.UntaggedVlanId) > 0 {
			vtep.createUntaggedVtepSenderListener()
		}
		if len(vxlan.VlanId) > 0 {
			vtep.createTaggedVtepSenderListener()
		}
	}

	return nil
}

// do not process packets which contain the vtep src mac
func (vtep *VtepDbEntry) filterPacket(packet gopacket.Packet, vlan uint16) bool {

	//ethernetL := packet.Layer(layers.LayerTypeEthernet)
	//vlanL := packet.Layer(layers.LayerTypeDot1Q)
	//if ethernetL != nil {
	//	ethernet := ethernetL.(*layers.Ethernet)
	//	if vlan == 0 && vlanL != nil {
	/*
		vxlan := GetVxlanDBEntry(vtep.Vni)
		if vxlan != nil {
			dot1q := vlanL.(layers.Dot1Q)

			foundVlan := false
			for _, vlan := range vxlan.VlanId {
				if dot1q.VLANIdentifier == vlan {
					foundVlan = true
				}
			}
			if !foundVlan {
				vtep.r
			}
	*/
	//		return true
	//	} else if vlan != 0 && vlanL == nil {
	//		return true
	//	} else if vlanL != nil && vlan != 0 {
	//		dot1q := vlanL.(*layers.Dot1Q)
	//		if dot1q.VLANIdentifier != vlan {
	//			return true
	//		}
	//	}

	//logger.Info("filterPacket pkt:", ethernet.SrcMAC, "vtep:", vtep.SrcMac)
	/*if ethernet.SrcMAC[0] == vtep.SrcMac[0] &&
		ethernet.SrcMAC[1] == vtep.SrcMac[1] &&
		ethernet.SrcMAC[2] == vtep.SrcMac[2] &&
		ethernet.SrcMAC[3] == vtep.SrcMac[3] &&
		ethernet.SrcMAC[4] == vtep.SrcMac[4] &&
		ethernet.SrcMAC[5] == vtep.SrcMac[5] {
		return true
	}*/
	//}
	return false
}

func (vtep *VtepDbEntry) snoop(data []byte) {
	p2 := gopacket.NewPacket(data, layers.LinkTypeEthernet, gopacket.Default)
	ethernetL := p2.Layer(layers.LayerTypeEthernet)
	if ethernetL != nil {
		ethernet, _ := ethernetL.(*layers.Ethernet)
		learnmac := ethernet.SrcMAC
		// fdb entry mac -> vtep ip interface
		logger.Debug(fmt.Sprintln("TODO Learning mac", learnmac, "against", strings.TrimRight(vtep.VtepName, "Int")))
		//asicDLearnFwdDbEntry(learnmac, vtep.VtepName, vtep.VtepIfIndex)
	}
}

func (vtep *VtepDbEntry) decapAndDispatchPkt(packet gopacket.Packet) {

	vxlanLayer := packet.Layer(layers.LayerTypeVxlan)
	vlanLayer := packet.Layer(layers.LayerTypeDot1Q)
	if vxlanLayer != nil {
		vxlan := vxlanLayer.(*layers.VXLAN)
		vlan := uint16(0)
		if vlanLayer != nil {
			v := vlanLayer.(*layers.Dot1Q)
			vlan = v.VLANIdentifier
		}
		buf := vxlan.LayerPayload()
		//logger.Info(fmt.Sprintf("Sending Packet to %s %#v", vtep.VtepName, buf))
		vtep.snoop(buf)
		//logger.Debug(vtep.VtepName, vtep.Vni, packet, buf)
		pktlen := len(buf)
		vtep.counters.Rxpkts++
		vtep.counters.Rxbytes += uint64(pktlen)
		if vtep.handle != nil && vlan == 0 {
			if err := vtep.handle.WritePacketData(buf); err != nil {
				rc := fmt.Sprintln("Error writing packet to interface", vtep.VtepName, err)
				logger.Err(rc)
				vtep.counters.Rxdroppkts++
				vtep.counters.Rxdropbytes += uint64(pktlen)
				vtep.counters.Lastrxdropreason = rc
				return
			}
			vtep.counters.Rxfwdpkts++
			vtep.counters.Rxfwdbytes += uint64(pktlen)
		} else if vlan != 0 {
			// TODO need to forward based on customer vlan, but as of 10/28/16
			// there is only one vlan in the map so can always send it
			for vlanId, handle := range vtep.taghandles {
				if vlanId == vlan {
					if err := handle.WritePacketData(buf); err != nil {
						// This is really the Int interface but users should not care about
						// internal details
						rc := fmt.Sprintln("Error writing packet to interface", fmt.Sprintf("%s.%d", vtep.VtepName, vlanId), err)
						logger.Err(rc)
						vtep.counters.Rxdroppkts++
						vtep.counters.Rxdropbytes += uint64(pktlen)
						vtep.counters.Lastrxdropreason = rc
						return
					}
					vtep.counters.Rxfwdpkts++
					vtep.counters.Rxfwdbytes += uint64(pktlen)
				}
			}
		}
	}
}

func (vtep *VtepDbEntry) encapAndDispatchPkt(packet gopacket.Packet) {
	//logger.Info("RX packet from Vtep", packet)
	//logger.Info("encapAndDispatchPkt:", vtep.SrcMac, vtep.DstMac)
	// outer ethernet header
	eth := layers.VxlanEthernet{
		layers.Ethernet{SrcMAC: vtep.SrcMac,
			DstMAC:       vtep.DstMac,
			EthernetType: layers.EthernetTypeIPv4,
		},
	}
	ip := layers.IPv4{
		Version: 4,
		IHL:     20,
		TOS:     uint8(vtep.TOS),
		//Length:     20 + uint16(origpktlen),
		Id:         0xd2c0,
		Flags:      layers.IPv4DontFragment, //IPv4Flag
		FragOffset: 0,                       //uint16
		TTL:        uint8(vtep.TTL),
		Protocol:   layers.IPProtocolUDP, //IPProtocol
		SrcIP:      vtep.SrcIp,
		DstIP:      vtep.DstIp,
	}

	udp := layers.UDP{
		SrcPort: layers.UDPPort(vtep.UDP), // TODO need a src port
		DstPort: layers.UDPPort(vtep.UDP),
		//Length:  8 + uint16(origpktlen),
	}
	udp.SetNetworkLayerForChecksum(&ip)

	vxlan := layers.VXLAN{
		Flags: 0x08,
	}
	vxlan.SetVNI(vtep.Vni)

	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	err := gopacket.SerializeLayers(buf, opts, &eth, &ip, &udp, &vxlan, gopacket.Payload(packet.Data()))
	pktlen := len(packet.Data())

	if err == nil {
		vtep.counters.Txpkts++
		vtep.counters.Txbytes += uint64(pktlen)

		// every vtep is tied to a port
		if p, ok := portDB[vtep.NextHop.IfName]; ok {
			phandle := p.handle
			if phandle != nil {

				//logger.Info(fmt.Sprintf("Rx Packet now encapsulating and sending packet to if", vtep.SrcIfName, buf.Bytes()))
				if err := phandle.WritePacketData(buf.Bytes()); err != nil {
					rc := fmt.Sprintln("Error writing packet to interface", vtep.NextHop.IfName, err)
					logger.Err(rc)
					vtep.counters.Txdroppkts++
					vtep.counters.Txdropbytes += uint64(pktlen)
					vtep.counters.Lasttxdropreason = rc
					return
				}
				vtep.snoop(packet.Data())

				vtep.counters.Txfwdpkts++
				vtep.counters.Txfwdbytes += uint64(pktlen)
			} else {
				rc := fmt.Sprintln("Unable to find vxlan pcap handle for next hop port", vtep.NextHop.IfName)
				vtep.counters.Txdroppkts++
				vtep.counters.Txdropbytes += uint64(pktlen)
				vtep.counters.Lasttxdropreason = rc
			}
		} else {
			rc := fmt.Sprintln("Unable to find vxlan port db handle for port", vtep.NextHop.IfName)
			vtep.counters.Txdroppkts++
			vtep.counters.Txdropbytes += uint64(pktlen)
			vtep.counters.Lasttxdropreason = rc
		}
	} else {
		rc := fmt.Sprintln("Unable to encap packet", vtep.NextHop.IfName, err)
		vtep.counters.Txdroppkts++
		vtep.counters.Txdropbytes += uint64(pktlen)
		vtep.counters.Lasttxdropreason = rc
	}
}
