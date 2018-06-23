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

//global.go
package drcp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"net"
)

var DRGlobalSystem DRSystemInfo

type TxCallback func(key IppDbKey, dmac net.HardwareAddr, data interface{})

type DRSystemInfo struct {
	// list of tx function which should be called for a given port
	TxCallbacks map[IppDbKey][]TxCallback
}

func (g *DRSystemInfo) DRSystemGlobalRegisterTxCallback(intf IppDbKey, f TxCallback) {
	g.TxCallbacks[intf] = append(g.TxCallbacks[intf], f)
}

func (g *DRSystemInfo) DRSystemGlobalDeRegisterTxCallback(intf IppDbKey) {
	delete(g.TxCallbacks, intf)
}

func DRSystemGlobalTxCallbackListGet(p *DRCPIpp) []TxCallback {

	key := IppDbKey{
		Name:   p.Name,
		DrName: p.dr.DrniName,
	}

	if fList, pok := DRGlobalSystem.TxCallbacks[key]; pok {
		return fList
	}

	// temporary function
	x := func(key IppDbKey, dmac net.HardwareAddr, data interface{}) {
		utils.GlobalLogger.Info(fmt.Sprintf("TX not registered for IPP port\n", key))
	}

	debugTxList := make([]TxCallback, 0)
	debugTxList = append(debugTxList, x)
	return debugTxList
}
