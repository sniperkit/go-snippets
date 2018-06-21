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

func (h *BFDHandler) SendBfdGlobalConfig(bfdGlobalConfig *bfdd.BfdGlobal) bool {
	gConf := server.GlobalConfig{
		Vrf:    bfdGlobalConfig.Vrf,
		Enable: bfdGlobalConfig.Enable,
	}
	h.server.GlobalConfigCh <- gConf
	return true
}

func (h *BFDHandler) SendBfdSessionConfig(bfdSessionConfig *bfdd.BfdSession) bool {
	sessionConf := server.SessionConfig{
		DestIp:    bfdSessionConfig.IpAddr,
		ParamName: bfdSessionConfig.ParamName,
		Interface: bfdSessionConfig.Interface,
		PerLink:   bfdSessionConfig.PerLink,
		Protocol:  bfddCommonDefs.ConvertBfdSessionOwnerStrToVal(bfdSessionConfig.Owner),
		Operation: bfddCommonDefs.CREATE,
	}
	h.server.SessionConfigCh <- sessionConf
	return true
}

func (h *BFDHandler) SendBfdSessionParamConfig(bfdSessionParamConfig *bfdd.BfdSessionParam) bool {
	sessionParamConf := server.SessionParamConfig{
		Name:                      bfdSessionParamConfig.Name,
		LocalMultiplier:           bfdSessionParamConfig.LocalMultiplier,
		DesiredMinTxInterval:      bfdSessionParamConfig.DesiredMinTxInterval,
		RequiredMinRxInterval:     bfdSessionParamConfig.RequiredMinRxInterval,
		RequiredMinEchoRxInterval: bfdSessionParamConfig.RequiredMinEchoRxInterval,
		DemandEnabled:             bfdSessionParamConfig.DemandEnabled,
		AuthenticationEnabled:     bfdSessionParamConfig.AuthenticationEnabled,
		AuthenticationType:        h.server.ConvertBfdAuthTypeStrToVal(bfdSessionParamConfig.AuthType),
		AuthenticationKeyId:       bfdSessionParamConfig.AuthKeyId,
		AuthenticationData:        bfdSessionParamConfig.AuthData,
	}
	h.server.SessionParamConfigCh <- sessionParamConf
	return true
}

func (h *BFDHandler) CreateBfdGlobal(bfdGlobalConf *bfdd.BfdGlobal) (bool, error) {
	if bfdGlobalConf == nil {
		err := errors.New("Invalid Global Configuration")
		return false, err
	}
	return h.SendBfdGlobalConfig(bfdGlobalConf), nil
}

func (h *BFDHandler) CreateBfdSession(bfdSessionConf *bfdd.BfdSession) (bool, error) {
	if bfdSessionConf == nil {
		err := errors.New("Invalid Session Configuration")
		return false, err
	}
	h.logger.Info("Create session config attrs:", bfdSessionConf)
	return h.SendBfdSessionConfig(bfdSessionConf), nil
}

func (h *BFDHandler) CreateBfdSessionParam(bfdSessionParamConf *bfdd.BfdSessionParam) (bool, error) {
	if bfdSessionParamConf == nil {
		err := errors.New("Invalid Session Param Configuration")
		return false, err
	}
	h.logger.Info("Create session param config attrs:", bfdSessionParamConf)
	return h.SendBfdSessionParamConfig(bfdSessionParamConf), nil
}
