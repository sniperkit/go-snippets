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
	"bfdd"
	"errors"
	"l3/bfd/bfddCommonDefs"
	"l3/bfd/server"
)

func (h *BFDHandler) SendBfdSessionDeleteConfig(bfdSessionConfig *bfdd.BfdSession) bool {
	sessionConf := server.SessionConfig{
		DestIp:    bfdSessionConfig.IpAddr,
		PerLink:   bfdSessionConfig.PerLink,
		Protocol:  bfddCommonDefs.ConvertBfdSessionOwnerStrToVal(bfdSessionConfig.Owner),
		Operation: bfddCommonDefs.DELETE,
	}
	h.server.SessionConfigCh <- sessionConf
	return true
}

func (h *BFDHandler) DeleteBfdGlobal(bfdGlobalConf *bfdd.BfdGlobal) (bool, error) {
	h.logger.Info("Delete global config attrs:", bfdGlobalConf)
	err := errors.New("BFD Global config delete not supported")
	return false, err
}

func (h *BFDHandler) DeleteBfdSession(bfdSessionConf *bfdd.BfdSession) (bool, error) {
	if bfdSessionConf == nil {
		err := errors.New("Invalid Session Configuration")
		return false, err
	}
	h.logger.Info("Delete session config attrs:", bfdSessionConf)
	return h.SendBfdSessionDeleteConfig(bfdSessionConf), nil
}

func (h *BFDHandler) DeleteBfdSessionParam(bfdSessionParamConf *bfdd.BfdSessionParam) (bool, error) {
	if bfdSessionParamConf == nil {
		err := errors.New("Invalid Session Param Configuration")
		return false, err
	}
	h.logger.Info("Delete session param config attrs:", bfdSessionParamConf)
	paramName := bfdSessionParamConf.Name
	h.server.SessionParamDeleteCh <- paramName
	return true, nil
}
