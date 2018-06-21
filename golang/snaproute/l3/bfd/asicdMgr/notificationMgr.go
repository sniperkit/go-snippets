package asicdMgr

import (
	"l3/bfd/server"
	"utils/clntUtils/clntDefs/asicdClntDefs"
	"utils/clntUtils/clntIntfs"
	"utils/logging"
)

type NotificationHdl struct {
	Server *server.BFDServer
}

/*
func initAsicdNotification() commonDefs.AsicdNotification {
	nMap := make(commonDefs.AsicdNotification)
	nMap = commonDefs.AsicdNotification{
		commonDefs.NOTIFY_L2INTF_STATE_CHANGE:           false,
		commonDefs.NOTIFY_IPV4_L3INTF_STATE_CHANGE:      false,
		commonDefs.NOTIFY_IPV6_L3INTF_STATE_CHANGE:      false,
		commonDefs.NOTIFY_VLAN_CREATE:                   true,
		commonDefs.NOTIFY_VLAN_DELETE:                   true,
		commonDefs.NOTIFY_VLAN_UPDATE:                   false,
		commonDefs.NOTIFY_LOGICAL_INTF_CREATE:           false,
		commonDefs.NOTIFY_LOGICAL_INTF_DELETE:           false,
		commonDefs.NOTIFY_LOGICAL_INTF_UPDATE:           false,
		commonDefs.NOTIFY_IPV4INTF_CREATE:               false,
		commonDefs.NOTIFY_IPV4INTF_DELETE:               false,
		commonDefs.NOTIFY_IPV6INTF_CREATE:               false,
		commonDefs.NOTIFY_IPV6INTF_DELETE:               false,
		commonDefs.NOTIFY_LAG_CREATE:                    true,
		commonDefs.NOTIFY_LAG_DELETE:                    true,
		commonDefs.NOTIFY_LAG_UPDATE:                    false,
		commonDefs.NOTIFY_IPV4NBR_MAC_MOVE:              false,
		commonDefs.NOTIFY_IPV6NBR_MAC_MOVE:              false,
		commonDefs.NOTIFY_IPV4_ROUTE_CREATE_FAILURE:     false,
		commonDefs.NOTIFY_IPV4_ROUTE_DELETE_FAILURE:     false,
		commonDefs.NOTIFY_IPV6_ROUTE_CREATE_FAILURE:     false,
		commonDefs.NOTIFY_IPV6_ROUTE_DELETE_FAILURE:     false,
		commonDefs.NOTIFY_VTEP_CREATE:                   false,
		commonDefs.NOTIFY_VTEP_DELETE:                   false,
		commonDefs.NOTIFY_MPLSINTF_STATE_CHANGE:         false,
		commonDefs.NOTIFY_MPLSINTF_CREATE:               false,
		commonDefs.NOTIFY_MPLSINTF_DELETE:               false,
		commonDefs.NOTIFY_PORT_CONFIG_MODE_CHANGE:       false,
		commonDefs.NOTIFY_PORT_ATTR_CHANGE:              false,
		commonDefs.NOTIFY_IPV4VIRTUAL_INTF_CREATE:       false,
		commonDefs.NOTIFY_IPV4VIRTUAL_INTF_DELETE:       false,
		commonDefs.NOTIFY_IPV6VIRTUAL_INTF_CREATE:       false,
		commonDefs.NOTIFY_IPV6VIRTUAL_INTF_DELETE:       false,
		commonDefs.NOTIFY_IPV4_VIRTUALINTF_STATE_CHANGE: false,
		commonDefs.NOTIFY_IPV6_VIRTUALINTF_STATE_CHANGE: false,
	}
	return nMap
}
*/

func NewNotificationHdl(server *server.BFDServer, logger logging.LoggerIntf) clntIntfs.NotificationHdl {
	return &NotificationHdl{server}
}

func (nHdl *NotificationHdl) ProcessNotification(msg clntIntfs.NotifyMsg) {
	switch msg.(type) {
	case asicdClntDefs.VlanNotifyMsg,
		asicdClntDefs.LagNotifyMsg:
		nHdl.Server.AsicdSubSocketCh <- msg
	}
}
