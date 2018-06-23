package server

import (
	"asicd/asicdCommonDefs"
	"fmt"
	stp "l2/stp/protocol"
	"utils/commonDefs"
	"utils/logging"
)

type STPConfigMsgType int8

const (
	STPConfigMsgCreateBridge STPConfigMsgType = iota + 1
	STPConfigMsgDeleteBridge
	STPConfigMsgUpdateBridgeMaxAge
	STPConfigMsgUpdateBridgeHelloTime
	STPConfigMsgUpdateBridgeForwardDelay
	STPConfigMsgUpdateBridgeTxHoldCount
	STPConfigMsgUpdateBridgePriority
	STPConfigMsgUpdateBridgeForceVersion
	STPConfigMsgUpdateBridgeDebugLevel
	STPConfigMsgCreatePort
	STPConfigMsgDeletePort
	STPConfigMsgUpdatePortPriority
	STPConfigMsgUpdatePortEnable
	STPConfigMsgUpdatePortPathCost
	STPConfigMsgUpdatePortProtocolMigration
	STPConfigMsgUpdatePortAdminPointToPoint
	STPConfigMsgUpdatePortAdminEdge
	STPConfigMsgUpdatePortAdminPathCost
	STPConfigMsgUpdatePortBpduGuard
	STPConfigMsgUpdatePortBridgeAssurance
	STPConfigMsgGlobalEnable
	STPConfigMsgGlobalDisable
)

type STPConfig struct {
	Msgtype STPConfigMsgType
	Msgdata interface{}
}

type STPServer struct {
	logger           *logging.Writer
	ConfigCh         chan STPConfig
	AsicdSubSocketCh chan commonDefs.AsicdNotifyMsg
}

func NewSTPServer(logger *logging.Writer) *STPServer {
	return &STPServer{
		logger:           logger,
		ConfigCh:         make(chan STPConfig),
		AsicdSubSocketCh: make(chan commonDefs.AsicdNotifyMsg),
	}
}

func (server *STPServer) InitServer() {
	//stp.ConnectToClients()
	stp.ConstructPortConfigMap()
	// TODO
	//go server.ListenToClientStateChanges()
	server.StartSTPSConfigNotificationListener()
}

/*
TODO
func (server *STPServer) ListenToClientStateChanges() {
	clientStatusListener := keepalive.InitDaemonStatusListener()
	if clientStatusListener != nil {
		go clientStatusListener.StartDaemonStatusListner()
		for {
			select {
			case clientStatus := <-clientStatusListener.DaemonStatusCh:
				mgr.logger.Info(fmt.Sprintln("Received client status: ", clientStatus.Name, clientStatus.Status))
				if mgr.IsReady() {
					switch clientStatus.Status {
					case sysdCommonDefs.STOPPED, sysdCommonDefs.RESTARTING:
						go mgr.DisconnectFromClient(clientStatus.Name)
					case sysdCommonDefs.UP:
						go mgr.ConnectToClient(clientStatus.Name)
					}
				}
			}
		}
	}
}
*/
// StartSTPSConfigNotificationListener
func (server *STPServer) StartSTPSConfigNotificationListener() {
	//server.InitDone <- true
	go func(s *STPServer) {
		stp.StpLogger("INFO", "Starting Config Event Listener")
		for {
			select {
			case stpConf, ok := <-server.ConfigCh:
				if ok {
					s.processStpConfig(stpConf)
				} else {
					// channel was closed
					return
				}
			case msg := <-server.AsicdSubSocketCh:
				s.processAsicdNotification(msg)
			}
		}
	}(server)
}

func (server *STPServer) processStpConfig(conf STPConfig) {

	switch conf.Msgtype {
	case STPConfigMsgCreateBridge:
		stp.StpLogger("INFO", "CONFIG: Create Bridge")
		config := conf.Msgdata.(*stp.StpBridgeConfig)
		stp.StpBridgeCreate(config)

	case STPConfigMsgDeleteBridge:
		stp.StpLogger("INFO", "CONFIG: Delete Bridge")
		config := conf.Msgdata.(*stp.StpBridgeConfig)
		stp.StpBridgeDelete(config)

	case STPConfigMsgUpdateBridgeMaxAge:
		stp.StpLogger("INFO", "CONFIG: Bridge Set Max Age")
		config := conf.Msgdata.(*stp.StpBridgeConfig)
		stp.StpBrgMaxAgeSet(config.IfIndex, config.MaxAge)

	case STPConfigMsgUpdateBridgeHelloTime:
		stp.StpLogger("INFO", "CONFIG: Bridge Set Hello Time")
		config := conf.Msgdata.(*stp.StpBridgeConfig)
		stp.StpBrgHelloTimeSet(config.IfIndex, config.HelloTime)

	case STPConfigMsgUpdateBridgeForwardDelay:
		stp.StpLogger("INFO", "CONFIG: Bridge Set Foward Delay")
		config := conf.Msgdata.(*stp.StpBridgeConfig)
		stp.StpBrgForwardDelaySet(config.IfIndex, config.ForwardDelay)

	case STPConfigMsgUpdateBridgeTxHoldCount:
		stp.StpLogger("INFO", "CONFIG: Bridge Set Tx Hold Count")
		config := conf.Msgdata.(*stp.StpBridgeConfig)
		stp.StpBrgTxHoldCountSet(config.IfIndex, uint16(config.TxHoldCount))

	case STPConfigMsgUpdateBridgePriority:
		stp.StpLogger("INFO", "CONFIG: Bridge Set Bridge Priority")
		config := conf.Msgdata.(*stp.StpBridgeConfig)
		stp.StpBrgPrioritySet(config.IfIndex, config.Priority)

	case STPConfigMsgUpdateBridgeForceVersion:
		stp.StpLogger("INFO", "CONFIG: Bridge Set Force Version")
		config := conf.Msgdata.(*stp.StpBridgeConfig)
		stp.StpBrgForceVersion(config.IfIndex, config.ForceVersion)

	case STPConfigMsgCreatePort:
		stp.StpLogger("INFO", "CONFIG: Port Create")
		config := conf.Msgdata.(*stp.StpPortConfig)
		stp.StpPortCreate(config)

	case STPConfigMsgDeletePort:
		stp.StpLogger("INFO", "CONFIG: Port Delete")
		config := conf.Msgdata.(*stp.StpPortConfig)
		stp.StpPortDelete(config)

	case STPConfigMsgUpdatePortPriority:
		stp.StpLogger("INFO", "CONFIG: Port Priority")
		config := conf.Msgdata.(*stp.StpPortConfig)
		stp.StpPortPrioritySet(config.IfIndex, config.BrgIfIndex, uint16(config.Priority))

	case STPConfigMsgUpdatePortEnable:
		stp.StpLogger("INFO", "CONFIG: Port Enable")
		config := conf.Msgdata.(*stp.StpPortConfig)
		stp.StpPortEnable(config.IfIndex, config.BrgIfIndex, config.Enable)

	case STPConfigMsgUpdatePortPathCost:
		stp.StpLogger("INFO", "CONFIG: Port Path Cost")
		config := conf.Msgdata.(*stp.StpPortConfig)
		stp.StpPortPortPathCostSet(config.IfIndex, config.BrgIfIndex, uint32(config.PathCost))

	case STPConfigMsgUpdatePortProtocolMigration:
		stp.StpLogger("INFO", "CONFIG: Port Protocol Migration")
		config := conf.Msgdata.(*stp.StpPortConfig)
		if config.ProtocolMigration == 1 {
			stp.StpPortProtocolMigrationSet(config.IfIndex, config.BrgIfIndex, true)
		} else {
			stp.StpPortProtocolMigrationSet(config.IfIndex, config.BrgIfIndex, false)
		}

	case STPConfigMsgUpdatePortAdminPointToPoint:
		stp.StpLogger("INFO", "CONFIG: Port Admin Point to Point UNSUPPORTED")
		//config := conf.Msgdata.(*stp.StpPortConfig)

	case STPConfigMsgUpdatePortAdminEdge:
		stp.StpLogger("INFO", "CONFIG: Port Admin Edge")
		config := conf.Msgdata.(*stp.StpPortConfig)
		stp.StpPortAdminEdgeSet(config.IfIndex, config.BrgIfIndex, config.AdminEdgePort)

	case STPConfigMsgUpdatePortAdminPathCost:
		stp.StpLogger("INFO", "CONFIG: Port Admin Path Cost UNSUPPORTED")
		//config := conf.Msgdata.(*stp.StpPortConfig)

	case STPConfigMsgUpdatePortBpduGuard:
		stp.StpLogger("INFO", "CONFIG: Port BPDU Guard")
		config := conf.Msgdata.(*stp.StpPortConfig)
		stp.StpPortBpduGuardSet(config.IfIndex, config.BrgIfIndex, config.BpduGuard)

	case STPConfigMsgUpdatePortBridgeAssurance:
		stp.StpLogger("INFO", "CONFIG: Port Bridge Assurance")
		config := conf.Msgdata.(*stp.StpPortConfig)
		stp.StpPortBridgeAssuranceSet(config.IfIndex, config.BrgIfIndex, config.BridgeAssurance)
		/*
			case STPConfigMsgGlobalEnable:
				stp.StpLogger("INFO", "CONFIG: Enable STP Global")
				stp.StpGlobalStateSet(true)
			case STPConfigMsgGlobalDisable:
				stp.StpLogger("INFO", "CONFIG: Disable STP Global")
				StpGlobalStateSet(false)
		*/
	}
}

func processLinkDownEvent(linkId int) {
	stp.StpLogger("INFO", fmt.Sprintln("STP EVT: Link Down", linkId))
	stp.StpPortLinkDown(int32(linkId))
}

func processLinkUpEvent(linkId int) {
	stp.StpLogger("INFO", fmt.Sprintln("STP EVT: Link Up", linkId))
	stp.StpPortLinkUp(int32(linkId))
}

func (server *STPServer) processAsicdNotification(msg commonDefs.AsicdNotifyMsg) {
	switch msg.(type) {
	case commonDefs.L2IntfStateNotifyMsg:
		l2Msg := msg.(commonDefs.L2IntfStateNotifyMsg)
		fmt.Printf("Msg linkstatus = %d msg port = %d\n", l2Msg.IfState, l2Msg.IfIndex)
		if l2Msg.IfState == asicdCommonDefs.INTF_STATE_DOWN {
			processLinkDownEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l2Msg.IfIndex)) //asicd always sends out link State events for PHY ports
		} else {
			processLinkUpEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l2Msg.IfIndex))
		}
	}
}
