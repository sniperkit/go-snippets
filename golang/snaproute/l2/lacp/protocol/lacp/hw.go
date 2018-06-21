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

// hw.go
package lacp

import (
	hwconst "asicd/asicdCommonDefs"
	"fmt"
	"l2/lacp/protocol/utils"
)

const (
	LACP_HW_INTF = "LacpHwIntf"
)

// convert the lacp port names name to asic format string list
func asicDPortBmpFormatGet(distPortList []string) string {
	s := ""
	dLength := len(distPortList)

	for i := 0; i < dLength; i++ {
		var num string

		ifindex := utils.GetIfIndexFromName(distPortList[i])
		num = fmt.Sprintf("%d", ifindex)
		if i == dLength-1 {
			s += num
		} else {
			s += num + ","
		}
	}
	return s

}

// convert the model value to asic value
func asicDHashModeGet(hashmode uint32) (laghash int32) {
	switch hashmode {
	case 0: //L2
		laghash = hwconst.HASH_SEL_SRCDSTMAC
		break
	case 1: //L2 + L3
		laghash = hwconst.HASH_SEL_SRCDSTIP
		break
	//case 2: //L3 + L4
	//break
	default:
		laghash = hwconst.HASH_SEL_SRCDSTMAC
	}
	return laghash
}

func initHwPortCreateDelCb() {
	RegisterLaPortCreateCb(LACP_HW_INTF, hwNotifyAggPortCreateDelete)
	RegisterLaPortDeleteCb(LACP_HW_INTF, hwNotifyAggPortCreateDelete)
}

func hwNotifyAggPortCreateDelete(ifIndex int32) {
	var p *LaAggPort
	if LaFindPortById(uint16(ifIndex), &p) {
		for _, client := range utils.GetAsicDPluginList() {
			var list []int32
			for _, port := range p.AggAttached.PortNumList {
				list = append(list, int32(port))
			}
			ok := client.UpdateLagCfgIntfList(p.AggAttached.AggName, list)
			if !ok {
				utils.GlobalLogger.Err("UpdateLagCfgIntfList: Error updating LAG cfg list in HW")
			}
		}
	}
}
