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

// init.go
package stp

import (
	"sync"
	"utils/logging"
)

var gLogger *logging.Writer
var portDbMutex *sync.Mutex

func init() {
	portDbMutex = &sync.Mutex{}
	PortConfigMap = make(map[int32]portConfig)
	PortMapTable = make(map[PortMapKey]*StpPort, 0)
	BridgeMapTable = make(map[BridgeKey]*Bridge, 0)
	StpPortConfigMap = make(map[int32]StpPortConfig, 0)
	StpBridgeConfigMap = make(map[int32]StpBridgeConfig, 0)

	// Init the state string maps
	TimerTypeStrStateMapInit()
	PtmMachineStrStateMapInit()
	PrxmMachineStrStateMapInit()
	PrsMachineStrStateMapInit()
	PrtMachineStrStateMapInit()
	BdmMachineStrStateMapInit()
	PimMachineStrStateMapInit()
	PpmmMachineStrStateMapInit()
	TcMachineStrStateMapInit()
	PtxmMachineStrStateMapInit()
	PstMachineStrStateMapInit()

	// create the logger used by this module
	gLogger, _ = logging.NewLogger("stpd", "STP", true)

}
