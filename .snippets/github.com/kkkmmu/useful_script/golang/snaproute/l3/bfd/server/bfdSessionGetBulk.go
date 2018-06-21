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

import (
	"strings"
)

func (server *BFDServer) GetBulkBfdSessionStates(idx int, cnt int) (int, int, []SessionState) {
	var nextIdx int
	var count int
	length := len(server.bfdGlobal.SessionsIdSlice)
	if idx+cnt >= length {
		count = length - idx
		nextIdx = 0
	} else {
		nextIdx = idx + count + 1
	}
	result := make([]SessionState, count)
	for i := idx; i < count; i++ {
		sessionId := server.bfdGlobal.SessionsIdSlice[i]
		result[i].SessionId = server.bfdGlobal.Sessions[sessionId].state.SessionId
		result[i].IpAddr = strings.Split(server.bfdGlobal.Sessions[sessionId].state.IpAddr, "%")[0]
		result[i].ParamName = server.bfdGlobal.Sessions[sessionId].state.ParamName
		result[i].Interface = server.bfdGlobal.Sessions[sessionId].state.Interface
		result[i].InterfaceSpecific = server.bfdGlobal.Sessions[sessionId].state.InterfaceSpecific
		result[i].PerLinkSession = server.bfdGlobal.Sessions[sessionId].state.PerLinkSession
		result[i].LocalMacAddr = server.bfdGlobal.Sessions[sessionId].state.LocalMacAddr
		result[i].RemoteMacAddr = server.bfdGlobal.Sessions[sessionId].state.RemoteMacAddr
		result[i].RegisteredProtocols = server.bfdGlobal.Sessions[sessionId].state.RegisteredProtocols
		result[i].SessionState = server.bfdGlobal.Sessions[sessionId].state.SessionState
		result[i].RemoteSessionState = server.bfdGlobal.Sessions[sessionId].state.RemoteSessionState
		result[i].LocalDiscriminator = server.bfdGlobal.Sessions[sessionId].state.LocalDiscriminator
		result[i].RemoteDiscriminator = server.bfdGlobal.Sessions[sessionId].state.RemoteDiscriminator
		result[i].LocalDiagType = server.bfdGlobal.Sessions[sessionId].state.LocalDiagType
		result[i].DesiredMinTxInterval = server.bfdGlobal.Sessions[sessionId].state.DesiredMinTxInterval
		result[i].RequiredMinRxInterval = server.bfdGlobal.Sessions[sessionId].state.RequiredMinRxInterval
		result[i].RemoteMinRxInterval = server.bfdGlobal.Sessions[sessionId].state.RemoteMinRxInterval
		result[i].DetectionMultiplier = server.bfdGlobal.Sessions[sessionId].state.DetectionMultiplier
		result[i].RemoteDetectionMultiplier = server.bfdGlobal.Sessions[sessionId].state.RemoteDetectionMultiplier
		result[i].DemandMode = server.bfdGlobal.Sessions[sessionId].state.DemandMode
		result[i].RemoteDemandMode = server.bfdGlobal.Sessions[sessionId].state.RemoteDemandMode
		result[i].AuthType = server.bfdGlobal.Sessions[sessionId].state.AuthType
		result[i].AuthSeqKnown = server.bfdGlobal.Sessions[sessionId].state.AuthSeqKnown
		result[i].ReceivedAuthSeq = server.bfdGlobal.Sessions[sessionId].state.ReceivedAuthSeq
		result[i].SentAuthSeq = server.bfdGlobal.Sessions[sessionId].state.SentAuthSeq
		result[i].NumTxPackets = server.bfdGlobal.Sessions[sessionId].state.NumTxPackets
		result[i].NumRxPackets = server.bfdGlobal.Sessions[sessionId].state.NumRxPackets
		result[i].ToDownCount = server.bfdGlobal.Sessions[sessionId].state.ToDownCount
		result[i].ToUpCount = server.bfdGlobal.Sessions[sessionId].state.ToUpCount
		result[i].UpTime = server.bfdGlobal.Sessions[sessionId].state.UpTime
	}
	return nextIdx, count, result
}

func (server *BFDServer) GetBfdSessionState(ipAddr string) (*SessionState, bool) {
	sessionState := new(SessionState)
	sessionId, found := server.FindBfdSessionContainingAddr(ipAddr)
	if found {
		sessionState.SessionId = server.bfdGlobal.Sessions[sessionId].state.SessionId
		sessionState.IpAddr = strings.Split(server.bfdGlobal.Sessions[sessionId].state.IpAddr, "%")[0]
		sessionState.ParamName = server.bfdGlobal.Sessions[sessionId].state.ParamName
		sessionState.Interface = server.bfdGlobal.Sessions[sessionId].state.Interface
		sessionState.Interface = server.bfdGlobal.Sessions[sessionId].state.Interface
		sessionState.InterfaceSpecific = server.bfdGlobal.Sessions[sessionId].state.InterfaceSpecific
		sessionState.PerLinkSession = server.bfdGlobal.Sessions[sessionId].state.PerLinkSession
		sessionState.LocalMacAddr = server.bfdGlobal.Sessions[sessionId].state.LocalMacAddr
		sessionState.RemoteMacAddr = server.bfdGlobal.Sessions[sessionId].state.RemoteMacAddr
		sessionState.RegisteredProtocols = server.bfdGlobal.Sessions[sessionId].state.RegisteredProtocols
		sessionState.SessionState = server.bfdGlobal.Sessions[sessionId].state.SessionState
		sessionState.RemoteSessionState = server.bfdGlobal.Sessions[sessionId].state.RemoteSessionState
		sessionState.LocalDiscriminator = server.bfdGlobal.Sessions[sessionId].state.LocalDiscriminator
		sessionState.RemoteDiscriminator = server.bfdGlobal.Sessions[sessionId].state.RemoteDiscriminator
		sessionState.LocalDiagType = server.bfdGlobal.Sessions[sessionId].state.LocalDiagType
		sessionState.DesiredMinTxInterval = server.bfdGlobal.Sessions[sessionId].state.DesiredMinTxInterval
		sessionState.RequiredMinRxInterval = server.bfdGlobal.Sessions[sessionId].state.RequiredMinRxInterval
		sessionState.RemoteMinRxInterval = server.bfdGlobal.Sessions[sessionId].state.RemoteMinRxInterval
		sessionState.DetectionMultiplier = server.bfdGlobal.Sessions[sessionId].state.DetectionMultiplier
		sessionState.RemoteDetectionMultiplier = server.bfdGlobal.Sessions[sessionId].state.RemoteDetectionMultiplier
		sessionState.DemandMode = server.bfdGlobal.Sessions[sessionId].state.DemandMode
		sessionState.RemoteDemandMode = server.bfdGlobal.Sessions[sessionId].state.RemoteDemandMode
		sessionState.AuthType = server.bfdGlobal.Sessions[sessionId].state.AuthType
		sessionState.AuthSeqKnown = server.bfdGlobal.Sessions[sessionId].state.AuthSeqKnown
		sessionState.ReceivedAuthSeq = server.bfdGlobal.Sessions[sessionId].state.ReceivedAuthSeq
		sessionState.SentAuthSeq = server.bfdGlobal.Sessions[sessionId].state.SentAuthSeq
		sessionState.NumTxPackets = server.bfdGlobal.Sessions[sessionId].state.NumTxPackets
		sessionState.NumRxPackets = server.bfdGlobal.Sessions[sessionId].state.NumRxPackets
		sessionState.ToDownCount = server.bfdGlobal.Sessions[sessionId].state.ToDownCount
		sessionState.ToUpCount = server.bfdGlobal.Sessions[sessionId].state.ToUpCount
		sessionState.UpTime = server.bfdGlobal.Sessions[sessionId].state.UpTime
	}

	return sessionState, found
}

func (server *BFDServer) GetBulkBfdSessionParamStates(idx int, cnt int) (int, int, []SessionParamState) {
	var nextIdx int
	var count int
	result := make([]SessionParamState, cnt)
	i := 0
	for _, sessionParam := range server.bfdGlobal.SessionParams {
		result[i].Name = sessionParam.state.Name
		result[i].NumSessions = sessionParam.state.NumSessions
		result[i].LocalMultiplier = sessionParam.state.LocalMultiplier
		result[i].DesiredMinTxInterval = sessionParam.state.DesiredMinTxInterval
		result[i].RequiredMinRxInterval = sessionParam.state.RequiredMinRxInterval
		result[i].RequiredMinEchoRxInterval = sessionParam.state.RequiredMinEchoRxInterval
		result[i].DemandEnabled = sessionParam.state.DemandEnabled
		result[i].AuthenticationEnabled = sessionParam.state.AuthenticationEnabled
		result[i].AuthenticationType = sessionParam.state.AuthenticationType
		result[i].AuthenticationKeyId = sessionParam.state.AuthenticationKeyId
		result[i].AuthenticationData = sessionParam.state.AuthenticationData
		i++
	}
	count = i
	nextIdx = 0
	return nextIdx, count, result
}

func (server *BFDServer) GetBfdSessionParamState(paramName string) (*SessionParamState, bool) {
	sessionParamState := new(SessionParamState)
	_, found := server.bfdGlobal.SessionParams[paramName]
	if found {
		sessionParamState.Name = server.bfdGlobal.SessionParams[paramName].state.Name
		sessionParamState.NumSessions = server.bfdGlobal.SessionParams[paramName].state.NumSessions
		sessionParamState.LocalMultiplier = server.bfdGlobal.SessionParams[paramName].state.LocalMultiplier
		sessionParamState.DesiredMinTxInterval = server.bfdGlobal.SessionParams[paramName].state.DesiredMinTxInterval
		sessionParamState.RequiredMinRxInterval = server.bfdGlobal.SessionParams[paramName].state.RequiredMinRxInterval
		sessionParamState.RequiredMinEchoRxInterval = server.bfdGlobal.SessionParams[paramName].state.RequiredMinEchoRxInterval
		sessionParamState.DemandEnabled = server.bfdGlobal.SessionParams[paramName].state.DemandEnabled
		sessionParamState.AuthenticationEnabled = server.bfdGlobal.SessionParams[paramName].state.AuthenticationEnabled
		sessionParamState.AuthenticationType = server.bfdGlobal.SessionParams[paramName].state.AuthenticationType
		sessionParamState.AuthenticationKeyId = server.bfdGlobal.SessionParams[paramName].state.AuthenticationKeyId
		sessionParamState.AuthenticationData = server.bfdGlobal.SessionParams[paramName].state.AuthenticationData
	}
	return sessionParamState, found
}
