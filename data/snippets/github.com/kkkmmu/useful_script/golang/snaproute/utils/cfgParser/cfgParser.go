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

package cfgParser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

func GetDmnPortFromClientJson(dmnName string, paramsFile string) (int, error) {
	var clientsList []ClientJson
	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		return 0, err
	}

	for _, client := range clientsList {
		if client.Name == dmnName {
			return client.Port, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("Unable to find the dmn %s in %s file", dmnName, paramsFile))
}

type ClntInfo struct {
	ClntDmnName string `json:ClntDmnName`
	PluginName  string `json:PluginName`
}

type PluginInfo struct {
	PluginName string `json:PluginName`
	ParamsFile string `json:ParamsFile`
}

type ClntInfos struct {
	PluginInfoList []PluginInfo `json:PluginInfoList`
	ClntInfoList   []ClntInfo   `json:ClntInfoList`
}

func GetDmnClntInfoFromClntInfoJson(dmnName string, fileName string) (string, string, error) {
	var clntInfos ClntInfos
	var pluginName string
	var paramsFile string

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", "", err
	}

	err = json.Unmarshal(bytes, &clntInfos)
	if err != nil {
		return "", "", err
	}

	for _, clntInfo := range clntInfos.ClntInfoList {
		if clntInfo.ClntDmnName == dmnName {
			pluginName = clntInfo.PluginName
		}
	}

	for _, pluginInfo := range clntInfos.PluginInfoList {
		if pluginName == pluginInfo.PluginName {
			paramsFile = pluginInfo.ParamsFile
		}
	}
	if pluginName != "" && paramsFile != "" {
		return pluginName, paramsFile, nil
	}
	return "", "", errors.New(fmt.Sprintf("Unable to find the dmn clnt %s in %s file", dmnName, fileName))
}
