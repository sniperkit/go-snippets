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

package statedbclient

import (
	"fmt"
	"models/objects"
	"utils/logging"
	"utils/statedbclient/flexswitch"
	"utils/statedbclient/ovs"
)

const (
	FlexSwitchPlugin = "Flexswitch"
	OVSPlugin        = "OvsDB"
)

type StateDBClient interface {
	Init() error
	AddObject(obj objects.ConfigObj) error
	DeleteObject(obj objects.ConfigObj) error
	UpdateObject(obj objects.ConfigObj) error
	DeleteAllObjects(obj objects.ConfigObj) error
}

func NewStateDBClient(plugin string, logger *logging.Writer) (StateDBClient, error) {
	var client StateDBClient
	if plugin == FlexSwitchPlugin {
		client = flexswitch.NewFSDBClient(logger)
	} else if plugin == OVSPlugin {
		client = ovs.NewOVSDBClient(logger)
	} else {
		logger.Err(fmt.Sprintf("Unknown plugin %s for State DB client", plugin))
	}

	if err := client.Init(); err != nil {
		logger.Err(fmt.Sprintf("Failed to instantiate State DB client for %s", plugin))
		return client, err
	}

	return client, nil
}
