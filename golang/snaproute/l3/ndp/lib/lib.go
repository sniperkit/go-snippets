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
package lib

import (
	"fmt"
	"l3/ndp/lib/flexswitch"
	"l3/ndp/lib/mock"
	"utils/commonDefs"
)

type NdpdClientIntf interface {
	DeleteNdpEntry(ipAddr string) (err error)
}

func NewNdpClient(plugin, paramsFile string, clntList []commonDefs.ClientJson, ndpHdl flexswitch.NdpdClientStruct) NdpdClientIntf {
	switch plugin {
	case commonDefs.FLEXSWITCH_PLUGIN:
		clntHdl := flexswitch.GetNdpdThriftClientHdl(paramsFile, clntList)
		if clntHdl == nil {
			fmt.Println("Unable to Connecte to Ndpd Client")
			return nil
		}
		flexswitch.InitNdpdSubscriber(ndpHdl)
		return &flexswitch.FSNdpdClient{clntHdl}
	case commonDefs.MOCK_PLUGIN:
		return &mockNdp.MockNdpClient{}
	default:
		return nil
	}

	return nil
}
