package asicdMgr

import (
	"l2/stp/server"
	"utils/commonDefs"
)

type NotificationHdl struct {
	Server *server.STPServer
}

func initAsicdNotification() commonDefs.AsicdNotification {
	nMap := make(commonDefs.AsicdNotification)
	nMap = commonDefs.AsicdNotification{
		commonDefs.NOTIFY_L2INTF_STATE_CHANGE: true,
	}
	return nMap
}

func NewNotificationHdl(server *server.STPServer) (commonDefs.AsicdNotificationHdl, commonDefs.AsicdNotification) {
	nMap := initAsicdNotification()
	return &NotificationHdl{server}, nMap
}

func (nHdl *NotificationHdl) ProcessNotification(msg commonDefs.AsicdNotifyMsg) {
	nHdl.Server.AsicdSubSocketCh <- msg
}
