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
	"l3/ospf/config"
	"ospfd"
	//    "l3/ospf/server"
	//    "utils/logging"
	//    "net"
)

func (h *OSPFHandler) UpdateOspfGlobal(origConf *ospfd.OspfGlobal, newConf *ospfd.OspfGlobal, attrset []bool, op []*ospfd.PatchOpInfo) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original global config attrs:", origConf))
	err := h.SendOspfGlobal(newConf)
	if err != nil {
		return false, err
	}
	h.logger.Info(fmt.Sprintln("New global config attrs:", newConf))
	return true, nil
}

func (h *OSPFHandler) UpdateOspfAreaEntry(origConf *ospfd.OspfAreaEntry, newConf *ospfd.OspfAreaEntry, attrset []bool, op []*ospfd.PatchOpInfo) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original area config attrs:", origConf))
	h.logger.Info(fmt.Sprintln("New area config attrs:", newConf))
	return true, nil
}

func (h *OSPFHandler) UpdateOspfIfEntry(origConf *ospfd.OspfIfEntry, newConf *ospfd.OspfIfEntry, attrset []bool, op []*ospfd.PatchOpInfo) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original interface config attrs:", origConf))
	h.logger.Info(fmt.Sprintln("New interface config attrs:", newConf))
	err := h.SendOspfIfConf(newConf, config.UPDATE)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *OSPFHandler) UpdateOspfIfMetricEntry(origConf *ospfd.OspfIfMetricEntry, newConf *ospfd.OspfIfMetricEntry, attrset []bool, op []*ospfd.PatchOpInfo) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original interface metric config attrs:", origConf))
	h.logger.Info(fmt.Sprintln("New interface metric config attrs:", newConf))
	return true, nil
}

/* This feature is not supported now
func (h *OSPFHandler) UpdateOspfVirtIfEntry(origConf *ospfd.OspfVirtIfEntry, newConf *ospfd.OspfVirtIfEntry, attrset []bool, op []*ospfd.PatchOpInfo) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original virtual interface config attrs:", origConf))
	h.logger.Info(fmt.Sprintln("New virtual interface config attrs:", newConf))
	return true, nil
}
*/
