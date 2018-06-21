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
package arpClient

import (
	"fmt"
	"l3/arp/clientMgr/flexswitch"
	"l3/arp/clientMgr/mock"
	"utils/commonDefs"
)

type ArpdClientIntf interface {
	ResolveArpIPV4(destNetIp string, vlanid int32) (err error)
	DeleteResolveArpIPv4(NbrIP string) (err error)
	DeleteArpEntry(ipAddr string) (err error)
	SendGarp(ifName string, macAddr string, ipAddr string) (err error)
}

func NewArpdClient(plugin, paramsFile string, clntList []commonDefs.ClientJson, arpHdl flexswitch.ArpdClientStruct) ArpdClientIntf {
	switch plugin {
	case commonDefs.FLEXSWITCH_PLUGIN:
		clntHdl := flexswitch.GetArpdThriftClientHdl(paramsFile, clntList)
		if clntHdl == nil {
			fmt.Println("Unable to Connect to Arpd Client")
			return nil
		}
		flexswitch.InitArpdSubscriber(arpHdl)
		return &flexswitch.FSArpdClientMgr{clntHdl}
	case commonDefs.MOCK_PLUGIN:
		return &mockArp.MockArpClient{}
	default:
		return nil
	}

	return nil
}
