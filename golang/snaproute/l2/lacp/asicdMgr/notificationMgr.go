package asicdMgr

import (
	"l2/lacp/server"
	"utils/commonDefs"
)

type NotificationHdl struct {
	Server *server.LAServer
}

func initAsicdNotification() commonDefs.AsicdNotification {
	nMap := make(commonDefs.AsicdNotification)
	nMap = commonDefs.AsicdNotification{
		commonDefs.NOTIFY_L2INTF_STATE_CHANGE: true,
		commonDefs.NOTIFY_VLAN_CREATE:         true,
		commonDefs.NOTIFY_VLAN_DELETE:         true,
		commonDefs.NOTIFY_VLAN_UPDATE:         true,
		commonDefs.NOTIFY_LAG_CREATE:          true,
		commonDefs.NOTIFY_LAG_DELETE:          true,
		commonDefs.NOTIFY_LAG_UPDATE:          true,
	}
	return nMap
}

func NewNotificationHdl(server *server.LAServer) (commonDefs.AsicdNotificationHdl, commonDefs.AsicdNotification) {
	nMap := initAsicdNotification()
	return &NotificationHdl{server}, nMap
}

func (nHdl *NotificationHdl) ProcessNotification(msg commonDefs.AsicdNotifyMsg) {
	nHdl.Server.AsicdSubSocketCh <- msg
}
