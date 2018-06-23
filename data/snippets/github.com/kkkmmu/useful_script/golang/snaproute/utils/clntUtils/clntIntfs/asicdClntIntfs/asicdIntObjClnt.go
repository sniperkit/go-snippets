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

package asicdClntIntfs

import (
	"utils/clntUtils/clntDefs/asicdClntDefs"
)

type AsicdIntObjClntIntf interface {
	CreateIPv4Neighbor(ipAddr string, macAddr string, vlanId int32, ifIndex int32) (int32, error)
	UpdateIPv4Neighbor(ipAddr string, macAddr string, vlanId int32, ifIndex int32) (int32, error)
	DeleteIPv4Neighbor(ipAddr string, macAddr string, vlanId int32, ifIndex int32) (int32, error)
	CreateIPv6Neighbor(ipAddr string, macAddr string, vlanId int32, ifIndex int32) (int32, error)
	UpdateIPv6Neighbor(ipAddr string, macAddr string, vlanId int32, ifIndex int32) (int32, error)
	DeleteIPv6Neighbor(ipAddr string, macAddr string, vlanId int32, ifIndex int32) (int32, error)

	GetBulkVlan(fromIndex int, count int) (*asicdClntDefs.VlanGetInfo, error)

	CreateLag(ifname string, hashType int32, ifIndexList string) (int32, error)
	DeleteLag(ifIndex int32) (int32, error)
	UpdateLag(ifIndex, hashType int32, ifIndexList string) (int32, error)
	CreateLagCfgIntfList(ifName string, ifIndexList []int32) (bool, error)
	UpdateLagCfgIntfList(ifName string, ifIndexList []int32) (bool, error)
	DeleteLagCfgIntfList(ifName string, ifIndexList []int32) (bool, error)

	GetBulkLag(fromIndex int, count int) (*asicdClntDefs.LagGetInfo, error)

	IsLinuxOnlyPlugin() (bool, error)
	GetAllPortsWithDirtyCache() ([]*asicdClntDefs.Port, error)

	OnewayCreateIPv4Route(ipv4RouteList []*asicdClntDefs.IPv4Route)
	OnewayDeleteIPv4Route(ipv4RouteList []*asicdClntDefs.IPv4Route)

	OnewayCreateIPv6Route(ipv6RouteList []*asicdClntDefs.IPv6Route)
	OnewayDeleteIPv6Route(ipv6RouteList []*asicdClntDefs.IPv6Route)
}
