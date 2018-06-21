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

// server.go
package server

import (
	"fmt"
	"infra/sysd/sysdCommonDefs"
	vxlan "l3/tunnel/vxlan/protocol"
	"net"
	"utils/keepalive"
	"utils/logging"
)

var SwitchMac [6]uint8
var NetSwitchMac net.HardwareAddr
var VxlanServer *VXLANServer

type VXLANServer struct {
	logger         *logging.Writer
	Configchans    *vxlan.VxLanConfigChannels
	DaemonStatusCh chan sysdCommonDefs.DaemonStatus
	Paramspath     string // location of params path
	CfgHandlerEvt  chan bool
}

func NewVXLANServer(l *logging.Writer, paramspath string, cfghandlerevt chan bool) *VXLANServer {

	if VxlanServer == nil {

		// setup server to monitor the daemons vxlan have a depenance on
		go keepalive.InitKeepAlive("vxland", paramspath)

		daemonstatuslistener := keepalive.InitDaemonStatusListener()
		if daemonstatuslistener != nil {
			go daemonstatuslistener.StartDaemonStatusListner()
		}

		VxlanServer = &VXLANServer{
			Paramspath: paramspath,
			logger:     l,
			Configchans: &vxlan.VxLanConfigChannels{
				Vxlancreate:               make(chan vxlan.VxlanConfig, 0),
				Vxlandelete:               make(chan vxlan.VxlanConfig, 0),
				Vxlanupdate:               make(chan vxlan.VxlanUpdate, 0),
				Vtepcreate:                make(chan vxlan.VtepConfig, 0),
				Vtepdelete:                make(chan vxlan.VtepConfig, 0),
				Vtepupdate:                make(chan vxlan.VtepUpdate, 0),
				VxlanAccessPortVlanUpdate: make(chan vxlan.VxlanAccessPortVlan, 0),
				VxlanNextHopUpdate:        make(chan vxlan.VxlanNextHopIp, 0),
				VxlanPortCreate:           make(chan vxlan.PortConfig, 0),
			},
			DaemonStatusCh: daemonstatuslistener.DaemonStatusCh,
			CfgHandlerEvt:  cfghandlerevt,
		}

		// listen for config messages from intf and server listener (thrift)
		VxlanServer.ConfigListener()

		// connect to the various servers in order to get additional information
		// such as connecting to RIB for next hop ip of the vtep dst ip, and
		// resolve the mac for the next hop ip
		for _, client := range vxlan.ClientIntf {
			client.SetServerChannels(VxlanServer.Configchans)
			err := client.ConnectToClients(paramspath+"clients.json", "")
			VxlanServer.logger.Info("NewServer: (ConnectToClients)", err)
			if err == nil {
				VxlanServer.CfgHandlerEvt <- true
			}
		}
	}
	return VxlanServer
}

// UpdateDaemonConnectDisconnect will take appropriate action when a server (daemon)
// goes down or up during the time at which vxland is up
func (s *VXLANServer) UpdateDaemonConnectDisconnect(daemonstatus sysdCommonDefs.DaemonStatus) {
	if daemonstatus.Name == "asicd" ||
		daemonstatus.Name == "ribd" ||
		daemonstatus.Name == "arpd" {
		s.logger.Info("Daemon Status", daemonstatus.Name, sysdCommonDefs.ConvertDaemonStateCodeToString(daemonstatus.Status))
		for _, client := range vxlan.ClientIntf {
			if daemonstatus.Status == sysdCommonDefs.UP {
				// should not return till all clients are connected
				err := client.ConnectToClients(s.Paramspath+"clients.json", daemonstatus.Name)

				if daemonstatus.Name == "ribd" {
					for _, vtep := range vxlan.GetVtepDBList() {
						// there is an assumtion here and it may be bad that
						// rib will repopulate with the correct reachability
						if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() >= vxlan.VxlanVtepStateNextHopInfo {
							client.RegisterReachability(vtep.DstIp)
						}
					}
				} else if daemonstatus.Name == "arpd" {
					// nothing to do?
				} else if daemonstatus.Name == "asicd" {
					for _, vxlan := range vxlan.GetVxlanDBList() {
						client.CreateVxlan(vxlan)
					}
					for _, vtep := range vxlan.GetVtepDBList() {
						returnEvt := make(chan vxlan.MachineEvent)
						if vtep.VxlanVtepMachineFsm.Machine.Curr.CurrentState() >= vxlan.VxlanVtepStateHwConfig {
							client.CreateVtep(vtep, returnEvt)
							<-returnEvt
						}
					}
				}
				// lets start listening to config msgs from confd now that all the
				// clients have been connected
				s.logger.Info("Clients connected", err)

				if err == nil {
					s.logger.Info("Sending Evt to Start Cfg handler loop")
					s.CfgHandlerEvt <- true
				}

			} else if daemonstatus.Status == sysdCommonDefs.RESTARTING ||
				daemonstatus.Status == sysdCommonDefs.STOPPED {
				// stop thrift server loop so that
				// confd does not send any config messages
				s.CfgHandlerEvt <- false
				client.DisconnectFromClients(daemonstatus.Name)
			}
		}
	}

}

func (s *VXLANServer) ConfigListener() {

	go func(cc *vxlan.VxLanConfigChannels) {
		for {
			select {

			case daemonstatus := <-s.DaemonStatusCh:
				s.UpdateDaemonConnectDisconnect(daemonstatus)

			case v := <-cc.Vxlancreate:
				vxlan.CreateVxLAN(&v)

			case v := <-cc.Vxlandelete:
				vxlan.DeleteVxLAN(&v, false)

			case v := <-cc.Vxlanupdate:
				vxlan.UpdateThriftVxLAN(&v)

			case v := <-cc.Vtepcreate:
				vxlan.CreateVtep(&v)

			case v := <-cc.Vtepdelete:
				vxlan.DeleteVtep(&v)

			case v := <-cc.Vtepupdate:
				vxlan.UpdateThriftVtep(&v)

			case <-cc.VxlanAccessPortVlanUpdate:
				// updates from client which are post create of vxlan

			case ipinfo := <-cc.VxlanNextHopUpdate:
				// updates from client which are triggered post create of vtep
				reachable := false
				if ipinfo.Command == vxlan.VxlanCommandCreate {
					reachable = true
				}
				//ip := net.ParseIP(fmt.Sprintf("%s.%s.%s.%s", uint8(ipinfo.Ip>>24&0xff), uint8(ipinfo.Ip>>16&0xff), uint8(ipinfo.Ip>>8&0xff), uint8(ipinfo.Ip>>0&0xff)))
				vxlan.HandleNextHopChange(ipinfo.Ip, ipinfo.NextHopIp, ipinfo.Intf, ipinfo.IntfName, reachable)

			case port := <-cc.VxlanPortCreate:
				// store all the valid physical ports
				if _, ok := vxlan.PortConfigMap[port.IfIndex]; !ok {
					var portcfg = &vxlan.PortConfig{
						Name:         port.Name,
						HardwareAddr: port.HardwareAddr,
						Speed:        port.Speed,
						PortNum:      port.PortNum,
						IfIndex:      port.IfIndex,
					}
					//s.logger.Info("Saving Port Config to db", *portcfg)
					vxlan.PortConfigMap[port.IfIndex] = portcfg
				}
			case intfinfo := <-cc.Vxlanintfinfo:
				for _, vtep := range vxlan.GetVtepDB() {
					s.logger.Info(fmt.Sprintln("received intf info", intfinfo, vtep))
					if vtep.SrcIfName == intfinfo.IntfName {

						vtep.VxlanVtepMachineFsm.VxlanVtepEvents <- vxlan.MachineEvent{
							E:    vxlan.VxlanVtepEventSrcInterfaceResolved,
							Src:  vxlan.VxlanVtepMachineModuleStr,
							Data: intfinfo,
						}
					}
				}
			}
		}
	}(s.Configchans)
}
