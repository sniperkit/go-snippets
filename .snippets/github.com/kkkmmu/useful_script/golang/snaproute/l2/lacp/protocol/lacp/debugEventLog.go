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

// debugEventLog this code is meant to serialize the logging States
package lacp

import (
	//"fmt"
	"l2/lacp/protocol/utils"
	"strings"
)

func (p *LaAggPort) LaPortLog(msg string) {
	if p.logEna {
		utils.GlobalLogger.Info(strings.Join([]string{p.IntfNum, "PORT", msg}, ":"))
	}
}

func (a *LaAggregator) LacpAggLog(msg string) {
	utils.GlobalLogger.Info(strings.Join([]string{a.AggName, "AGG", msg}, ":"))
}

func (txm *LacpTxMachine) LacpTxmLog(msg string) {
	if txm.Machine.Curr.IsLoggerEna() {
		p := txm.p
		// always set to debug as this will be an excessive log
		utils.GlobalLogger.Debug(strings.Join([]string{p.IntfNum, "TXM", msg}, ":"))
	}
}

func (cdm *LacpCdMachine) LacpCdmLog(msg string) {
	if cdm.Machine.Curr.IsLoggerEna() {
		p := cdm.p
		utils.GlobalLogger.Debug(strings.Join([]string{p.IntfNum, "CDM", msg}, ":"))
	}
}

func (cdm *LacpPartnerCdMachine) LacpCdmLog(msg string) {
	if cdm.Machine.Curr.IsLoggerEna() {
		p := cdm.p
		utils.GlobalLogger.Debug(strings.Join([]string{p.IntfNum, "PCDM", msg}, ":"))
	}
}

func (ptxm *LacpPtxMachine) LacpPtxmLog(msg string) {
	if ptxm.Machine.Curr.IsLoggerEna() {
		p := ptxm.p
		utils.GlobalLogger.Info(strings.Join([]string{p.IntfNum, "PTXM", msg}, ":"))
	}
}

func (rxm *LacpRxMachine) LacpRxmLog(msg string) {
	if rxm.Machine.Curr.IsLoggerEna() {
		p := rxm.p
		// always set to debug as this will be an excessive log
		utils.GlobalLogger.Debug(strings.Join([]string{p.IntfNum, "RXM", msg}, ":"))
	}
}

func (muxm *LacpMuxMachine) LacpMuxmLog(msg string) {
	if muxm.Machine.Curr.IsLoggerEna() {
		p := muxm.p
		utils.GlobalLogger.Info(strings.Join([]string{p.IntfNum, "MUXM", msg}, ":"))
	}
}

func (mr *LampMarkerResponderMachine) LampMarkerResponderLog(msg string) {
	if mr.Machine.Curr.IsLoggerEna() {
		p := mr.p
		utils.GlobalLogger.Info(strings.Join([]string{p.IntfNum, "MARKER RESPONDER", msg}, ":"))
	}
}
