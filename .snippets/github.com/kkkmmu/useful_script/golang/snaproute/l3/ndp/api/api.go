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
package api

import (
	"errors"
	"l3/ndp/config"
	"l3/ndp/server"
	"sync"
)

var ndpApi *NDPApiLayer = nil
var once sync.Once

type NDPApiLayer struct {
	server *server.NDPServer
}

func InitComplete() bool {
	if ndpApi == nil {
		return false
	}
	if ndpApi.server == nil {
		return false
	}
	return true
}

/*  Singleton instance should be accessible only within api
 */
func getApiInstance() *NDPApiLayer {
	once.Do(func() {
		ndpApi = &NDPApiLayer{}
	})
	return ndpApi
}

func Init(svr *server.NDPServer) {
	ndpApi = getApiInstance()
	ndpApi.server = svr
}

func SendL3PortNotification(ifIndex int32, state, ipAddr string) {
	ndpApi.server.IpIntfCh <- &config.IPIntfNotification{
		IfIndex:   ifIndex,
		Operation: state,
		IpAddr:    ipAddr,
	}
}

func SendIPIntfNotfication(ifIndex int32, ipaddr, intfRef, msgType string) {
	ndpApi.server.IpIntfCh <- &config.IPIntfNotification{
		IfIndex:   ifIndex,
		IpAddr:    ipaddr,
		IntfRef:   intfRef,
		Operation: msgType,
	}
}

func SendVlanNotification(oper string, vlanId int32, vlanIfIndex int32, vlanName string, untagPorts []int32, tagPorts []int32) {
	ndpApi.server.VlanCh <- &config.VlanNotification{
		Operation:   oper,
		VlanId:      vlanId,
		VlanIfIndex: vlanIfIndex,
		VlanName:    vlanName,
		UntagPorts:  untagPorts,
		TagPorts:    tagPorts,
	}
}

func SendMacMoveNotification(ipAddr string, ifIndex, vlanId int32) {
	ndpApi.server.MacMoveCh <- &config.MacMoveNotification{ipAddr, ifIndex, vlanId}
}

func GetAllNeigborEntries(from, count int) (int, int, []config.NeighborConfig) {
	n, c, result := ndpApi.server.GetNeighborEntries(from, count)
	return n, c, result
}

func GetNeighborEntry(ipAddr string) *config.NeighborConfig {
	return ndpApi.server.GetNeighborEntry(ipAddr)
}

func CreateGlobalConfig(vrf string, retransmit uint32, reachableTime uint32, raTime uint8) (bool, error) {
	if ndpApi.server == nil {
		return false, errors.New("Server is not initialized")
	}
	rv, err := ndpApi.server.NdpConfig.Validate(vrf, retransmit, reachableTime, raTime)
	if err != nil {
		return rv, err
	}
	ndpApi.server.GlobalCfg <- server.NdpConfig{vrf, reachableTime, retransmit, raTime}
	return true, nil
}

func UpdateGlobalConfig(vrf string, retransmit uint32, reachableTime uint32, raTime uint8) (bool, error) {
	return CreateGlobalConfig(vrf, retransmit, reachableTime, raTime)
}

func GetNDPGlobalState(vrf string) (*config.GlobalState, error) {
	return ndpApi.server.GetGlobalState(vrf), nil
}

func GetAllNdpIntfState(from, count int) (int, int, []config.InterfaceEntries) {
	n, c, result := ndpApi.server.GetInterfaceNeighborEntries(from, count)
	return n, c, result
}

func GetNdpIntfState(intfRef string) *config.InterfaceEntries {
	return ndpApi.server.GetInterfaceNeighborEntry(intfRef)
}

func SendDeleteByIfName(intfRef string) {
	ndpApi.server.ActionCh <- &config.ActionData{
		Type:    config.DELETE_BY_IFNAME,
		IntfRef: intfRef,
	}
}

func SendDeleteByNeighborIp(ipAddr string) {
	ndpApi.server.ActionCh <- &config.ActionData{
		Type:  config.DELETE_BY_IPADDR,
		NbrIp: ipAddr,
	}
}

func SendRefreshByIfName(intfRef string) {
	ndpApi.server.ActionCh <- &config.ActionData{
		Type:    config.REFRESH_BY_IFNAME,
		IntfRef: intfRef,
	}
}

func SendRefreshByNeighborIp(ipAddr string) {
	ndpApi.server.ActionCh <- &config.ActionData{
		Type:  config.REFRESH_BY_IPADDR,
		NbrIp: ipAddr,
	}
}

func SendVirtualIpNotification(IfIndex, ParentIfIndex int32, IpAddr, MacAddr, IfName, msgType string) {
	ndpApi.server.VirtualIpCh <- &config.VirtualIpInfo{
		IfIndex:       IfIndex,
		ParentIfIndex: ParentIfIndex,
		IpAddr:        IpAddr,
		MacAddr:       MacAddr,
		IfName:        IfName,
		MsgType:       msgType,
	}
}

func SendVirtualIpStateMsg(IfIndex int32, IpAddr, msgType string) {
	ndpApi.server.VirtualIpCh <- &config.VirtualIpInfo{
		IfIndex: IfIndex,
		IpAddr:  IpAddr,
		MsgType: msgType,
	}
}
