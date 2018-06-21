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
package drcp

import (
	"fmt"
	"l2/lacp/protocol/utils"
	"strings"
)

func (dr *DistributedRelay) LaDrLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR", fmt.Sprintf("%s", dr.DrniName), msg}, ":"))
	}
}

func (p *DRCPIpp) LaIppLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR IPP", fmt.Sprintf("%s", p.Name), msg}, ":"))
	}
}

func (rxm *RxMachine) DrcpRxmLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR RXM", fmt.Sprintf("%s", rxm.p.Name), msg}, ":"))
	}
}

func (ptxm *PtxMachine) DrcpPtxmLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR PTXM", fmt.Sprintf("%s", ptxm.p.Name), msg}, ":"))
	}
}

func (psm *PsMachine) DrcpPsmLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR PSM", fmt.Sprintf("%s", psm.dr.DrniName), msg}, ":"))
	}
}

func (gm *GMachine) DrcpGmLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR GM", fmt.Sprintf("%s", gm.dr.DrniName), msg}, ":"))
	}
}

func (am *AMachine) DrcpAmLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR AM", fmt.Sprintf("%s", am.dr.DrniName), msg}, ":"))
	}
}

func (txm *TxMachine) DrcpTxmLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR TXM", fmt.Sprintf("%s", txm.p.Name), msg}, ":"))
	}
}

func (nism *NetIplShareMachine) DrcpNetIplSharemLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR NETIPLSHARE", fmt.Sprintf("%s", nism.p.Name), msg}, ":"))
	}
}

func (iam *IAMachine) DrcpIAmLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR IAM", fmt.Sprintf("%s", iam.p.Name), msg}, ":"))
	}
}

func (igm *IGMachine) DrcpIGmLog(msg string) {
	if utils.GlobalLogger != nil {
		utils.GlobalLogger.Info(strings.Join([]string{"DR IGM", fmt.Sprintf("%s", igm.p.Name), msg}, ":"))
	}
}
