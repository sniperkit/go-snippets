// vxlanArpd.go
package snapclient

import (
	"arpd"
	"arpdInt"
	"fmt"
	vxlan "l3/tunnel/vxlan/protocol"
	"net"
	"strings"
)

type ArpdClient struct {
	VXLANClientBase
	ClientHdl *arpd.ARPDServicesClient
}

var arpdclnt ArpdClient

func (intf VXLANSnapClient) ResolveNextHopMac(nexthopip net.IP, nexthopIfName string, macchan chan<- vxlan.MachineEvent) {
	if arpdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		arpentrystate, err := arpdclnt.ClientHdl.GetArpEntryState(nexthopip.String())
		logger.Debug(fmt.Sprintln("calling GetArpEntryState", nexthopip, nexthopip.String(), arpentrystate, err))
		if err == nil && !strings.Contains(arpentrystate.MacAddr, "incomplete") && arpentrystate.MacAddr != "" {
			nexthopmac, _ := net.ParseMAC(arpentrystate.MacAddr)
			event := vxlan.MachineEvent{
				E:    vxlan.VxlanVtepEventNextHopInfoMacResolved,
				Src:  vxlan.VXLANSnapClientStr,
				Data: nexthopmac,
			}
			macchan <- event
		} else {
			logger.Debug(fmt.Sprintln("calling ResolveArpIPV4", nexthopip, nexthopIfName))
			portstate, _ := asicdclnt.ClientHdl.GetPortState(nexthopIfName)
			//arpdclnt.ClientHdl.ResolveArpIPV4(nexthopip.String(), arpdInt.Int(portstate.Pvid))
			arpdclnt.ClientHdl.ResolveArpIPv4(nexthopip.String(), arpdInt.Int(portstate.IfIndex))
		}
		intf.thriftmutex.Unlock()

	}
}

func (intf VXLANSnapClient) UnresolveNextHopMac(nexthopip net.IP, nexthopIfName string) {
	if arpdclnt.ClientHdl != nil {
		intf.thriftmutex.Lock()
		arpdclnt.ClientHdl.DeleteResolveArpIPv4(nexthopip.String())
		intf.thriftmutex.Unlock()
	}
}
