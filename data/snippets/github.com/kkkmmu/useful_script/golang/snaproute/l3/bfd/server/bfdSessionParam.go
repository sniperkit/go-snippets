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

package server

func (server *BFDServer) createDefaultSessionParam() error {
	paramName := "default"
	_, exist := server.bfdGlobal.SessionParams[paramName]
	if !exist {
		sessionParam := &BfdSessionParam{}
		sessionParam.state.Name = paramName
		sessionParam.state.LocalMultiplier = DEFAULT_DETECT_MULTI
		sessionParam.state.DesiredMinTxInterval = DEFAULT_DESIRED_MIN_TX_INTERVAL
		sessionParam.state.RequiredMinRxInterval = DEFAULT_REQUIRED_MIN_RX_INTERVAL
		sessionParam.state.RequiredMinEchoRxInterval = DEFAULT_REQUIRED_MIN_ECHO_RX_INTERVAL
		sessionParam.state.DemandEnabled = false
		sessionParam.state.AuthenticationEnabled = false
		server.bfdGlobal.SessionParams[paramName] = sessionParam
		server.bfdGlobal.NumSessionParams++
		server.UpdateBfdSessionsUsingParam(sessionParam.state.Name)
	}
	server.logger.Info("Created default session param")
	return nil
}

func (server *BFDServer) processSessionParamConfig(paramConfig SessionParamConfig) error {
	sessionParam, exist := server.bfdGlobal.SessionParams[paramConfig.Name]
	if !exist {
		server.logger.Info("Creating session param: ", paramConfig.Name)
		sessionParam = &BfdSessionParam{}
		server.bfdGlobal.SessionParams[paramConfig.Name] = sessionParam
	} else {
		server.logger.Info("Updating session param: ", paramConfig.Name)
	}
	sessionParam.state.Name = paramConfig.Name
	sessionParam.state.LocalMultiplier = paramConfig.LocalMultiplier
	sessionParam.state.DesiredMinTxInterval = paramConfig.DesiredMinTxInterval * 1000
	sessionParam.state.RequiredMinRxInterval = paramConfig.RequiredMinRxInterval * 1000
	sessionParam.state.RequiredMinEchoRxInterval = paramConfig.RequiredMinEchoRxInterval * 1000
	sessionParam.state.DemandEnabled = paramConfig.DemandEnabled
	sessionParam.state.AuthenticationEnabled = paramConfig.AuthenticationEnabled
	sessionParam.state.AuthenticationType = paramConfig.AuthenticationType
	sessionParam.state.AuthenticationKeyId = paramConfig.AuthenticationKeyId
	sessionParam.state.AuthenticationData = paramConfig.AuthenticationData
	if !exist {
		server.bfdGlobal.NumSessionParams++
	}
	server.UpdateBfdSessionsUsingParam(sessionParam.state.Name)
	return nil
}

func (server *BFDServer) processSessionParamDelete(paramName string) error {
	_, exist := server.bfdGlobal.SessionParams[paramName]
	if exist {
		server.logger.Info("Deleting session param: ", paramName)
		delete(server.bfdGlobal.SessionParams, paramName)
		server.bfdGlobal.NumSessionParams--
		server.UpdateBfdSessionsUsingParam(paramName)
	}
	return nil
}
