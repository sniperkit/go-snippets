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
package vxlan

import (
	"utils/logging"
)

var ClientIntf []VXLANClientIntf
var logger *logging.Writer

// vtep vlan membership
var PortConfigMap map[int32]*PortConfig

// set the global logger instance
func SetLogger(l *logging.Writer) {
	logger = l
}

func GetLogger() *logging.Writer {
	return logger
}

func init() {
	// initialize the various db maps
	vtepDB = make(map[VtepDbKey]*VtepDbEntry, 0)
	vtepDbList = make([]*VtepDbEntry, 0)
	vxlanDB = make(map[uint32]*VxlanDbEntry, 0)
	vxlanDbList = make([]*VxlanDbEntry, 0)

	PortConfigMap = make(map[int32]*PortConfig, 0)
	portDB = make(map[string]*VxlanPort, 0)

	VxlanVtepMachineStrStateMapInit()

}
