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

package api

import (
	"errors"
	"l2/lldp/config"
	"l2/lldp/server"
	"l2/lldp/utils"
	"sync"
)

type ApiLayer struct {
	server *server.LLDPServer
}

var lldpapi *ApiLayer = nil
var once sync.Once

/*  Singleton instance should be accesible only within api
 */
func getInstance() *ApiLayer {
	once.Do(func() {
		lldpapi = &ApiLayer{}
	})
	return lldpapi
}

func Init(svr *server.LLDPServer) {
	lldpapi = getInstance()
	lldpapi.server = svr
}

//@TODO: Create for LLDP Interface will only be called during auto-create if an entry is not present in the DB
// If it is present then we will read it during restart
// During update we need to check whether there is any entry in the runtime information or not
func validateExistingIntfConfig(intfRef string) (int32, bool, error) {
	ifIndex, exists := lldpapi.server.EntryExist(intfRef)
	if !exists {
		return ifIndex, exists, errors.New("Update cannot be performed for " + intfRef +
			" as LLDP Server doesn't have any info for the ifIndex")
	}
	return ifIndex, exists, nil
}

func SendIntfConfig(intfRef, txrxMode string, enable bool) (bool, error) {
	var txrxModeEnum uint8
	// Validate ifIndex before sending the config to server
	ifIndex, proceed, err := validateExistingIntfConfig(intfRef)
	if !proceed {
		return proceed, err
	}
	switch txrxMode {
	case config.TX_RX_MODE_TxRx:
		txrxModeEnum = config.TXRX
	case config.TX_RX_MODE_TxOnly:
		txrxModeEnum = config.TX_ONLY
	case config.TX_RX_MODE_RxOnly:
		txrxModeEnum = config.RX_ONLY
	default:
		return false, errors.New("Invalid TxRxMode string provided")
	}
	lldpapi.server.IntfCfgCh <- &config.IntfConfig{ifIndex, enable, txrxModeEnum}
	return proceed, err
}

func UpdateIntfConfig(intfRef, txrxMode string, enable bool) (bool, error) {
	var txrxModeEnum uint8
	ifIndex, proceed, err := validateExistingIntfConfig(intfRef)
	if !proceed {
		return proceed, err
	}
	switch txrxMode {
	case config.TX_RX_MODE_TxRx:
		txrxModeEnum = config.TXRX
	case config.TX_RX_MODE_TxOnly:
		txrxModeEnum = config.TX_ONLY
	case config.TX_RX_MODE_RxOnly:
		txrxModeEnum = config.RX_ONLY
	default:
		return false, errors.New("Invalid TxRxMode string provided")
	}
	lldpapi.server.IntfCfgCh <- &config.IntfConfig{ifIndex, enable, txrxModeEnum}
	return proceed, err
}

func SendGlobalConfig(vrf, txrxMode string, enable, snoopAndDrop bool, tranmitInterval int32) (bool, error) {
	var txrxModeEnum uint8
	if lldpapi.server.Global != nil {
		return false, errors.New("Create/Delete on Global Object is not allowed, please do Update")
	}
	switch txrxMode {
	case config.TX_RX_MODE_TxRx:
		txrxModeEnum = config.TXRX
	case config.TX_RX_MODE_TxOnly:
		txrxModeEnum = config.TX_ONLY
	case config.TX_RX_MODE_RxOnly:
		txrxModeEnum = config.RX_ONLY
	default:
		return false, errors.New("Invalid TxRxMode string provided")
	}
	debug.Logger.Debug("LLDP API received auto-create global config:", vrf, enable, tranmitInterval, txrxMode, snoopAndDrop)
	lldpapi.server.GblCfgCh <- &config.Global{vrf, enable, tranmitInterval, txrxModeEnum, snoopAndDrop}
	debug.Logger.Debug("LLDP API pushed the global config on channel and returning true to confgMgr for create")
	return true, nil
}

func UpdateGlobalConfig(vrf, txrxMode string, enable, snoopAndDrop bool, tranmitInterval int32) (bool, error) {
	var txrxModeEnum uint8
	if lldpapi.server.Global == nil {
		return false, errors.New("Update can only be performed if the global object for LLDP is created")
	}
	switch txrxMode {
	case config.TX_RX_MODE_TxRx:
		txrxModeEnum = config.TXRX
	case config.TX_RX_MODE_TxOnly:
		txrxModeEnum = config.TX_ONLY
	case config.TX_RX_MODE_RxOnly:
		txrxModeEnum = config.RX_ONLY
	default:
		return false, errors.New("Invalid TxRxMode string provided")
	}
	debug.Logger.Debug("LLDP API received global config:", vrf, enable, tranmitInterval, txrxModeEnum, snoopAndDrop)
	lldpapi.server.GblCfgCh <- &config.Global{vrf, enable, tranmitInterval, txrxModeEnum, snoopAndDrop}
	debug.Logger.Debug("LLDP API pushed the global config on channel and returning true to confgMgr")
	return true, nil
}

func SendPortStateChange(ifIndex int32, state string) {
	lldpapi.server.IfStateCh <- &config.PortState{ifIndex, state}
}

func GetIntfs(idx int, cnt int) (int, int, []config.Intf) {
	n, c, result := lldpapi.server.GetIntfs(idx, cnt)
	return n, c, result
}

func GetIntfStates(idx int, cnt int) (int, int, []config.IntfState) {
	n, c, result := lldpapi.server.GetIntfStates(idx, cnt)
	return n, c, result
}

func GetIntfState(intfRef string) *config.IntfState {
	return lldpapi.server.GetIntfState(intfRef)
}

func UpdateCache(sysInfo *config.SystemInfo) {
	lldpapi.server.UpdateCacheCh <- sysInfo
}

func GetLLDPGlobalState(vrf string) (*config.GlobalState, error) {
	return lldpapi.server.GetGlobalState(vrf), nil
}

func SendPortAttrChange(ifIndex int32, desc string, pvid int32) {
	if pvid == 0 {
		lldpapi.server.PortAttrCh <- &config.PortAttrMsg{ifIndex, desc, config.INVALID_PVID}
	} else {
		lldpapi.server.PortAttrCh <- &config.PortAttrMsg{ifIndex, desc, pvid}
	}
}
