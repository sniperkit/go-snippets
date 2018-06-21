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

package rpc

import (
	"fmt"
	"ospfd"
	//    "l3/ospf/config"
	//    "l3/ospf/server"
	//    "utils/logging"
	//    "net"
)

func (h *OSPFHandler) DeleteOspfGlobal(ospfGlobalConf *ospfd.OspfGlobal) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete global config attrs:", ospfGlobalConf))
	return true, nil
}

func (h *OSPFHandler) DeleteOspfAreaEntry(ospfAreaConf *ospfd.OspfAreaEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete Area Config attrs:", ospfAreaConf))
	return true, nil
}

func (h *OSPFHandler) DeleteOspfIfEntry(ospfIfConf *ospfd.OspfIfEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete interface config attrs:", ospfIfConf))
	return true, nil
}

func (h *OSPFHandler) DeleteOspfIfMetricEntry(ospfIfMetricConf *ospfd.OspfIfMetricEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete interface metric config attrs:", ospfIfMetricConf))
	return true, nil
}

/* This feature is not supported now
func (h *OSPFHandler) DeleteOspfVirtIfEntry(ospfVirtIfConf *ospfd.OspfVirtIfEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete virtual interface config attrs:", ospfVirtIfConf))
	return true, nil
}
*/
