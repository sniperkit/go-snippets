//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package clntIntfs

import (
	"utils/cfgParser"
	"utils/logging"
)

type NotificationHdl interface {
	ProcessNotification(msg NotifyMsg)
}

type NotifyMsg interface {
}

type BaseClntInitParams struct {
	Logger     logging.LoggerIntf
	NHdl       NotificationHdl
	ParamsFile string
	PluginName string
}

const (
	ClntInfoFile   string = "clntInfo.json"
	FlexswitchClnt string = "Flexswitch"
	DellCPSClnt    string = "DellCPS"
)

func NewBaseClntInitParams(dmnName string, logger logging.LoggerIntf, nHdl NotificationHdl, paramsDir string) (*BaseClntInitParams, error) {
	var (
		pluginName, paramsFile string
		err                    error
	)

	clntInfoFile := paramsDir + ClntInfoFile
	pluginName, paramsFile, err = cfgParser.GetDmnClntInfoFromClntInfoJson(dmnName, clntInfoFile)
	if err != nil {
		return nil, err
	}
	return &BaseClntInitParams{
		Logger:     logger,
		NHdl:       nHdl,
		ParamsFile: paramsFile,
		PluginName: pluginName,
	}, nil
}
