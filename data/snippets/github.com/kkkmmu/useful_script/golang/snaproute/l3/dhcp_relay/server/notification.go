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
	"utils/commonDefs"
	"utils/logging"
)

type NotificationHdl struct {
	dmn *DmnServer
}

func initAsicdNotification() commonDefs.AsicdNotification {
	nMap := make(commonDefs.AsicdNotification)
	nMap = commonDefs.AsicdNotification{
		commonDefs.NOTIFY_L2INTF_STATE_CHANGE:       false,
		commonDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE:  true,
		commonDefs.NOTIFY_IPV6_L3INTF_STATE_CHANGE:  true,
		commonDefs.NOTIFY_VLAN_CREATE:               false,
		commonDefs.NOTIFY_VLAN_DELETE:               false,
		commonDefs.NOTIFY_VLAN_UPDATE:               false,
		commonDefs.NOTIFY_LOGICAL_INTF_CREATE:       false,
		commonDefs.NOTIFY_LOGICAL_INTF_DELETE:       false,
		commonDefs.NOTIFY_LOGICAL_INTF_UPDATE:       false,
		commonDefs.NOTIFY_IPV4INTF_CREATE:           true,
		commonDefs.NOTIFY_IPV4INTF_DELETE:           true,
		commonDefs.NOTIFY_IPV6INTF_CREATE:           true,
		commonDefs.NOTIFY_IPV6INTF_DELETE:           true,
		commonDefs.NOTIFY_LAG_CREATE:                false,
		commonDefs.NOTIFY_LAG_DELETE:                false,
		commonDefs.NOTIFY_LAG_UPDATE:                false,
		commonDefs.NOTIFY_IPV4NBR_MAC_MOVE:          false,
		commonDefs.NOTIFY_IPV4_ROUTE_CREATE_FAILURE: false,
		commonDefs.NOTIFY_IPV4_ROUTE_DELETE_FAILURE: false,
	}
	return nMap
}

func (srvr *DmnServer) NewNotificationHdl(dmn *DmnServer, logger *logging.Writer) (commonDefs.AsicdNotificationHdl, commonDefs.AsicdNotification) {
	nMap := initAsicdNotification()
	return &NotificationHdl{dmn}, nMap
}

func (nHdl *NotificationHdl) ProcessNotification(msg commonDefs.AsicdNotifyMsg) {
	nHdl.dmn.AsicdSubSocketCh <- msg
}
